package main

import (
	"fmt"
	"runtime"
)

//go:generate goast write impl fanin.go
//go:generate goast write impl sendslice.go

type IntChans []<-chan int
type Ints []int

func main() {
	runtime.GOMAXPROCS(8)

	done := make(chan struct{})
	defer close(done)

	in := intRange(100).Send(done)

	workers := IntChans{
		sq(done, in),
		sq(done, in),
		sq(done, in),
		sq(done, in),
		sq(done, in),
		sq(done, in),
		sq(done, in),
		sq(done, in),
		sq(done, in),
		sq(done, in),
		sq(done, in)}

	out := workers.FanIn(done)
	for range intRange(13) {
		fmt.Println(<-out)
	}

	done <- struct{}{}
}

func sq(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func intRange(max int) (r Ints) {
	for i := 0; i < max; i++ {
		r = append(r, i)
	}
	return
}
