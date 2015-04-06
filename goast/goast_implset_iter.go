package main

func (s implSet) All(fn func(ImplMap) bool) bool {
	for _, v := range s {
		if !fn(v) {
			return false
		}
	}
	return true
}
func (s implSet) Any(fn func(ImplMap) bool) bool {
	for _, v := range s {
		if fn(v) {
			return true
		}
	}
	return false
}
func (s implSet) Count(fn func(ImplMap) bool) int {
	count := 0
	for _, v := range s {
		if fn(v) {
			count += 1
		}
	}
	return count
}
func (s implSet) Each(fn func(ImplMap)) {
	for _, v := range s {
		fn(v)
	}
}
func (s *implSet) Extract(fn func(ImplMap) bool) (removed implSet) {
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
func (s implSet) First(fn func(ImplMap) bool) (match ImplMap, found bool) {
	for _, v := range s {
		if fn(v) {
			match = v
			found = true
			break
		}
	}
	return
}
func (s implSet) Fold(initial ImplMap, fn func(ImplMap, ImplMap) ImplMap) ImplMap {
	folded := initial
	for _, v := range s {
		folded = fn(folded, v)
	}
	return folded
}
func (s implSet) FoldR(initial ImplMap, fn func(ImplMap, ImplMap) ImplMap) ImplMap {
	folded := initial
	for i := len(s) - 1; i >= 0; i-- {
		folded = fn(folded, s[i])
	}
	return folded
}
func (s implSet) Where(fn func(ImplMap) bool) (result implSet) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
func (s implSet) Zip(in ...implSet) (result []implSet) {
	minLen := len(s)
	for _, x := range in {
		if len(x) < minLen {
			minLen = len(x)
		}
	}
	for i := 0; i < minLen; i++ {
		row := implSet{s[i]}
		for _, x := range in {
			row = append(row, x[i])
		}
		result = append(result, row)
	}
	return
}
