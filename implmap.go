package main

import (
	"fmt"
	"go/ast"
)

//go:generate goast write impl gen\maputil.go

type ImplMap map[string]ast.Expr

func NewImplMap() ImplMap {
	var im ImplMap = make(map[string]ast.Expr)
	return im
}

func (imp ImplMap) Store(ident string, expr ast.Expr) (ok bool, err error) {
	//Already mapped type, ensure that spec matches expected identifier
	if val, exists := imp[ident]; exists {
		ok = EquivalentExprs(val, expr)
		if !ok {
			err = fmt.Errorf("Cannot implement identifier %s as %s, already mapped to %s", ident, ExprString(expr), ExprString(val))
		}
		return
	}

	imp[ident] = expr
	ok = true
	return
}

func (imp ImplMap) String() (val string) {
	for k, v := range imp {
		val += fmt.Sprintf("%s->%s\n", k, ExprString(v))
	}
	return
}
