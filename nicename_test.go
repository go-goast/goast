package main

import (
	"go/parser"
	"testing"
)

type nicenameTest struct {
	expr, expect, note string
}

func Test_NiceName(t *testing.T) {

	tests := []nicenameTest{
		{"chan Thing", "ThingChan", ""},
		{"chan <- Thang", "ThangSendChan", ""},
		{"<-chan Thong", "ThongRecvChan", ""},
		{"chan string", "StringChan", "Should capitalize identifiers"},
		{"[]int", "IntSlice", ""},
		{"*User", "UserPointer", ""},
		{"map[int]string", "StringMapByInt", ""}}

	for _, test := range tests {
		e, _ := parser.ParseExpr(test.expr)
		if nice := NiceName(e); nice != test.expect {
			t.Errorf("Found %s, expected %s. %s", nice, test.expect, test.note)
		}
	}
}
