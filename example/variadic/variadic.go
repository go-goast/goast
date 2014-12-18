package main

type T interface{}
type Variadic func(...T) T

func (v Variadic) Bind(b T) (bound Variadic) {
	bound = func(arg ...T) T {
		args := append([]T{b}, arg...)
		return v(args...)
	}
	return
}

func (v Variadic) BindTail(b ...T) (bound Variadic) {
	bound = func(arg ...T) T {
		args := append(arg, b...)
		return v(args...)
	}
	return
}
