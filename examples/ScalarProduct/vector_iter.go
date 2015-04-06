package main

func (s Vector) All(fn func(int64) bool) bool {
	for _, v := range s {
		if !fn(v) {
			return false
		}
	}
	return true
}
func (s Vector) Any(fn func(int64) bool) bool {
	for _, v := range s {
		if fn(v) {
			return true
		}
	}
	return false
}
func (s Vector) Count(fn func(int64) bool) int {
	count := 0
	for _, v := range s {
		if fn(v) {
			count += 1
		}
	}
	return count
}
func (s Vector) Each(fn func(int64)) {
	for _, v := range s {
		fn(v)
	}
}
func (s *Vector) Extract(fn func(int64) bool) (removed Vector) {
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
func (s Vector) First(fn func(int64) bool) (match int64, found bool) {
	for _, v := range s {
		if fn(v) {
			match = v
			found = true
			break
		}
	}
	return
}
func (s Vector) Fold(initial int64, fn func(int64, int64) int64) int64 {
	folded := initial
	for _, v := range s {
		folded = fn(folded, v)
	}
	return folded
}
func (s Vector) FoldR(initial int64, fn func(int64, int64) int64) int64 {
	folded := initial
	for i := len(s) - 1; i >= 0; i-- {
		folded = fn(folded, s[i])
	}
	return folded
}
func (s Vector) Where(fn func(int64) bool) (result Vector) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
func (s Vector) Zip(in ...Vector) (result []Vector) {
	minLen := len(s)
	for _, x := range in {
		if len(x) < minLen {
			minLen = len(x)
		}
	}
	for i := 0; i < minLen; i++ {
		row := Vector{s[i]}
		for _, x := range in {
			row = append(row, x[i])
		}
		result = append(result, row)
	}
	return
}
