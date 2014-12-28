package main

//go:generate goast write impl jsonstringer.go

type User struct {
	id    int
	Name  string
	Email []string
}

func main() {
	u := User{3, "Katie", []string{"katie@example.com"}}
	println(u.Json())
}
