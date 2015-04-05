package main

func (s Vectors) All(fn func(Vector) bool) bool {
	for _, v := range s {
		if !fn(v) {
			return false
		}
	}
	return true
}
func (s Vectors) Any(fn func(Vector) bool) bool {
	for _, v := range s {
		if fn(v) {
			return true
		}
	}
	return false
}
func (s Vectors) Count(fn func(Vector) bool) int {
	count := 0
	for _, v := range s {
		if fn(v) {
			count += 1
		}
	}
	return count
}
func (s Vectors) Each(fn func(Vector)) {
	for _, v := range s {
		fn(v)
	}
}
func (s *Vectors) Extract(fn func(Vector) bool) (removed Vectors) {
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
func (s Vectors) First(fn func(Vector) bool) (match Vector, found bool) {
	for _, v := range s {
		if fn(v) {
			match = v
			found = true
			break
		}
	}
	return
}
func (s Vectors) Fold(initial Vector, fn func(Vector, Vector) Vector) Vector {
	folded := initial
	for _, v := range s {
		folded = fn(folded, v)
	}
	return folded
}
func (s Vectors) FoldR(initial Vector, fn func(Vector, Vector) Vector) Vector {
	folded := initial
	for i := len(s) - 1; i >= 0; i-- {
		folded = fn(folded, s[i])
	}
	return folded
}
func (s Vectors) Where(fn func(Vector) bool) (result Vectors) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
func (s Vectors) Zip(in ...Vectors) (result []Vectors) {
	minLen := len(s)
	for _, x := range in {
		if len(x) < minLen {
			minLen = len(x)
		}
	}
	for i := 0; i < minLen; i++ {
		row := Vectors{s[i]}
		for _, x := range in {
			row = append(row, x[i])
		}
		result = append(result, row)
	}
	return
}
