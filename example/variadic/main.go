package main

//go:generate goast write impl variadic.go

type ReduceInts func(...int) int

func main() {

	var sumFn ReduceInts = Sum
	println(sumFn(1, 2, 3))

	sumPlus2 := sumFn.Bind(2)

	println(sumPlus2(1, 2, 3))
}

func Sum(ints ...int) (result int) {
	for _, i := range ints {
		result += i
	}
	return
}
