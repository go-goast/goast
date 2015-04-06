package main

import (
	"sort"
	"go/ast"
)

type typeSetSorter struct {
	typeSet
	LessFunc	func(*ast.TypeSpec, *ast.TypeSpec) bool
}

func (s typeSetSorter) Less(i, j int) bool {
	return s.LessFunc(s.typeSet[i], s.typeSet[j])
}
func (s typeSet) Len() int {
	return len(s)
}
func (s typeSet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s typeSet) Sort(less func(*ast.TypeSpec, *ast.TypeSpec) bool) {
	sort.Sort(typeSetSorter{s, less})
}
