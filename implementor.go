/*
Copyright 2014 James Garfield. All rights reserved.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"go/ast"
	"sort"
	"strings"
)

//go:generate goast write impl gen/sliceutil.go

type Implementor struct {
	TypeProvider *Context
}

func NewImplementor(typeProvider *Context) *Implementor {
	imp := &Implementor{typeProvider}
	return imp
}

type typeSet []*ast.TypeSpec

//When implementing a generic file, types get matched from most-to-least complex
//As types are solved for, any 'subtypes' they have are also solved for
//e.g. type Collection []I gets solved and both Collection and I are known at the end
//Using this strategy eliminates a lot of otherwise ambiguous possibilities of type combinations
//because by the time you are down to the simple types, you have most likely solved them
type typesByComplexity struct {
	typeSet
	*Context
}

func (a typesByComplexity) Less(i, j int) bool {
	return a.Context.Complexity(a.typeSet[i]) > a.Context.Complexity(a.typeSet[j])
}

func (imp *Implementor) Transform(gen *Context) (result SourceSet, ok bool, errors []error) {

	isRelatedType := func(t *ast.TypeSpec) bool { return strings.Contains(t.Name.Name, "_") }
	isImplType := func(t *ast.TypeSpec) bool { return !isRelatedType(t) }

	var (
		candidateTypes typeSet = imp.TypeProvider.Types()
		genTypes       typeSet = gen.Types()
		implTypes              = genTypes.Where(isImplType)
		relatedTypes           = genTypes.Where(isRelatedType)
	)

	sort.Sort(typesByComplexity{candidateTypes, imp.TypeProvider})
	sort.Sort(typesByComplexity{implTypes, gen})

	if implTypes.Len() == 0 {
		errors = append(errors, fmt.Errorf("Invalid generic specification: No Types!"))
		return
	}

	primaryGeneric := implTypes[0]

	implContext := ContextPair{gen, imp.TypeProvider}

	//Test each type in the provider file for implementation
	for _, c := range candidateTypes {

		ok, resultMap, err := Implement(implContext, NewImplMap(), c, primaryGeneric)

		//If this candidate can't implement the primary generic type, there is no more to do
		//Save the error in case there is no implementation possible
		if !ok {
			errors = append(errors, err)
			continue
		}

		//Types can only be matched once, so no need to keep the selected candidate
		specTypes := candidateTypes.Where(func(t *ast.TypeSpec) bool {
			return t.Name.Name != c.Name.Name
		})

		matched := 0
		for _, g := range implTypes {
			//types that already have mappings are already solved for
			if _, ok := resultMap[g.Name.Name]; ok {
				matched++
				continue
			}

			//check each specification type to see if it satisfies the requirements for this generic type
			for _, s := range specTypes {
				ok, resultMap, err = Implement(implContext, resultMap, s, g)
				if ok {
					matched++
					break
				} else {
					errors = append(errors, err)
				}
			}

			//being unable to satify a generic type indicates we can't implement with the current combination of types
			//Save an error in case there is no implementation possible
			if _, ok := resultMap[g.Name.Name]; !ok {
				errors = append(errors, fmt.Errorf("Unable to satisfy generic type %s with any of the specification types.", g.Name.Name))
				break
			}
		}

		//If all generic types are mapped, rewrite the generic AST with the provided types
		if matched == implTypes.Len() {

			implAst, err := gen.Clone()
			if err != nil {
				errors = append(errors, err)
				continue
			}

			//Filter implemented types out
			//Do this prior to renaming related types so we can still identify them
			//ast.FilterFile filters out import statements...always, so use custom filter method https://github.com/golang/go/issues/9248
			filterTypeSpecs(implAst.File, func(t *ast.TypeSpec) bool { return !isImplType(t) })

			//Generate names for all related types
			relatedTypes.Each(func(t *ast.TypeSpec) {
				specName := imp.relatedTypeName(t, resultMap)
				resultMap.Store(t.Name.Name, ast.NewIdent(specName))
			})

			ast.Walk(ImplRewriter{resultMap}, implAst.File)

			imports := ImportsOfImplMap(imp.TypeProvider, resultMap)
			for _, i := range imports {
				implAst.AddImportFromSpec(i)
			}

			//ensure that implementation is in the correct package
			implAst.SetPackage(imp.TypeProvider.File.Name.Name)

			name := c.Name.Name

			result = append(result, &SourceCode{implAst, name})
		}
	}

	//If there are any results, an implementation was found and errors can be cleared
	if result.Len() != 0 {
		errors = []error{}
		ok = true
	}

	return
}

func (imp *Implementor) relatedTypeName(t *ast.TypeSpec, imap ImplMap) string {

	var (
		implExpr    ast.Expr
		found       bool
		implName    string
		partialName string = t.Name.Name
	)

	ast.Inspect(t, func(node ast.Node) bool {
		if id, ok := node.(*ast.Ident); ok {
			implExpr, found = imap[id.Name]
		}
		return !found
	})

	switch exprType := implExpr.(type) {
	case *ast.Ident:
		implName = exprType.Name

	default:
		implName = ExprString(implExpr)
	}

	return strings.Replace(partialName, "_", implName, -1)
}

//Find and return a field with a given name within a field list
func FieldByName(list *ast.FieldList, name string) (field *ast.Field, found bool) {
	for _, field = range list.List {
		for _, ident := range field.Names {
			if found = (ident.Name == name); found {
				return
			}
		}
	}
	field = nil
	return
}

func ImportsOfImplMap(ctx *Context, imap ImplMap) (result ImportSpecs) {
	for _, x := range imap {
		result = append(result, ctx.ImportsOf(x)...)
	}
	return
}
