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
func (s typeSet) All(fn func(*ast.TypeSpec) bool) bool {
	for _, v := range s {
		if !fn(v) {
			return false
		}
	}
	return true
}
func (s typeSet) Any(fn func(*ast.TypeSpec) bool) bool {
	for _, v := range s {
		if fn(v) {
			return true
		}
	}
	return false
}
func (s typeSet) Count(fn func(*ast.TypeSpec) bool) int {
	count := 0
	for _, v := range s {
		if fn(v) {
			count += 1
		}
	}
	return count
}
func (s typeSet) Each(fn func(*ast.TypeSpec)) {
	for _, v := range s {
		fn(v)
	}
}
func (s typeSet) First(fn func(*ast.TypeSpec) bool) (match *ast.TypeSpec, found bool) {
	for _, v := range s {
		if fn(v) {
			match = v
			found = true
			break
		}
	}
	return
}
func (s typeSet) Where(fn func(*ast.TypeSpec) bool) (result typeSet) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
func (s *typeSet) Extract(fn func(*ast.TypeSpec) bool) (removed typeSet) {
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
