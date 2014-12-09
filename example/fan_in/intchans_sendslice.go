package main

func (s IntChans) Send(done <-chan struct{}) <-chan (<-chan int) {
	out := make(chan (<-chan int))
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
