package main

import "go/ast"

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
func (s *funcDecls) Extract(fn func(*ast.FuncDecl) bool) (removed funcDecls) {
	pos := 0
	kept := *s
	for i := 0; i < len(kept); i++ {
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
func (s funcDecls) Fold(initial *ast.FuncDecl, fn func(*ast.FuncDecl, *ast.FuncDecl) *ast.FuncDecl) *ast.FuncDecl {
	folded := initial
	for _, v := range s {
		folded = fn(folded, v)
	}
	return folded
}
func (s funcDecls) FoldR(initial *ast.FuncDecl, fn func(*ast.FuncDecl, *ast.FuncDecl) *ast.FuncDecl) *ast.FuncDecl {
	folded := initial
	for i := len(s) - 1; i >= 0; i-- {
		folded = fn(folded, s[i])
	}
	return folded
}
func (s funcDecls) Where(fn func(*ast.FuncDecl) bool) (result funcDecls) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
func (s funcDecls) Zip(in ...funcDecls) (result []funcDecls) {
	minLen := len(s)
	for _, x := range in {
		if len(x) < minLen {
			minLen = len(x)
		}
	}
	for i := 0; i < minLen; i++ {
		row := funcDecls{s[i]}
		for _, x := range in {
			row = append(row, x[i])
		}
		result = append(result, row)
	}
	return
}
