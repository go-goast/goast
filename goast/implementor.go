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

//go:generate goast write impl --prefix=goast_ goast.net/x/iter
//go:generate goast write impl --prefix=goast_ goast.net/x/sort

type Implementor struct {
	TypeProvider *Context
}

func NewImplementor(typeProvider *Context) *Implementor {
	imp := &Implementor{typeProvider}
	return imp
}

type typeSet []*ast.TypeSpec

type implSet []ImplMap

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

		ok, primaryMap, err := Implement(implContext, NewImplMap(), c, primaryGeneric)

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

		impls := implSet{primaryMap}

		for _, g := range implTypes {
			subimpls := implSet{}

			foundMatch := false

			for _, currentMap := range impls {
				//types that already have mappings are already solved for
				if _, ok := currentMap[g.Name.Name]; ok {
					subimpls = append(subimpls, currentMap)
					if !foundMatch {
						foundMatch = true
						matched += 1
					}
					continue
				}

				//check each specification type to see if it satisfies the requirements for this generic type
				for _, s := range specTypes {
					ok, resultMap, err := Implement(implContext, currentMap, s, g)
					if ok {
						subimpls = append(subimpls, resultMap)
						if !foundMatch {
							foundMatch = true
							matched += 1
						}
						continue
					} else {
						errors = append(errors, err)
					}
				}

			}

			//being unable to satify a generic type indicates we can't implement with the current combination of types
			//Save an error in case there is no implementation possible
			if len(subimpls) == 0 {
				errors = append(errors, fmt.Errorf("Unable to satisfy generic type %s with any of the specification types.", g.Name.Name))
				break
			}

			//iterate on the new set of impls during the next iteration
			impls = subimpls

		}

		//If all generic types are mapped, rewrite the generic AST with the provided types
		if matched == implTypes.Len() {
			impPkg := &ast.Package{
				Name:  imp.TypeProvider.File.Name.Name,
				Files: make(map[string]*ast.File),
			}

			relatedImpl := map[string]bool{}

			name := c.Name.Name

			for n, currentMap := range impls {
				implAst, err := gen.Clone()
				if err != nil {
					errors = append(errors, err)
					continue
				}

				//Filter impl types & previously implemented related types out of the current ast
				//Do this prior to renaming related types so we can still identify them
				//ast.FilterFile filters out import statements...always, so use custom filter method https://github.com/golang/go/issues/9248
				filterTypeSpecs(implAst.File, func(t *ast.TypeSpec) bool {
					if isRelatedType(t) {
						specName := imp.relatedTypeName(t, currentMap)
						_, exist := relatedImpl[specName]
						relatedImpl[specName] = true
						return !exist
					}
					return !isImplType(t)
				})

				//Generate names for all related types
				relatedTypes.Each(func(t *ast.TypeSpec) {
					specName := imp.relatedTypeName(t, currentMap)
					id := ast.NewIdent(specName)
					currentMap.Store(t.Name.Name, id)
				})

				ast.Walk(ImplRewriter{currentMap}, implAst.File)

				imports := ImportsOfImplMap(imp.TypeProvider, currentMap)
				for _, i := range imports {
					implAst.AddImportFromSpec(i)
				}

				//ensure that implementation is in the correct package
				implAst.SetPackage(imp.TypeProvider.File.Name.Name)
				fileName := fmt.Sprintf("%s_%d.go", name, n)
				impPkg.Files[fileName] = implAst.File
			}

			mergedAst := ast.MergePackageFiles(impPkg, ast.FilterFuncDuplicates|ast.FilterImportDuplicates)
			mergedContext, _ := gen.Clone()
			mergedContext.File = mergedAst

			result = append(result, &SourceCode{mergedContext, name})
		}
	}

	//If there are any results, an implementation was found and errors can be cleared
	if len(result) != 0 {
		errors = []error{}
		ok = true
	}

	return
}

func (imp *Implementor) relatedTypeName(t *ast.TypeSpec, imap ImplMap) string {

	var (
		implExpr    ast.Expr
		relatedName string = t.Name.Name
	)

	n := strings.Count(relatedName, "_")

	names := []ast.Expr{}

	ast.Inspect(t, func(node ast.Node) bool {
		if id, ok := node.(*ast.Ident); ok {
			if implExpr, found := imap[id.Name]; found {
				names = append(names, implExpr)
			}
		}
		return n != len(names)
	})

	for _, expr := range names {
		var implName string
		switch exprType := expr.(type) {
		case *ast.Ident:
			implName = exprType.Name

		default:
			implName = NiceName(implExpr)
		}
		relatedName = strings.Replace(relatedName, "_", implName, 1)
	}

	return relatedName
}

func NiceName(e ast.Expr) string {
	//TODO: Woefully inadaquate. Total failure for function types, interfaces, struct types

	// *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
	switch t := e.(type) {
	case *ast.Ident:
		return strings.Title(ExprString(e))

	case *ast.ParenExpr:
		return NiceName(t.X)

	case *ast.SelectorExpr:
		return NiceName(t.X)

	case *ast.StarExpr:
		return NiceName(t.X) + "Pointer"

	case *ast.ChanType:
		switch t.Dir {
		case ast.SEND:
			return NiceName(t.Value) + "SendChan"
		case ast.RECV:
			return NiceName(t.Value) + "RecvChan"
		default:
			return NiceName(t.Value) + "Chan"
		}

	case *ast.ArrayType:
		return NiceName(t.Elt) + "Slice"

	case *ast.MapType:
		return NiceName(t.Value) + "MapBy" + strings.Title(NiceName(t.Key))

	default:
		return strings.Title(ExprString(e))
	}
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
