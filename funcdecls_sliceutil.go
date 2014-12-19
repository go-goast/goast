package main

import "go/ast"

func (s funcDecls) Len() int {
	return len(s)
}
func (s funcDecls) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
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
