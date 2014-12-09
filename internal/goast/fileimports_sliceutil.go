package goast

import "go/ast"

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
func (s fileImports) Where(fn func(*ast.ImportSpec) bool) (result fileImports) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
