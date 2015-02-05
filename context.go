package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"path"
	"path/filepath"
	"strconv"
)

//go:generate goast write impl gen/sliceutil.go

type Context struct {
	*ast.File
	*token.FileSet
	*ast.Package
	ast.CommentMap
}

func (c *Context) Clone() (clone *Context, err error) {
	var b bytes.Buffer
	printer.Fprint(&b, c.FileSet, c.File)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "main.go", b.String(), parser.ParseComments)
	if err != nil {
		return
	}

	clone = &Context{file, fset, c.Package, ast.NewCommentMap(fset, file, file.Comments)}
	return
}

type ImportSpecs []*ast.ImportSpec
type funcDecls []*ast.FuncDecl

type fileDecls []ast.Decl

func (s fileDecls) MapToFuncDecl(fn func(ast.Decl) (*ast.FuncDecl, bool)) (result []*ast.FuncDecl) {
	for _, v := range s {
		if r, ok := fn(v); ok {
			result = append(result, r)
		}
	}
	return
}

func (s fileDecls) MapToGenDecl(fn func(ast.Decl) (*ast.GenDecl, bool)) (result []*ast.GenDecl) {
	for _, v := range s {
		if r, ok := fn(v); ok {
			result = append(result, r)
		}
	}
	return
}

func (s fileDecls) MapToTypeSpec(fn func(ast.Decl) (*ast.TypeSpec, bool)) (result []*ast.TypeSpec) {
	for _, v := range s {
		if r, ok := fn(v); ok {
			result = append(result, r)
		}
	}
	return
}

func (c *Context) Lookup(ident string) (obj *ast.Object, ok bool) {
	if c.Package != nil && c.Package.Scope != nil {
		obj = c.Package.Scope.Lookup(ident)
	} else {
		obj = c.File.Scope.Lookup(ident)
	}
	ok = obj != nil
	return
}

func (c *Context) LookupImport(ident string) (i *ast.ImportSpec, ok bool) {
	var imports ImportSpecs = c.File.Imports
	i, ok = imports.First(importSpecMatchesIdentifier(ident))
	return
}

func importSpecMatchesIdentifier(ident string) func(*ast.ImportSpec) bool {
	fn := func(ips *ast.ImportSpec) bool {
		ok := (ips.Name != nil && ips.Name.Name == ident)
		if ok {
			return ok
		}

		importPath, strerr := strconv.Unquote(ips.Path.Value)
		if strerr != nil {
			println("Parsing error while unquoting import path: ", ips.Path.Value)
			return false
		}

		return path.Base(importPath) == ident
	}
	return fn
}

func (c *Context) LookupMethod(rcvr, method string) (f *ast.FuncDecl, ok bool) {
	var funcs funcDecls = c.Funcs()
	f, ok = funcs.First(funcDeclIsMethod(rcvr, method))
	return
}

//Provide a function that determines if a given FuncDecl matches rcvr.method
func funcDeclIsMethod(rcvr, method string) func(*ast.FuncDecl) bool {
	fn := func(f *ast.FuncDecl) bool {
		if f.Name.Name != method || f.Recv == nil || f.Recv.NumFields() == 0 {
			return false
		}
		if name, ok := methodRecieverTypeIdentifier(f); ok && rcvr == name {
			return ok
		}
		return false
	}
	return fn
}

func methodRecieverTypeIdentifier(f *ast.FuncDecl) (name string, ok bool) {
	if f.Recv == nil || f.Recv.NumFields() == 0 {
		return
	}

	//Need to declare separately from initializtion or else
	//the function won't create a closure around find
	//and won't be able to recurse...it's a little janky
	var find func(ast.Expr) (string, bool)
	find = func(e ast.Expr) (string, bool) {
		switch t := e.(type) {
		case *ast.Ident:
			return t.Name, true
		case *ast.StarExpr:
			return find(t.X)
		case *ast.SelectorExpr:
			return find(t.X)
		default:
			return "", false
		}
	}
	return find(f.Recv.List[0].Type)
}

func (c *Context) LookupType(ident string) (t *ast.TypeSpec, ok bool) {
	if obj, exists := c.Lookup(ident); exists {
		if t, ok = obj.Decl.(*ast.TypeSpec); ok {
			return
		}
	}
	return
}

func (c *Context) SetPackage(name string) {
	c.File.Name = ast.NewIdent(name)
}

func (c *Context) Funcs() (funcs []*ast.FuncDecl) {
	var decls fileDecls = c.File.Decls
	funcs = decls.MapToFuncDecl(declAsFuncDecl)
	return
}

func (c *Context) Types() []*ast.TypeSpec {
	var decls fileDecls = c.File.Decls
	types := decls.MapToTypeSpec(declAsTypeSpec)
	return types
}

//Parse a given source file, and its enclosing package directory
func NewFilePackageContext(sourceFile string) (*Context, error) {
	fset := token.NewFileSet()
	packagePath := filepath.Dir(sourceFile)
	pkgs, err := parser.ParseDir(fset, packagePath, nil, 0)
	if err != nil {
		return nil, err
	}
	for _, pkg := range pkgs {
		if file, exists := pkg.Files[sourceFile]; exists {
			cmap := ast.NewCommentMap(fset, file, file.Comments)
			return &Context{file, fset, pkg, cmap}, nil
		}
	}
	return nil, nil
}

//Parse just a given source file, do not include package
func NewFileContext(sourceFile string) (*Context, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, sourceFile, nil, 0)
	if err != nil {
		return nil, err
	}
	cmap := ast.NewCommentMap(fset, file, file.Comments)
	return &Context{file, fset, nil, cmap}, nil
}

//Parse a source string as a given filename
func NewSourceStringContext(source, name string) (*Context, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, name, source, 0)
	if err != nil {
		return nil, err
	}
	cmap := ast.NewCommentMap(fset, file, file.Comments)
	return &Context{file, fset, nil, cmap}, nil
}
