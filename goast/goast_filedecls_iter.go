package main

import "go/ast"

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
func (s *fileDecls) Extract(fn func(ast.Decl) bool) (removed fileDecls) {
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
func (s fileDecls) Fold(initial ast.Decl, fn func(ast.Decl, ast.Decl) ast.Decl) ast.Decl {
	folded := initial
	for _, v := range s {
		folded = fn(folded, v)
	}
	return folded
}
func (s fileDecls) FoldR(initial ast.Decl, fn func(ast.Decl, ast.Decl) ast.Decl) ast.Decl {
	folded := initial
	for i := len(s) - 1; i >= 0; i-- {
		folded = fn(folded, s[i])
	}
	return folded
}
func (s fileDecls) Where(fn func(ast.Decl) bool) (result fileDecls) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
func (s fileDecls) Zip(in ...fileDecls) (result []fileDecls) {
	minLen := len(s)
	for _, x := range in {
		if len(x) < minLen {
			minLen = len(x)
		}
	}
	for i := 0; i < minLen; i++ {
		row := fileDecls{s[i]}
		for _, x := range in {
			row = append(row, x[i])
		}
		result = append(result, row)
	}
	return
}
