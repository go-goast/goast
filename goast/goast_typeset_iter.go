package main

import "go/ast"

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
func (s *typeSet) Extract(fn func(*ast.TypeSpec) bool) (removed typeSet) {
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
func (s typeSet) Fold(initial *ast.TypeSpec, fn func(*ast.TypeSpec, *ast.TypeSpec) *ast.TypeSpec) *ast.TypeSpec {
	folded := initial
	for _, v := range s {
		folded = fn(folded, v)
	}
	return folded
}
func (s typeSet) FoldR(initial *ast.TypeSpec, fn func(*ast.TypeSpec, *ast.TypeSpec) *ast.TypeSpec) *ast.TypeSpec {
	folded := initial
	for i := len(s) - 1; i >= 0; i-- {
		folded = fn(folded, s[i])
	}
	return folded
}
func (s typeSet) Where(fn func(*ast.TypeSpec) bool) (result typeSet) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
func (s typeSet) Zip(in ...typeSet) (result []typeSet) {
	minLen := len(s)
	for _, x := range in {
		if len(x) < minLen {
			minLen = len(x)
		}
	}
	for i := 0; i < minLen; i++ {
		row := typeSet{s[i]}
		for _, x := range in {
			row = append(row, x[i])
		}
		result = append(result, row)
	}
	return
}
