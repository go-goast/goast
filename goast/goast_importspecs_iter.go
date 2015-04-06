package main

import "go/ast"

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
func (s *ImportSpecs) Extract(fn func(*ast.ImportSpec) bool) (removed ImportSpecs) {
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
func (s ImportSpecs) Fold(initial *ast.ImportSpec, fn func(*ast.ImportSpec, *ast.ImportSpec) *ast.ImportSpec) *ast.ImportSpec {
	folded := initial
	for _, v := range s {
		folded = fn(folded, v)
	}
	return folded
}
func (s ImportSpecs) FoldR(initial *ast.ImportSpec, fn func(*ast.ImportSpec, *ast.ImportSpec) *ast.ImportSpec) *ast.ImportSpec {
	folded := initial
	for i := len(s) - 1; i >= 0; i-- {
		folded = fn(folded, s[i])
	}
	return folded
}
func (s ImportSpecs) Where(fn func(*ast.ImportSpec) bool) (result ImportSpecs) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
func (s ImportSpecs) Zip(in ...ImportSpecs) (result []ImportSpecs) {
	minLen := len(s)
	for _, x := range in {
		if len(x) < minLen {
			minLen = len(x)
		}
	}
	for i := 0; i < minLen; i++ {
		row := ImportSpecs{s[i]}
		for _, x := range in {
			row = append(row, x[i])
		}
		result = append(result, row)
	}
	return
}
