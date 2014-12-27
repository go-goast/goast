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
	"go/token"
	"golang.org/x/tools/astutil"
	"sort"
	"strconv"
	"strings"
)

//go:generate goast write impl gen/sliceutil.go

//String identifer for the empty interface; used it Implementor.TypeMap
const (
	emptyInterface = "0"
)

type Implementor struct {
	TypeProvider *Context
	TypeMap      map[string]ast.Expr
	Generic      *Context
}

func NewImplementor(typeProvider *Context) *Implementor {
	imp := &Implementor{TypeProvider: typeProvider, TypeMap: make(map[string]ast.Expr)}
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

	//Test each type in the provider file for implementation
	for _, c := range candidateTypes {

		//reset local state
		//imp.Generic & imp.TypeMap get modified in-place during parts of this process,
		//so for for each iteration, you need a new copy
		if genClone, err := gen.Clone(); err != nil {
			errors = append(errors, err)
			return
		} else {
			imp.Generic = genClone
		}

		imp.TypeMap = make(map[string]ast.Expr)

		//If this candidate can't implement the primary generic type, there is no more to do
		//Save the error in case there is no implementation possible
		if ok, err := imp.implementType(primaryGeneric, c); !ok {
			errors = append(errors, err)
			continue
		}

		imp.storeTypeMapping(primaryGeneric.Name.Name, c.Name)

		//Types can only be matched once, so no need to keep the selected candidate
		specTypes := candidateTypes.Where(func(t *ast.TypeSpec) bool {
			return t.Name.Name != c.Name.Name
		})

		matched := 0
		for _, g := range implTypes {
			//types that already have mappings are already solved for
			if _, ok := imp.TypeMap[g.Name.Name]; ok {
				matched++
				continue
			}

			//check each specification type to see if it satisfies the requirements for this generic type
			for _, s := range specTypes {
				ok, err := imp.implementType(g, s)
				if ok {
					imp.storeTypeMapping(g.Name.Name, s.Name)
					matched++
					break
				} else {
					errors = append(errors, err)
				}
			}

			//being unable to satify a generic type indicates we can't implement with the current combination of types
			//Save an error in case there is no implementation possible
			if _, ok := imp.TypeMap[g.Name.Name]; !ok {
				errors = append(errors, fmt.Errorf("Unable to satisfy generic type %s with any of the specification types.", g.Name.Name))
				break
			}
		}

		//If all generic types are mapped, rewrite the generic AST with the provided types
		if matched == implTypes.Len() {

			//Filter implemented types out
			//Do this prior to renaming related types so we can still identify them
			//ast.FilterFile filters out import statements...always, so use custom filter method https://github.com/golang/go/issues/9248
			filterTypeSpecs(imp.Generic.File, func(t *ast.TypeSpec) bool { return !isImplType(t) })

			//Generate names for all related types
			relatedTypes.Each(func(t *ast.TypeSpec) {
				specName := imp.relatedTypeName(t)
				imp.storeTypeMapping(t.Name.Name, ast.NewIdent(specName))
			})

			ast.Walk(imp, imp.Generic.File)

			//ensure that implementation is in the correct package
			imp.Generic.SetPackage(imp.TypeProvider.File.Name.Name)

			name := c.Name.Name

			result = append(result, &SourceCode{imp.Generic, name})
		}
	}

	//If there are any results, an implementation was found and errors can be cleared
	if result.Len() != 0 {
		errors = []error{}
		ok = true
	}

	return
}

func (imp *Implementor) implementType(gen, spec *ast.TypeSpec) (bool, error) {
	return imp.implementExpr(gen.Type, spec.Type)
}

func (imp *Implementor) implementExpr(gen, spec ast.Expr) (ok bool, err error) {
	// *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
	switch genType := gen.(type) {
	case *ast.Ident:
		return imp.implementIdent(genType, spec)

	case *ast.Ellipsis:
		//TODO Augement matching to allow ident-slipthru
		if specType, ok := spec.(*ast.Ellipsis); ok {
			return imp.implementExpr(genType.Elt, specType.Elt)
		}

	case *ast.ParenExpr:
		return imp.implementExpr(genType.X, spec)

	case *ast.SelectorExpr:
		//TODO Augement matching to allow ident-slipthru
		if specType, ok := spec.(*ast.SelectorExpr); ok {
			return imp.implementExpr(genType.X, specType.X)
		}

	case *ast.StarExpr:
		if specType, ok := spec.(*ast.StarExpr); ok {
			return imp.implementExpr(genType.X, specType.X)
		}

	case *ast.ArrayType:
		if specType, ok := spec.(*ast.ArrayType); ok {
			return imp.implementExpr(genType.Elt, specType.Elt)
		}

	case *ast.ChanType:
		if specType, ok := spec.(*ast.ChanType); ok {
			if genType.Dir != specType.Dir && specType.Arrow != token.NoPos {
				err = fmt.Errorf("Cannot implement generic channel. Channel directions do not match. Expected: %s or bidirectional, Found: %s", ExprString(genType), ExprString(specType))
				return false, err
			}
			return imp.implementExpr(genType.Value, specType.Value)
		}

	case *ast.FuncType:
		if specType, ok := spec.(*ast.FuncType); ok {
			if ok, err := imp.implementFieldList(genType.Params, specType.Params); !ok {
				return false, fmt.Errorf("Cannot implement generic function because of a parameter list error: %s", err)
			}

			if ok, err := imp.implementFieldList(genType.Results, specType.Results); !ok {
				return false, fmt.Errorf("Cannot implement generic function because of a result list error: %s", err)
			}
			return true, nil
		}

	case *ast.InterfaceType:
		return imp.implementInterfaceType(genType, spec)

	case *ast.MapType:
		if specType, ok := spec.(*ast.MapType); ok {
			if ok, err := imp.implementExpr(genType.Key, specType.Key); !ok {
				return false, fmt.Errorf("Cannot implement generic map because of a key type error: %s", err)
			}
			if ok, err := imp.implementExpr(genType.Value, specType.Value); !ok {
				return false, fmt.Errorf("Cannot implement generic map because of a value type error: %s", err)
			}
			return true, nil
		}

	case *ast.StructType:
		if specType, ok := spec.(*ast.StructType); ok {
			return imp.implementStruct(genType, specType)
		}

	default:
		err = fmt.Errorf("Invalid Expression. %s", ExprString(gen))
		return
	}

	err = fmt.Errorf("Cannot implement generic expression %s with matching specification expression %s", ExprString(gen), ExprString(spec))
	return
}

func (imp *Implementor) implementIdent(gen *ast.Ident, spec ast.Expr) (ok bool, err error) {

	//Already mapped type, ensure that spec matches expected identifier
	if val, exists := imp.TypeMap[gen.Name]; exists {
		ok = EquivalentExprs(val, spec)
		if !ok {
			err = fmt.Errorf("Cannot implement identifier %s as %s, already mapped to %s", gen.Name, ExprString(spec), ExprString(val))
		}
		return
	}

	if genType, isType := imp.Generic.LookupType(gen.Name); isType {
		if ok = isEmptyInterface(genType.Type); ok {
			imp.storeTypeMapping(gen.Name, spec)
			return
		} else if ok, err = imp.implementExpr(genType.Type, spec); ok {
			imp.storeTypeMapping(gen.Name, spec)
			return
		}
		err = fmt.Errorf("Cannot implement %s with %s. Error: %s", gen.Name, ExprString(spec), err)
		return
	}

	specIdent, ok := spec.(*ast.Ident)
	if !ok {
		err = fmt.Errorf("Cannot implement ident %s with non-ident %+v", gen.Name, spec)
		return
	}

	if ok = (gen.Name == specIdent.Name); ok {
		imp.storeTypeMapping(gen.Name, spec)
		return
	}

	err = fmt.Errorf("Cannot implement %s with %s", gen.Name, specIdent.Name)
	return
}

func (imp *Implementor) implementInterfaceType(gen *ast.InterfaceType, spec ast.Expr) (ok bool, err error) {

	if ok = isEmptyInterface(gen); !ok {
		err = fmt.Errorf("Non-empty interface types are not supported in generic implementations\n%s", ExprString(gen))
		return
	}

	imp.storeTypeMapping(emptyInterface, spec)
	return
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

func (imp *Implementor) implementStruct(gen, spec *ast.StructType) (ok bool, err error) {

	genCount := gen.Fields.NumFields()
	//Empty generic structs match any other stuct
	if genCount == 0 {
		ok = true
		return
	}

	specCount := spec.Fields.NumFields()
	//Specification structs must have at least as many fields as generic stucts to be able to match
	if specCount < genCount {
		err = fmt.Errorf("Not enough fields to implement struct")
		return
	}

	//Check that the specification struct implements all fields in the generic struct
	//TODO How are embedded types handled?
	//TODO Is there a way to support _ named fields? What would this mean?
	for _, field := range gen.Fields.List {
		for _, name := range field.Names {
			var nameMatch *ast.Field
			nameMatch, ok = FieldByName(spec.Fields, name.Name)
			if !ok {
				err = fmt.Errorf("Missing field %s", name.Name)
				return
			}
			if ok, err = imp.implementExpr(field.Type, nameMatch.Type); !ok {
				return
			}
		}
	}
	ok = true
	return
}

func (imp *Implementor) implementFieldList(gen, spec *ast.FieldList) (ok bool, err error) {
	if gen == nil && spec == nil {
		ok = true
		return
	}

	if gen == nil {
		err = fmt.Errorf("Generic FieldList Empty")
		return
	}

	if spec == nil {
		err = fmt.Errorf("Specification FieldList Empty")
		return
	}

	if len(gen.List) != len(spec.List) {
		err = fmt.Errorf("FieldLists do not match in length")
		return
	}

	for i, genField := range gen.List {
		specField := spec.List[i]
		if ok, err = imp.implementExpr(genField.Type, specField.Type); !ok {
			return
		}
	}

	return
}

func (imp *Implementor) storeTypeMapping(genId string, spec ast.Expr) (ok bool, err error) {
	imp.implementSpecificationExpr(spec)
	println("Map ", genId, " -> ", ExprString(spec))
	imp.TypeMap[genId] = spec
	ok = true
	return
}

//Traverse a specification expression to ensure that any imported types within have their packages added to imports
func (imp *Implementor) implementSpecificationExpr(node ast.Expr) (ok bool, err error) {
	// *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
	switch t := node.(type) {
	case *ast.Ident:
		if specType, isType := imp.TypeProvider.LookupType(t.Name); isType {
			return imp.implementSpecificationExpr(specType.Type)
		}
	case *ast.ParenExpr:
		return imp.implementSpecificationExpr(t.X)

	case *ast.SelectorExpr:
		return imp.implementSpecificationSelector(t)

	case *ast.StarExpr:
		return imp.implementSpecificationExpr(t.X)

	case *ast.Ellipsis:
		return imp.implementSpecificationExpr(t.Elt)

	case *ast.ArrayType:
		return imp.implementSpecificationExpr(t.Elt)

	case *ast.ChanType:
		return imp.implementSpecificationExpr(t.Value)

	case *ast.FuncType:
		if ok, err = imp.implementSpecificationFieldList(t.Params); !ok {
			return
		}
		return imp.implementSpecificationFieldList(t.Results)

	case *ast.InterfaceType:
		return imp.implementSpecificationFieldList(t.Methods)

	case *ast.MapType:
		if ok, err = imp.implementSpecificationExpr(t.Key); !ok {
			return
		}
		return imp.implementSpecificationExpr(t.Value)
	}
	err = fmt.Errorf("Unimplementable specification expression %s", ExprString(node))
	return
}

func (imp *Implementor) implementSpecificationFieldList(list *ast.FieldList) (ok bool, err error) {
	for _, field := range list.List {
		if ok, err = imp.implementSpecificationExpr(field.Type); !ok {
			return
		}
	}
	return
}

func (imp *Implementor) implementSpecificationSelector(node *ast.SelectorExpr) (ok bool, err error) {
	id, isIdent := node.X.(*ast.Ident)
	if !isIdent {
		err = fmt.Errorf("Cannot implement selector expression %s because the selected expression is not an identifier: %s", node.Sel.Name, ExprString(node.X))
		return
	}

	ok = true

	//Check to see if a package identifier is being selected
	//If so, the implemented specification need to have that
	//package added to its imports
	if match, found := imp.TypeProvider.LookupImport(id.Name); found {
		newImport, _ := strconv.Unquote(match.Path.Value)
		astutil.AddImport(token.NewFileSet(), imp.Generic.File, newImport)
	}
	return
}

//Visits nodes in the generic ast to rewrite with specification types
func (imp *Implementor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch t := node.(type) {
	case *ast.ArrayType:
		return imp.visitArrayType(t)

	case *ast.ChanType:
		return imp.visitChanType(t)

	case *ast.Ellipsis:
		return imp.visitEllipsis(t)

	case *ast.Ident:
		return imp.visitIdent(t)

	case *ast.MapType:
		return imp.visitMapType(t)

	case *ast.ParenExpr:
		return imp.visitParenExpr(t)

	case *ast.StarExpr:
		return imp.visitStarExpr(t)

	case *ast.Field:
		return imp.visitField(t)

	default:
		return imp
	}
}

func (imp *Implementor) visitArrayType(node *ast.ArrayType) ast.Visitor {
	if t, ok := imp.replacementType(node.Elt); ok {
		node.Elt = t
		return nil
	}
	return imp
}

func (imp *Implementor) visitChanType(node *ast.ChanType) ast.Visitor {
	if t, ok := imp.replacementType(node.Value); ok {
		switch t.(type) {
		case *ast.ChanType:
			//it seems that the syntax of channels of channels gets a bit ambiguous after a while?
			//surrounding the specification type in parens seems to clear that up
			node.Value = &ast.ParenExpr{X: t}
		default:
			node.Value = t
		}

		return nil
	}
	return imp
}

func (imp *Implementor) visitEllipsis(node *ast.Ellipsis) ast.Visitor {
	if t, ok := imp.replacementType(node.Elt); ok {
		node.Elt = t
		return nil
	}
	return imp
}

func (imp *Implementor) visitIdent(node *ast.Ident) ast.Visitor {
	if t, ok := imp.replacementType(node); ok {
		if id, ok := t.(*ast.Ident); ok {
			node.Name = id.Name
		} else {
			fmt.Printf("Invalid ident replacement for %s: %s", node.Name, ExprString(t))
		}
		return nil
	}
	return imp
}

func (imp *Implementor) visitMapType(node *ast.MapType) ast.Visitor {
	if t, ok := imp.replacementType(node.Key); ok {
		node.Key = t
		return nil
	}

	if t, ok := imp.replacementType(node.Value); ok {
		node.Value = t
		return nil
	}

	return imp
}

func (imp *Implementor) visitParenExpr(node *ast.ParenExpr) ast.Visitor {
	if t, ok := imp.replacementType(node.X); ok {
		node.X = t
		return nil
	}
	return imp
}

func (imp *Implementor) visitStarExpr(node *ast.StarExpr) ast.Visitor {
	if t, ok := imp.replacementType(node.X); ok {
		node.X = t
		return nil
	}
	return imp
}

func (imp *Implementor) visitField(node *ast.Field) ast.Visitor {
	if t, ok := imp.replacementType(node.Type); ok {
		node.Type = t
		return nil
	}
	return imp
}

//Determines what if anything a given ast node should be replaced with
func (imp *Implementor) replacementType(node ast.Node) (ast.Expr, bool) {
	switch t := node.(type) {
	case *ast.Ident:
		if val, ok := imp.TypeMap[t.Name]; ok {
			return val, true
		}
		return nil, false

	case *ast.InterfaceType:
		if isEmptyInterface(node) {
			return imp.TypeMap[emptyInterface], true
		}
		return nil, false

	default:
		return nil, false
	}
}

func (imp *Implementor) relatedTypeName(t *ast.TypeSpec) string {

	var (
		implExpr    ast.Expr
		found       bool
		implName    string
		partialName string = t.Name.Name
	)

	ast.Inspect(t, func(node ast.Node) bool {
		if id, ok := node.(*ast.Ident); ok {
			implExpr, found = imp.TypeMap[id.Name]
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
