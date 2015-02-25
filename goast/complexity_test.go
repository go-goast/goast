package main

import (
	"fmt"
	"testing"
)

type ComplexityCompare struct {
	Smaller, Bigger string
}

func Test_Complexity(t *testing.T) {

	tests := []ComplexityCompare{
		{"type T interface{}", "type T []interface{}"},
		{"type T interface{}", "type T <- chan interface{}"},
		{"type T <- chan interface{}", "type T []<- chan interface{}"},
		{"type T interface{}", "type T map[string]interface{}"},
		{"type T map[string]interface{}", `	type T map[A]interface{}
											type A interface{}`},
		{"type T interface{}", "type T struct{ id int }"},
	}

	for _, cmp := range tests {
		if ok, err := compareComplexities(cmp); !ok {
			t.Error(err)
		}
	}

}

func compareComplexities(cmp ComplexityCompare) (ok bool, err error) {
	src := "package main\n" + cmp.Smaller
	small, err := NewSourceStringContext(src, "smaller.go")
	if err != nil {
		return
	}

	src = "package main\n" + cmp.Bigger
	big, err := NewSourceStringContext(src, "bigger.go")
	if err != nil {
		return
	}

	smType := small.Types()[0]
	bgType := big.Types()[0]

	if ok = (small.Complexity(smType) < big.Complexity(bgType)); ok {
		return
	}

	err = fmt.Errorf("%s (%d) >= %s (%d)", ExprString(smType.Type), small.Complexity(smType), ExprString(bgType.Type), big.Complexity(bgType))
	return
}
