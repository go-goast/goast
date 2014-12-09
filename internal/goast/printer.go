package goast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"
)

type FilePrinter struct {
	*ast.File
	depth int
}

func (f FilePrinter) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}
	return nil
}

type PrintGenDecl struct {
	*ast.GenDecl
}

func (g PrintGenDecl) String() string {

	strs := []string{}
	for _, spec := range g.Specs {
		switch s := spec.(type) {
		case *ast.ImportSpec:
			if s.Name != nil {
				strs = append(strs, fmt.Sprintf("ImportSpec: %s -> %s", s.Name.Name, s.Path.Value))
			} else {
				strs = append(strs, fmt.Sprintf("ImportSpec: %s", s.Path.Value))
			}

		case *ast.TypeSpec:
			strs = append(strs, fmt.Sprintf("TypeSpec: %s -> %s", s.Name.Name, ExprString(s.Type)))

		case *ast.ValueSpec:
			for i, name := range s.Names {
				strs = append(strs, fmt.Sprintf("ValueSpec: %s -> %s", name.Name, ExprString(s.Values[i])))
			}
		}
	}
	if len(strs) > 1 {
		return "\n\t" + strings.Join(strs, "\n\t")
	}
	return strs[0]

}

type PrintFuncDecl struct {
	*ast.FuncDecl
}

func (f PrintFuncDecl) String() string {
	if f.Recv != nil && f.Recv.NumFields() > 0 {
		return fmt.Sprintf("(%s) %s -> %s", ExprString(f.Recv.List[0].Type), f.Name.Name, ExprString(f.Type))
	}
	return fmt.Sprintf("%s -> %s", f.Name.Name, ExprString(f.Type))
}

func PrintDecls(file *ast.File) {
	for _, d := range file.Decls {
		switch t := d.(type) {
		case *ast.GenDecl:
			fmt.Printf("GenDecl: %s\n", PrintGenDecl{t})
		case *ast.FuncDecl:
			fmt.Printf("FuncDecl: %s\n", PrintFuncDecl{t})
		}

	}
}

func ExprString(e ast.Expr) string {
	var b bytes.Buffer
	printer.Fprint(&b, token.NewFileSet(), e)
	return b.String()
}
