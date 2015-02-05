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
)

func Implement(cp ContextPair, known ImplMap, spec, gen *ast.TypeSpec) (ok bool, result ImplMap, err error) {

	result = known.Copy()

	ok, err = implementType(cp, result, gen, spec)
	if ok {
		ok, err = result.Store(gen.Name.Name, spec.Name)
	}
	return
}

func implementType(cp ContextPair, known ImplMap, gen, spec *ast.TypeSpec) (bool, error) {
	return implementExpr(cp, known, gen.Type, spec.Type)
}

func implementExpr(cp ContextPair, known ImplMap, gen, spec ast.Expr) (ok bool, err error) {
	// *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
	switch genType := gen.(type) {
	case *ast.Ident:
		return implementIdent(cp, known, genType, spec)

	case *ast.Ellipsis:
		//TODO Augement matching to allow ident-slipthru
		if specType, ok := spec.(*ast.Ellipsis); ok {
			return implementExpr(cp, known, genType.Elt, specType.Elt)
		}

	case *ast.ParenExpr:
		return implementExpr(cp, known, genType.X, spec)

	case *ast.SelectorExpr:
		//TODO Augement matching to allow ident-slipthru
		if specType, ok := spec.(*ast.SelectorExpr); ok {
			return implementExpr(cp, known, genType.X, specType.X)
		}

	case *ast.StarExpr:
		if specType, ok := spec.(*ast.StarExpr); ok {
			return implementExpr(cp, known, genType.X, specType.X)
		}

	case *ast.ArrayType:
		if specType, ok := spec.(*ast.ArrayType); ok {
			return implementExpr(cp, known, genType.Elt, specType.Elt)
		}

	case *ast.ChanType:
		if specType, ok := spec.(*ast.ChanType); ok {
			if genType.Dir != specType.Dir && specType.Arrow != token.NoPos {
				err = fmt.Errorf("Cannot implement generic channel. Channel directions do not match. Expected: %s or bidirectional, Found: %s", ExprString(genType), ExprString(specType))
				return false, err
			}
			return implementExpr(cp, known, genType.Value, specType.Value)
		}

	case *ast.FuncType:
		if specType, ok := spec.(*ast.FuncType); ok {
			if ok, err := implementFieldList(cp, known, genType.Params, specType.Params); !ok {
				return false, fmt.Errorf("Cannot implement generic function because of a parameter list error: %s", err)
			}

			if ok, err := implementFieldList(cp, known, genType.Results, specType.Results); !ok {
				return false, fmt.Errorf("Cannot implement generic function because of a result list error: %s", err)
			}
			return true, nil
		}

	case *ast.InterfaceType:
		return implementInterfaceType(cp, known, genType, spec)

	case *ast.MapType:
		if specType, ok := spec.(*ast.MapType); ok {
			if ok, err := implementExpr(cp, known, genType.Key, specType.Key); !ok {
				return false, fmt.Errorf("Cannot implement generic map because of a key type error: %s", err)
			}
			if ok, err := implementExpr(cp, known, genType.Value, specType.Value); !ok {
				return false, fmt.Errorf("Cannot implement generic map because of a value type error: %s", err)
			}
			return true, nil
		}

	case *ast.StructType:
		if specType, ok := spec.(*ast.StructType); ok {
			return implementStruct(cp, known, genType, specType)
		}

	default:
		err = fmt.Errorf("Invalid Expression. %s", ExprString(gen))
		return
	}

	err = fmt.Errorf("Cannot implement generic expression %s with matching specification expression %s", ExprString(gen), ExprString(spec))
	return
}

func implementIdent(cp ContextPair, known ImplMap, gen *ast.Ident, spec ast.Expr) (ok bool, err error) {

	if genType, isType := cp.Generic.LookupType(gen.Name); isType {
		if ok = isEmptyInterface(genType.Type); ok {
			known.Store(gen.Name, spec)
			return
		} else if ok, err = implementExpr(cp, known, genType.Type, spec); ok {
			known.Store(gen.Name, spec)
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
		return
	}

	err = fmt.Errorf("Cannot implement %s with %s", gen.Name, specIdent.Name)
	return
}

func implementInterfaceType(cp ContextPair, known ImplMap, gen *ast.InterfaceType, spec ast.Expr) (ok bool, err error) {

	if ok = isEmptyInterface(gen); !ok {
		err = fmt.Errorf("Non-empty interface types are not supported in generic implementations\n%s", ExprString(gen))
		return
	}

	//ok, err = known.Store(emptyInterface, spec)
	return
}

func implementStruct(cp ContextPair, known ImplMap, gen, spec *ast.StructType) (ok bool, err error) {

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
			if ok, err = implementExpr(cp, known, field.Type, nameMatch.Type); !ok {
				return
			}
		}
	}
	ok = true
	return
}

func implementFieldList(cp ContextPair, known ImplMap, gen, spec *ast.FieldList) (ok bool, err error) {
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
		if ok, err = implementExpr(cp, known, genField.Type, specField.Type); !ok {
			return
		}
	}

	return
}
