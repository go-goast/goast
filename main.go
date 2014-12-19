/*
Copyright 2014 James Garfield. All rights reserved.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v1"
	"os"
	"path/filepath"
)

const VERSION = "0.1.2"

func main() {
	var (
		app = kingpin.New("goast", "An AST utility for Go")

		writeCmd         = app.Command("write", "Generate code with various AST transformations")
		writeImpl        = writeCmd.Command("impl", "Generate an implementation of a generically defined file")
		writeImplGeneric = writeImpl.Arg("generic", "Generic file to implement").Required().String()
		writeImplSpec    = writeImpl.Arg("spec", "Spec file that provides types to the generic file").Default(os.ExpandEnv("$GOFILE")).String()

		printCmd       = app.Command("print", "Print various representations of an ast to stdout")
		printDecls     = printCmd.Command("decls", "Print a summary of the top level declarations of a file")
		printDeclsFile = printDecls.Arg("file", "File to inspect").Required().String()
	)

	app.Version(version())

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case writeImpl.FullCommand():
		implement(*writeImplGeneric, *writeImplSpec)

	case printDecls.FullCommand():
		printFileDecls(*printDeclsFile)

	default:
		app.Usage(os.Stdout)
	}

}

func implement(genericFile, specFile string) {

	typeProvider, err := NewFilePackageContext(specFile)

	if err != nil {
		fmt.Println("Error type provider file: ", err)
		return
	}

	imp := NewImplementor(typeProvider)

	fmt.Printf("Implement %s on %s\n", genericFile, specFile)
	RewriteFile(genericFile, filepath.Dir(specFile), imp)
}

func printFileDecls(path string) {
	println("Printing ", path)
	c, err := NewFileContext(path)
	if err != nil {
		println(err)
		return
	}
	PrintDecls(c.File)
}

func version() string {
	return VERSION
}
