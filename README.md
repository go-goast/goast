goast
=====

goast is a Go AST (Abstract Syntax Tree) based meta-programming utility with the aim of providing safe, idiomatic code generation facilities by taking advantage of Go's native AST abilities and the go:generate build directive.

goast's core philosophies are

* Compile time type safety
* No runtime type casting (see previous)
* Avoid runtime reflection
* Prefer pure Go over syntax extensions
* No text templates (see previous)
* Prefer inference over annotation

The functionality of goast is currently built on the following axiom and proposition

1. The empty interface (`interface{}`) can be replaced with any other type
2. Any composite type composed at least partially of the empty interface (e.g. `map[string]interface{}`) can be replaced with any other composite type of the same structure with the empty interface swapped out for a concrete type (e.g. `map[string][]int64`)


## Installing

Installing goast is as easy installing any other go package; Use the `go get` command to install from the canonical import path:
```
go get goast.net/x/goast
```


## Usage

goast is designed as a general purpose command line utility for reading, examining, writing, and working with Go files and their AST. It follows a `cli command subcommand` pattern similar to git and go. 

This document will focus primarily on a single subcommand `goast write impl`, but you can get more information on any other commands by typing `goast help`.

*Please note: While functional, goast is still in an alpha/RFC phase so more subcommands may be added over time and subcommand structure may change (hopefully not)*


### Simple example

Consider the following generic implementation of a filtering method on a slice of empty interface


```go
//file: slicefilter.go
package gen

type T interface{}
type Slice []T
func (s Slice) Where(fn func(T)bool) (result Slice) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return
}
```

This code compiles, and accurately reflects the algorithm that any slice type might implement. There can be some initial confusion where a new gopher would expect the following to work

```go
package main

import "fmt"

type Slice []interface{}

func (s Slice) Where(fn func(interface{}) bool) (result Slice) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return
}

func main() {

	var s Slice = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	evens := s.Where(func(val interface{}) bool {
		i := val.(int)
		return i%2 == 0
	})

	fmt.Printf("Evens: %+v\n", evens)
}
```

Our intrepid gopher is a little sad about needing to use type casting, but is surprised to discover the following compile error `main.go:18: cannot use []int literal (type []int) as type Slice in assignment` ([Playground](http://play.golang.org/p/S5elysXjrx)). Surely a slice of ints provides the same interface as a slice of empty interface? The explanation for this is that an `[]int` cannot be assigned to an `[]interface{}` because they have different memory layouts. This is unfortunate because algorithmically the code is correct and it seems intuitive that it could be interpreted that way. 

It is possible to [rewrite](http://play.golang.org/p/IpvXBKSxpS) the above by using `var s Slice = []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9}`, but converting back and forth between slice types on top of already needing to do type casting is onerous and a great place for bugs to creep in.


`goast write impl` exploits the fact that you can write compilable, generic Go code in that manner to provide typesafe code generation of a specific implementation of that meta-code.

With goast, a developer only needs to provide the following code, and the go generate command will provide the rest

```go
//file main.go
package main

import "fmt"

//go:generate goast write impl slicefilter.go

type Ints []int

func main() {
	var s Ints = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	evens := s.Where(func(val int) bool { return val%2 == 0 })
	fmt.Printf("Evens: %+v\n", evens)
}


```

The `go:generate` build directive instructs goast to write an implementation of the code in slicefilter.go for the types provided in main.go, resulting in the following file being generated

```go
//file ints_slicefilter.go
package main

func (s Ints) Where(fn func(int)bool) (result Ints) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return
}
```
[Playground](http://play.golang.org/p/ATVduqekFh)

A reader following the code closely might have noticed that there is no longer a need for type casting within the callback function to `Where`. goast infers concrete type information from your source file and rewrites the the meta-code to be type-safe for your types, thus providing the ability to write correct, generic code once, and use it in an opt-in, type-safe manner.

A more complete slice iteration library can be seen in the [iter](http://goast.net/x/iter) package.

**A Note About Paths:** The examples above use relative paths to implment the code in a single, local file. goast also supports fetching code from the `$GOPATH`. The above example could be rewritten to use the  `iter` package like so:

```go
package main

import "fmt"

//go:generate goast write impl goast.net/x/iter

type Ints []int

func main() {
	var s Ints = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	evens := s.Where(func(val int) bool { return i%2 == 0 })
	fmt.Printf("Evens: %+v\n", evens)
}
```

### Complex Example

Vanilla iteration patterns aren't exciting enough for you? 

Want something more Go-centric, maybe even concurrency or parallel oriented? 

How about a Fan-Out/Fan-In Concurrent Pipeline?

If you're not familiar with pipelines, [this](http://blog.golang.org/pipelines) is a great primer on the pattern and pitfalls of implementing it. The concept is strightforward, but it's not the kind of thing I can see trusting myself to write repeatedly.

With goast you can have type-safe pipeline facilities for any type by defining a channel for that type and generating from the corresponding meta-code. A simple implementation of this pattern is provided in the [pipeline](http://goast.net/x/pipeline) package.

The following is a simple example of concurrent, parallel squaring of a range of numbers

```go
package main

import (
	"fmt"
	"runtime"
)

//go:generate goast write impl goast.net/x/pipeline
type IntPipe chan int

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	done := make(chan bool)
	defer close(done)
	workers := runtime.NumCPU()

	var ch IntPipe = make(chan int)
	go func(ch IntPipe) {
		for i := 0; i < 1000; i++ {
			ch <- i
		}
	}(ch)

	//Fan squaring out over a number of workers and collect the first 20 results
	result := ch.
		Fan(done, workers, func(n int) int { return n * n }).
		Collect(done, 20)

	done <- true
	for _, i := range result {
		fmt.Println(i)
	}
}
```

### Related Types

goast also allows for the generation of derived composite types based on the provided type information, and calls these Related Types. Unlike other types seen in examples so far, the end-user does not need to provide an implementation of the Related Types, only the concrete types used to replace the generic segments of it. To accomlish this, goast also needs to ability to generate new, unique names for Related Types, and allows the meta-code developer to specify how these names should look using a convention of leaving a space for type name replacement with an '_' in the identifier of the Related type

The easy example of this would be a slice sorter Related Type

```go
package sort // import "goast.net/x/sort"

import "sort"

type I interface{}
type Slice []I

//Related Type: Enables sorting as a slice operation
type _Sorter struct {
	Slice
	LessFunc func(I, I) bool
}

func (s _Sorter) Less(i, j int) bool {
	return s.LessFunc(s.Slice[i], s.Slice[j])
}

func (s Slice) Len() int {
	return len(s)
}

func (s Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//Sorts the slice in-place using the related type sorter
func (s Slice) Sort(less func(I, I) bool) {
	sort.Sort(_Sorter{s, less})
}
```

When that is compiled against a type such as `type Contacts []*Contact`, the `Contacts` type gets a `Sort(func(*Contact, *Contact)bool)` method, and a Related Type that looks like the following is generated
```go
type ContactsSorter struct {
	Contacts
	LessFunc	func(*contact, *Contact) bool
}
```

This enables easy arbitrary sorts

```go
//go:generate goast write impl goast.net/x/sort

type Contact struct {
	Id    int
	First string
	Last  string
	Email string
}

type Contacts []*Contact

func main() {

	set := getContacts()
	
	set.Sort(func(a, b *Contact)bool {
		return a.Id < b.Id
	})

	set.Sort(func(a, b *Contact)bool {
		return a.Last > b.Last
	})
}
```

### Multiple Type Parameters

Any number of types are allowed to be specified in template code. For each desired type, assign a new identifier to `interface{}`.

The following defines a type that *maps any type to a slice of any other type*

```go
type K interface{}
type V interface{}
type SliceMap map[K][]V
```

### Partial Structs

`goast write impl` also allows for an amount of structural duck-typing with partially defined structs. This provides a type-safe ability to say "If it looks like X, it can do Y". 

A toy example of this would a 'quittable struct'

```go
//quittable.go
package main

type T struct {
	quit chan bool
}

func (t *T) Quit() {
	t.quit <- true
}
```

Implemented against

```go
package main

//go:generate goast write impl quittable.go
type Process {
	data <- chan string
	value int
	quit chan bool

}
```

### File Naming Control

It can be useful for organizational purposes for generated files to have a naming scheme that identifies them as a generated file. `goast` provides the `--prefix` and `--suffix` flags on the `impl` sub-command to control this behavior.

Example: The following produces a file named `gen_ints_iter.go`

```go
package main

//goast write impl --prefix=gen_ goast.net/x/iter

type Ints []int
```

## Roadmap

goast is still in an alpha/RFC stage of development. Some features that are planned for v1 are

* Projection [Issue](https://github.com/go-goast/goast/issues/4)
* Pruning. [Issue](https://github.com/go-goast/goast/issues/6)
* Support for comments. [Issue](https://github.com/go-goast/goast/issues/5)


## History and acknowledgements

I originally got interested in code generation as a method of genericty in Go when I learned about the [gen](http://clipperhouse.github.io/gen/) package from clipperhouse. When to Go team first announced Go 1.4 and the go:generate proposal, it planted the seed of the idea for goast in my brain and initiated my research into how it might work. In the intervening time I found [gotgo](https://github.com/droundy/gotgo), and more recently (and also quite close to my goals) the [gonerics](https://github.com/bouk/gonerics) package. As projects in the same area as what goast explores, they were all valuble for research and inspiration, as well for providing a contrast against which I wanted to differentiate.


## Licence

goast is released under a GPL v2 licence.

