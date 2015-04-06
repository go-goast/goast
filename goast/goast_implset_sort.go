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
