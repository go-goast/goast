package main

import (
	"fmt"
	"runtime"
)

//go:generate goast write impl gen/pipeline.go

type IntPipe <-chan int

func init() {
	runtime.GOMAXPROCS(8)
}

func main() {
	done := make(chan bool)
	defer close(done)

	result := rangePipe(done, 0, 1000).
		Fan(done, 10, sq). //Fan the squaring calculation out over 10 workers, send results back on a single pipeline
		Collect(done, 20)  //Collect the first 20 returned results

	for _, i := range result {
		fmt.Println(i)
	}
}

func sq(n int) int {
	return n * n
}

func rangePipe(done <-chan bool, min, max int) IntPipe {
	out := make(chan int)
	go (func() {
		for i := min; i <= max; i++ {
			select {
			case out <- i:
			case <-done:
			}
		}
	})()
	return out
}
