package main

import (
	"fmt"
	"go/ast"
)

type ImplRewriter struct {
	ImplMap
}

func (imr ImplRewriter) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch t := node.(type) {
	case *ast.ArrayType:
		return imr.visitArrayType(t)

	case *ast.ChanType:
		return imr.visitChanType(t)

	case *ast.Ellipsis:
		return imr.visitEllipsis(t)

	case *ast.Ident:
		return imr.visitIdent(t)

	case *ast.MapType:
		return imr.visitMapType(t)

	case *ast.ParenExpr:
		return imr.visitParenExpr(t)

	case *ast.StarExpr:
		return imr.visitStarExpr(t)

	case *ast.Field:
		return imr.visitField(t)

	default:
		return imr
	}
}

func (imr ImplRewriter) visitArrayType(node *ast.ArrayType) ast.Visitor {
	if t, ok := imr.replacementType(node.Elt); ok {
		node.Elt = t
		return nil
	}
	return imr
}

func (imr ImplRewriter) visitChanType(node *ast.ChanType) ast.Visitor {
	if t, ok := imr.replacementType(node.Value); ok {
		switch t.(type) {
		case *ast.ChanType:
			//it seems that the syntax of channels of channels gets a bit ambiguous after a while?
			//surrounding the specification type in parens seems to clear that up
			node.Value = &ast.ParenExpr{X: t}
		default:
			node.Value = t
		}

		return nil
	}
	return imr
}

func (imr ImplRewriter) visitEllipsis(node *ast.Ellipsis) ast.Visitor {
	if t, ok := imr.replacementType(node.Elt); ok {
		node.Elt = t
		return nil
	}
	return imr
}

func (imr ImplRewriter) visitIdent(node *ast.Ident) ast.Visitor {
	if t, ok := imr.replacementType(node); ok {
		if id, ok := t.(*ast.Ident); ok {
			node.Name = id.Name
		} else {
			fmt.Printf("Invalid ident replacement for %s: %s", node.Name, ExprString(t))
		}
		return nil
	}
	return imr
}

func (imr ImplRewriter) visitMapType(node *ast.MapType) ast.Visitor {
	var done bool
	if t, ok := imr.replacementType(node.Key); ok {
		node.Key = t
		done = true
	}

	if t, ok := imr.replacementType(node.Value); ok {
		node.Value = t
		done = true
	}

	if done {
		return nil
	}

	return imr
}

func (imr ImplRewriter) visitParenExpr(node *ast.ParenExpr) ast.Visitor {
	if t, ok := imr.replacementType(node.X); ok {
		node.X = t
		return nil
	}
	return imr
}

func (imr ImplRewriter) visitStarExpr(node *ast.StarExpr) ast.Visitor {
	if t, ok := imr.replacementType(node.X); ok {
		node.X = t
		return nil
	}
	return imr
}

func (imr ImplRewriter) visitField(node *ast.Field) ast.Visitor {
	if t, ok := imr.replacementType(node.Type); ok {
		node.Type = t
		return nil
	}
	return imr
}

//Determines what if anything a given ast node should be replaced with
func (imr ImplRewriter) replacementType(node ast.Node) (ast.Expr, bool) {
	switch t := node.(type) {
	case *ast.Ident:
		if val, ok := imr.ImplMap[t.Name]; ok {
			return val, true
		}
		return nil, false

	default:
		return nil, false
	}
}
