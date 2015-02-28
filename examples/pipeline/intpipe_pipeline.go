package main

import (
	"sync"
)

type IntPipeFan []IntPipe

func (pip IntPipe) Collect(done <-chan bool, num int) (result []int) {
	for i := 0; i < num; i++ {
		select {
		case val := <-pip:
			result = append(result, val)
		case <-done:
			return
		}
	}
	return
}
func (pip IntPipe) Fan(done <-chan bool, workers int, fn func(int) int) IntPipe {
	return pip.FanOut(done, workers, fn).FanIn(done)
}
func (pip IntPipe) FanOut(done <-chan bool, workers int, fn func(int) int) IntPipeFan {
	fan := IntPipeFan{}
	for i := 0; i < workers; i++ {
		fan = append(fan, pip.worker(done, fn))
	}
	return fan
}
func (pip IntPipe) Filter(done <-chan bool, fn func(int) bool) IntPipe {
	out := make(chan int)
	go func() {
		defer close(out)
		for val := range pip {
			if fn(val) {
				select {
				case out <- val:
				case <-done:
					return
				}
			}
		}
	}()
	return out
}
func (pip IntPipe) Pipe(done <-chan bool, fn func(int) int) IntPipe {
	out := make(chan int)
	go func() {
		defer close(out)
		for val := range pip {
			select {
			case out <- fn(val):
			case <-done:
				return
			}
		}
	}()
	return out
}
func (pip IntPipe) worker(done <-chan bool, fn func(int) int) IntPipe {
	out := make(chan int)
	go func() {
		defer close(out)
		for val := range pip {
			select {
			case out <- fn(val):
			case <-done:
				return
			}
		}
	}()
	return out
}
func (fan IntPipeFan) FanIn(done <-chan bool) IntPipe {
	var wg sync.WaitGroup
	out := make(chan int)
	output := func(pl IntPipe) {
		defer wg.Done()
		for val := range pl {
			select {
			case out <- val:
			case <-done:
				return
			}
		}
	}
	wg.Add(len(fan))
	for _, val := range fan {
		go output(val)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
func (fan IntPipeFan) Filter(done <-chan bool, fn func(int) bool) (result IntPipeFan) {
	for _, pipe := range fan {
		result = append(result, pipe.Filter(done, fn))
	}
	return
}
func (fan IntPipeFan) Pipe(done <-chan bool, fn func(int) int) (result IntPipeFan) {
	for _, pipe := range fan {
		result = append(result, pipe.Pipe(done, fn))
	}
	return
}
