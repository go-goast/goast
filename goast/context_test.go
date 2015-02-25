package main

import (
	"go/ast"
	"testing"
)

func Test_Lookup(t *testing.T) {
	c, _ := NewFileContext("context.go")

	for _, i := range []string{"Context", "ImportSpecs", "NewFileContext"} {
		if _, ok := c.Lookup(i); !ok {
			t.Error("Failed to find ", i)
		}
	}
}

func Test_LookupImport(t *testing.T) {
	c, _ := NewFileContext("context_test.go")

	for _, i := range []string{"testing"} {
		if _, ok := c.LookupImport(i); !ok {
			t.Error("Failed to find ", i)
		}
	}
}

func Test_LookupType(t *testing.T) {
	c, _ := NewFileContext("context.go")

	for _, i := range []string{"Context", "ImportSpecs"} {
		if _, ok := c.LookupType(i); !ok {
			t.Error("Failed to find ", i)
		}
	}
}

func Test_LookupMethod(t *testing.T) {
	c, _ := NewFileContext("context.go")
	tests := []struct {
		Rcvr string
		Name string
	}{
		{"Context", "Lookup"},
		{"Context", "LookupMethod"},
		{"Context", "LookupType"},
	}
	for _, v := range tests {
		if _, ok := c.LookupMethod(v.Rcvr, v.Name); !ok {
			t.Error("Failed to find ", v)
		}
	}
}

func Test_Funcs(t *testing.T) {
	c, _ := NewFileContext("context_test.go")
	var funcs funcDecls = c.Funcs()

	nameIs := func(name string) func(*ast.FuncDecl) bool {
		return func(f *ast.FuncDecl) bool {
			return f.Name.Name == name
		}
	}

	for _, v := range []string{"Test_Funcs", "Test_Lookup", "Test_LookupType"} {
		if !funcs.Any(nameIs(v)) {
			t.Error("Did not find ", v)
		}
	}
}

type typeThatUsesAnImport map[string]*testing.T
type typeThatHasNoImport map[string]int

func Test_ImportsOf(t *testing.T) {
	c, _ := NewFileContext("context_test.go")

	tp, _ := c.LookupType("typeThatUsesAnImport")

	imports := c.ImportsOf(tp.Type)
	if imports.Len() != 1 {
		t.Error("Invalid number of imports. Expected 1, got ", imports.Len())
	}

	tp, _ = c.LookupType("typeThatHasNoImport")

	imports = c.ImportsOf(tp.Type)
	if imports.Len() != 0 {
		t.Error("Invalid number of imports. Expected 0, got ", imports.Len())
	}

}
