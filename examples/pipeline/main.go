package main

import (
	"fmt"
	"runtime"
)

//go:generate goast write impl goast.net/x/pipeline

type IntPipe chan int

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	done := make(chan bool)
	defer close(done)

	workers := runtime.NumCPU()

	var result Ints = rangePipe(done, 0, 1000).
		Fan(done, workers, sq). //Fan the squaring calculation out over 10 workers, send results back on a single pipeline
		Collect(done, 20)       //Collect the first 20 returned results

	for _, i := range result {
		fmt.Println(i)
	}
}

//Sends the provided range of ints out over the returned IntPipe
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

func sq(n int) int {
	return n * n
}
