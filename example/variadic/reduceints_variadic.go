package main

func (v ReduceInts) Bind(b int) (bound ReduceInts) {
	bound = func(arg ...int) int {
		args := append([]int{b}, arg...)
		return v(args...)
	}
	return
}
func (v ReduceInts) BindTail(b ...int) (bound ReduceInts) {
	bound = func(arg ...int) int {
		args := append(arg, b...)
		return v(args...)
	}
	return bound
}
