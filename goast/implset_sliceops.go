package main

import (
	"sort"
)

type implSetSorter struct {
	implSet
	LessFunc	func(ImplMap, ImplMap) bool
}

func (s implSetSorter) Less(i, j int) bool {
	return s.LessFunc(s.implSet[i], s.implSet[j])
}
func (s implSet) Len() int {
	return len(s)
}
func (s implSet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s implSet) Sort(less func(ImplMap, ImplMap) bool) {
	sort.Sort(implSetSorter{s, less})
}
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
func (s implSet) Where(fn func(ImplMap) bool) (result implSet) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
func (s *implSet) Extract(fn func(ImplMap) bool) (removed implSet) {
	pos := 0
	kept := *s
	for i := 0; i < kept.Len(); i++ {
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
