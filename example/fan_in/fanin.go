package main

import (
	"sync"
)

type E interface{}
type Fan []<-chan E

func (f Fan) FanIn(done <-chan struct{}) <-chan E {
	var wg sync.WaitGroup
	out := make(chan E)

	output := func(c <-chan E) {
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
