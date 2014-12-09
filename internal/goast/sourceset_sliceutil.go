package goast

func (s SourceSet) Len() int {
	return len(s)
}
func (s SourceSet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
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
func (s SourceSet) Where(fn func(*SourceCode) bool) (result SourceSet) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
