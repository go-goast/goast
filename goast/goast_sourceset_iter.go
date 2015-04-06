package main

func (s SourceSet) All(fn func(*SourceCode) bool) bool {
	for _, v := range s {
		if !fn(v) {
			return false
		}
	}
	return true
}
func (s SourceSet) Any(fn func(*SourceCode) bool) bool {
	for _, v := range s {
		if fn(v) {
			return true
		}
	}
	return false
}
func (s SourceSet) Count(fn func(*SourceCode) bool) int {
	count := 0
	for _, v := range s {
		if fn(v) {
			count += 1
		}
	}
	return count
}
func (s SourceSet) Each(fn func(*SourceCode)) {
	for _, v := range s {
		fn(v)
	}
}
func (s *SourceSet) Extract(fn func(*SourceCode) bool) (removed SourceSet) {
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
func (s SourceSet) First(fn func(*SourceCode) bool) (match *SourceCode, found bool) {
	for _, v := range s {
		if fn(v) {
			match = v
			found = true
			break
		}
	}
	return
}
func (s SourceSet) Fold(initial *SourceCode, fn func(*SourceCode, *SourceCode) *SourceCode) *SourceCode {
	folded := initial
	for _, v := range s {
		folded = fn(folded, v)
	}
	return folded
}
func (s SourceSet) FoldR(initial *SourceCode, fn func(*SourceCode, *SourceCode) *SourceCode) *SourceCode {
	folded := initial
	for i := len(s) - 1; i >= 0; i-- {
		folded = fn(folded, s[i])
	}
	return folded
}
func (s SourceSet) Where(fn func(*SourceCode) bool) (result SourceSet) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
func (s SourceSet) Zip(in ...SourceSet) (result []SourceSet) {
	minLen := len(s)
	for _, x := range in {
		if len(x) < minLen {
			minLen = len(x)
		}
	}
	for i := 0; i < minLen; i++ {
		row := SourceSet{s[i]}
		for _, x := range in {
			row = append(row, x[i])
		}
		result = append(result, row)
	}
	return
}
