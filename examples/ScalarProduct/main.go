package main

import (
	"fmt"
)

//Type detection in impl finds all valid type targets
//goast.net/x/iter looks for a type that matchs []interface{}
//Both Vector and Vectors match this, so both will have `iter` implemented
//go:generate goast write impl goast.net/x/iter

type Vector []int64
type Vectors []Vector

func (v Vector) Product() int64 {
	return v.Fold(1, func(a, b int64) int64 { return a * b })
}

func (v Vector) Sum() int64 {
	return v.Fold(0, func(a, b int64) int64 { return a + b })
}

func (v Vector) ScalarProduct(vex ...Vector) int64 {
	var zipped Vectors = v.Zip(vex...)
	return zipped.Fold(Vector{}, func(acc, vec Vector) Vector {
		return append(acc, vec.Product())
	}).Sum()
}

func main() {
	x := Vector{1, 2, 3, 4, 5}
	y := Vector{10, 10, 10, 10, 10}

	sp := x.ScalarProduct(y)
	//(1*10) + (2*10) + (3*10) + (4*10) + (5*10) == 150
	fmt.Printf("Scalar Product: %d\n", sp)
}
