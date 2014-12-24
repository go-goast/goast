/*
Copyright 2014 James Garfield. All rights reserved.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package gen

import (
	"sort"
)

type I interface{}
type Slice []I

type _Sorter struct {
	Slice
	LessFunc func(I, I) bool
}

func (s _Sorter) Less(i, j int) bool {
	return s.LessFunc(s.Slice[i], s.Slice[j])
}

func (s Slice) Len() int {
	return len(s)
}

func (s Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Slice) All(fn func(I) bool) bool {
	for _, v := range s {
		if !fn(v) {
			return false
		}
	}
	return true
}

func (s Slice) Any(fn func(I) bool) bool {
	for _, v := range s {
		if fn(v) {
			return true
		}
	}
	return false
}

func (s Slice) Count(fn func(I) bool) int {
	count := 0
	for _, v := range s {
		if fn(v) {
			count += 1
		}
	}
	return count
}

func (s Slice) Each(fn func(I)) {
	for _, v := range s {
		fn(v)
	}
}

func (s Slice) First(fn func(I) bool) (match I, found bool) {
	for _, v := range s {
		if fn(v) {
			match = v
			found = true
			break
		}
	}
	return
}

func (s Slice) Sort(less func(I, I) bool) {
	sort.Sort(_Sorter{s, less})
}

func (s Slice) Where(fn func(I) bool) (result Slice) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
