package main

func (s Ints) Send(done <-chan struct{}) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, val := range s {
			select {
			case out <- val:
			case <-done:
				close(out)
				return
			}
		}
	}()
	return out
}
