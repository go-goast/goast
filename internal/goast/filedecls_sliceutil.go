package goast

import "go/ast"

func (s fileDecls) Len() int {
	return len(s)
}
func (s fileDecls) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
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
