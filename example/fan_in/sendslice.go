package main

type I interface{}
type Slice []I

func (s Slice) Send(done <-chan struct{}) <-chan I {
	out := make(chan I)
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
