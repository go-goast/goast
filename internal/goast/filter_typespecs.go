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

//Filters out top level type declarations out of an ast
//ast.FilterFile cannot be used because it filters out import statements
//This is a known issue and the suggestion is to roll your own
//See: https://github.com/golang/go/issues/9248
func filterTypeSpecs(file *ast.File, fn func(string) bool) {
	i := 0
	for _, d := range file.Decls {
		if filterNodeForTypeSpec(d, fn) {
			file.Decls[i] = d
			i++
		}
	}
	file.Decls = file.Decls[0:i]
}

func filterNodeForTypeSpec(node ast.Node, fn func(string) bool) bool {
	switch t := node.(type) {
	case *ast.GenDecl:
		t.Specs = filterSpecsForTypeSpec(t.Specs, fn)
		return len(t.Specs) > 0

	case *ast.TypeSpec:
		return fn(t.Name.Name)

	default:
		return true
	}
}

func filterSpecsForTypeSpec(specs []ast.Spec, fn func(string) bool) []ast.Spec {
	i := 0
	for _, s := range specs {
		if filterNodeForTypeSpec(s, fn) {
			specs[i] = s
			i++
		}
	}
	specs = specs[0:i]
	return specs
}
