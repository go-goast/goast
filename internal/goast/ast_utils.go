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

package goast

import (
	"go/ast"
)

func declAsFuncDecl(d ast.Decl) (f *ast.FuncDecl, ok bool) {
	f, ok = d.(*ast.FuncDecl)
	return
}

func declAsTypeSpec(d ast.Decl) (t *ast.TypeSpec, ok bool) {
	if g, isGenDecl := d.(*ast.GenDecl); isGenDecl {
		t, ok = g.Specs[0].(*ast.TypeSpec)
		return t, ok
	}
	return
}

func isEmptyInterface(node ast.Node) bool {
	i, ok := node.(*ast.InterfaceType)
	if !ok {
		return false
	}

	isEmpty := i.Methods != nil && i.Methods.NumFields() == 0
	return isEmpty
}

func EquivalentExprs(a, b ast.Expr) bool {
	// *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
	switch aType := a.(type) {
	case *ast.Ident:
		if bType, ok := b.(*ast.Ident); ok {
			return aType.Name == bType.Name
		}

	case *ast.StarExpr:
		if bType, ok := b.(*ast.StarExpr); ok {
			return EquivalentExprs(aType.X, bType.X)
		}

	case *ast.ArrayType:
		if bType, ok := b.(*ast.ArrayType); ok {
			return EquivalentExprs(aType.Elt, bType.Elt)
		}

	case *ast.ChanType:
		if bType, ok := b.(*ast.ChanType); ok {
			return EquivalentExprs(aType.Value, bType.Value) && aType.Dir == bType.Dir
		}

	case *ast.FuncType:
		if bType, ok := b.(*ast.FuncType); ok {
			return equivalentFieldList(aType.Params, bType.Params) && equivalentFieldList(aType.Results, bType.Results)
		}

	case *ast.InterfaceType:
		if _, ok := b.(*ast.InterfaceType); ok {
			return true
		}

	case *ast.MapType:
		if bType, ok := b.(*ast.MapType); ok {
			return EquivalentExprs(aType.Key, bType.Key) && EquivalentExprs(aType.Value, bType.Value)
		}

	case *ast.StructType:
		if bType, ok := b.(*ast.StructType); ok {
			return equivalentFieldList(aType.Fields, bType.Fields)
		}
	}
	return false
}

func equivalentFieldList(a, b *ast.FieldList) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a.List) != len(b.List) {
		return false
	}

	for i, aField := range a.List {
		bField := b.List[i]
		if EquivalentExprs(aField.Type, bField.Type) {
			return false
		}
	}

	return true
}
