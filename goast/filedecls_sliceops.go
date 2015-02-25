package main

import (
	"sort"
	"go/ast"
)

type fileDeclsSorter struct {
	fileDecls
	LessFunc	func(ast.Decl, ast.Decl) bool
}

func (s fileDeclsSorter) Less(i, j int) bool {
	return s.LessFunc(s.fileDecls[i], s.fileDecls[j])
}
func (s fileDecls) Len() int {
	return len(s)
}
func (s fileDecls) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s fileDecls) Sort(less func(ast.Decl, ast.Decl) bool) {
	sort.Sort(fileDeclsSorter{s, less})
}
func (s fileDecls) All(fn func(ast.Decl) bool) bool {
	for _, v := range s {
		if !fn(v) {
			return false
		}
	}
	return true
}
func (s fileDecls) Any(fn func(ast.Decl) bool) bool {
	for _, v := range s {
		if fn(v) {
			return true
		}
	}
	return false
}
func (s fileDecls) Count(fn func(ast.Decl) bool) int {
	count := 0
	for _, v := range s {
		if fn(v) {
			count += 1
		}
	}
	return count
}
func (s fileDecls) Each(fn func(ast.Decl)) {
	for _, v := range s {
		fn(v)
	}
}
func (s fileDecls) First(fn func(ast.Decl) bool) (match ast.Decl, found bool) {
	for _, v := range s {
		if fn(v) {
			match = v
			found = true
			break
		}
	}
	return
}
func (s fileDecls) Where(fn func(ast.Decl) bool) (result fileDecls) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
func (s *fileDecls) Extract(fn func(ast.Decl) bool) (removed fileDecls) {
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
