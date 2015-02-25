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
	"go/ast"
	"math"
)

//Complexity is a measure of 'how deep' an ast node goes.
// type I interface{}; Complexity == 1
// type Collection []I; Comlexity == 2 (1 for the array level, and 1 for the type it implements)
func (c *Context) Complexity(node *ast.TypeSpec) int {
	return c.complexityOfExpr(node.Type)
}

func (c *Context) complexityOfExpr(node ast.Expr) int {
	// *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
	switch nodeType := node.(type) {
	case *ast.Ident:
		if t, ok := c.LookupType(nodeType.Name); ok {
			return c.Complexity(t) + 1
		}
		return 1

	case *ast.Ellipsis:
		return c.complexityOfExpr(nodeType.Elt)

	case *ast.ParenExpr:
		return c.complexityOfExpr(nodeType.X)

	case *ast.SelectorExpr:
		return c.complexityOfExpr(nodeType.X)

	case *ast.StarExpr:
		return c.complexityOfExpr(nodeType.X)

	case *ast.ArrayType:
		return c.complexityOfExpr(nodeType.Elt) + 1

	case *ast.ChanType:
		return c.complexityOfExpr(nodeType.Value) + 1

	case *ast.FuncType:
		return intMax(
			c.complexityOfFieldList(nodeType.Params),
			c.complexityOfFieldList(nodeType.Results)) + 1

	case *ast.InterfaceType:
		return 1

	case *ast.MapType:
		return intMax(
			c.complexityOfExpr(nodeType.Key),
			c.complexityOfExpr(nodeType.Value)) + 1

	case *ast.StructType:
		return c.complexityOfFieldList(nodeType.Fields) + 1

	default:
		return 0
	}
}

func (c *Context) complexityOfFieldList(node *ast.FieldList) (max int) {
	if node == nil {
		return 0
	}

	for _, nodeField := range node.List {
		max = intMax(max, c.complexityOfExpr(nodeField.Type))
	}
	return max
}

func intMax(x, y int) int {
	return int(math.Max(float64(x), float64(y)))
}
