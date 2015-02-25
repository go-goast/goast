package main

import (
	"go/ast"
	"go/parser"
	"testing"
)

func Test_ImplMap(t *testing.T) {

	var imp ImplMap = make(map[string]ast.Expr)

	imp["X"], _ = parser.ParseExpr("string")
	imp["Y"], _ = parser.ParseExpr("float64")

	println(imp.String())

}

type ImplementTest struct {
	Solves    bool
	Gen, Spec string
}

func Test_Implement(t *testing.T) {
	tests := []ImplementTest{
		{true, "type T interface{}", "type Email string"},
		{true, "type T interface{}", "type Celcius float64"},
		{true, `type S []T
				type T interface{}`, "type Bools []bool"},
		{true, `type S []T
				type T interface{}`, "type Exprs []ast.Expr"},
		{true, "type T struct{}", `type User struct{
										Id	int
										Name string}`},
		{true, `type T struct{
					quit chan bool}`, `type StringProcess struct {
										    in   <-chan string
										    out  chan<- string
										    quit chan bool
										}`},
		{true, `type M map[K]V
				type K interface{}
				type V interface{}`, "type IntMap map[string]int"},
		{true, `type M map[K]V
				type K interface{}
				type V interface{}`, "type IntMap map[string]ast.Expr"},
	}

	for _, tst := range tests {
		if ok, err := trySolving(tst); !ok {
			t.Error(err)
		}
	}
}

func trySolving(tst ImplementTest) (ok bool, err error) {
	src := "package main\n" + tst.Gen
	generic, err := NewSourceStringContext(src, "gen.go")
	if err != nil {
		return
	}

	src = "package main\n" + tst.Spec
	provider, err := NewSourceStringContext(src, "provider.go")
	if err != nil {
		return
	}

	cp := ContextPair{generic, provider}
	genType := generic.Types()[0]
	specType := provider.Types()[0]

	var imp ImplMap = make(map[string]ast.Expr)
	satisfied, result, err := Implement(cp, imp, specType, genType)
	ok = (satisfied == tst.Solves)
	println(result.String())
	return
}
