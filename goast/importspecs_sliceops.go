package main

import (
	"sort"
	"go/ast"
)

type ImportSpecsSorter struct {
	ImportSpecs
	LessFunc	func(*ast.ImportSpec, *ast.ImportSpec) bool
}

func (s ImportSpecsSorter) Less(i, j int) bool {
	return s.LessFunc(s.ImportSpecs[i], s.ImportSpecs[j])
}
func (s ImportSpecs) Len() int {
	return len(s)
}
func (s ImportSpecs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ImportSpecs) Sort(less func(*ast.ImportSpec, *ast.ImportSpec) bool) {
	sort.Sort(ImportSpecsSorter{s, less})
}
func (s ImportSpecs) All(fn func(*ast.ImportSpec) bool) bool {
	for _, v := range s {
		if !fn(v) {
			return false
		}
	}
	return true
}
func (s ImportSpecs) Any(fn func(*ast.ImportSpec) bool) bool {
	for _, v := range s {
		if fn(v) {
			return true
		}
	}
	return false
}
func (s ImportSpecs) Count(fn func(*ast.ImportSpec) bool) int {
	count := 0
	for _, v := range s {
		if fn(v) {
			count += 1
		}
	}
	return count
}
func (s ImportSpecs) Each(fn func(*ast.ImportSpec)) {
	for _, v := range s {
		fn(v)
	}
}
func (s ImportSpecs) First(fn func(*ast.ImportSpec) bool) (match *ast.ImportSpec, found bool) {
	for _, v := range s {
		if fn(v) {
			match = v
			found = true
			break
		}
	}
	return
}
func (s ImportSpecs) Where(fn func(*ast.ImportSpec) bool) (result ImportSpecs) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
func (s *ImportSpecs) Extract(fn func(*ast.ImportSpec) bool) (removed ImportSpecs) {
	pos := 0
	kept := *s
	for i := 0; i < kept.Len(); i++ {
		if fn(kept[i]) {
			removed = append(removed, kept[i])
		} else {
			kept[pos] = kept[i]
			pos++
		}
	}
	kept = kept[:pos:pos]
	*s = kept
	return removed
}
