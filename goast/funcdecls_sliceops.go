package main

import (
	"sort"
	"go/ast"
)

type funcDeclsSorter struct {
	funcDecls
	LessFunc	func(*ast.FuncDecl, *ast.FuncDecl) bool
}

func (s funcDeclsSorter) Less(i, j int) bool {
	return s.LessFunc(s.funcDecls[i], s.funcDecls[j])
}
func (s funcDecls) Len() int {
	return len(s)
}
func (s funcDecls) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s funcDecls) Sort(less func(*ast.FuncDecl, *ast.FuncDecl) bool) {
	sort.Sort(funcDeclsSorter{s, less})
}
func (s funcDecls) All(fn func(*ast.FuncDecl) bool) bool {
	for _, v := range s {
		if !fn(v) {
			return false
		}
	}
	return true
}
func (s funcDecls) Any(fn func(*ast.FuncDecl) bool) bool {
	for _, v := range s {
		if fn(v) {
			return true
		}
	}
	return false
}
func (s funcDecls) Count(fn func(*ast.FuncDecl) bool) int {
	count := 0
	for _, v := range s {
		if fn(v) {
			count += 1
		}
	}
	return count
}
func (s funcDecls) Each(fn func(*ast.FuncDecl)) {
	for _, v := range s {
		fn(v)
	}
}
func (s funcDecls) First(fn func(*ast.FuncDecl) bool) (match *ast.FuncDecl, found bool) {
	for _, v := range s {
		if fn(v) {
			match = v
			found = true
			break
		}
	}
	return
}
func (s funcDecls) Where(fn func(*ast.FuncDecl) bool) (result funcDecls) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
func (s *funcDecls) Extract(fn func(*ast.FuncDecl) bool) (removed funcDecls) {
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
