package main

//go:generate goast write impl quittable.go

type StringProcess struct {
	in   <-chan string
	out  chan<- string
	quit chan bool
}

func main() {

}
