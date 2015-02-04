package main

import "go/ast"

func (mp ImplMap) Copy() ImplMap {
	var newMap ImplMap = make(map[string]ast.Expr)
	for k, v := range mp {
		newMap[k] = v
	}
	return newMap
}
func (mp *ImplMap) Init() {
	var newMap ImplMap = make(map[string]ast.Expr)
	mp = &newMap
}
