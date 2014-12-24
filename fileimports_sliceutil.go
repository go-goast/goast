package main

import (
	"sort"
	"go/ast"
)

type fileImportsSorter struct {
	fileImports
	LessFunc	func(*ast.ImportSpec, *ast.ImportSpec) bool
}

func (s fileImportsSorter) Less(i, j int) bool {
	return s.LessFunc(s.fileImports[i], s.fileImports[j])
}
func (s fileImports) Len() int {
	return len(s)
}
func (s fileImports) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s fileImports) All(fn func(*ast.ImportSpec) bool) bool {
	for _, v := range s {
		if !fn(v) {
			return false
		}
	}
	return true
}
func (s fileImports) Any(fn func(*ast.ImportSpec) bool) bool {
	for _, v := range s {
		if fn(v) {
			return true
		}
	}
	return false
}
func (s fileImports) Count(fn func(*ast.ImportSpec) bool) int {
	count := 0
	for _, v := range s {
		if fn(v) {
			count += 1
		}
	}
	return count
}
func (s fileImports) Each(fn func(*ast.ImportSpec)) {
	for _, v := range s {
		fn(v)
	}
}
func (s fileImports) First(fn func(*ast.ImportSpec) bool) (match *ast.ImportSpec, found bool) {
	for _, v := range s {
		if fn(v) {
			match = v
			found = true
			break
		}
	}
	return
}
func (s fileImports) Sort(less func(*ast.ImportSpec, *ast.ImportSpec) bool) {
	sort.Sort(fileImportsSorter{s, less})
}
func (s fileImports) Where(fn func(*ast.ImportSpec) bool) (result fileImports) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
