package main

func (t *StringProcess) Quit() {
	t.quit <- true
}
