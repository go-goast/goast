package main

import (
	"sync"
)

func (f IntChans) FanIn(done <-chan struct{}) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)
	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}
	wg.Add(len(f))
	for _, c := range f {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
