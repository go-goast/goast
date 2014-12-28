package main

type T struct {
	quit chan bool
}

func (t *T) Quit() {
	t.quit <- true
}
