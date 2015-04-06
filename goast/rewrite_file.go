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
	"bytes"
	"fmt"
	"go/printer"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"strings"
)

//go:generate goast write impl github.com/jamesgarfield/sliceops

//name does not definitely indicate the resulting filename
//it merely acts as a unique identifer that can be used in the filename
type SourceCode struct {
	*Context
	Name string
}

type SourceSet []*SourceCode

type AstTransform interface {
	Transform(*Context) (SourceSet, bool, []error)
}

type writeConfig struct {
	Prefix, Suffix string
}

func RewriteFile(genericSourceFile, outputDirectory string, t AstTransform, cfg writeConfig) {

	gen, err := NewFileContext(genericSourceFile)
	if err != nil {
		printErrors([]error{err})
	}

	codes, ok, errors := t.Transform(gen)
	if !ok {
		printErrors(errors)
		return
	}

	codes.Each(func(s *SourceCode) {
		srcName := strings.TrimSuffix(filepath.Base(genericSourceFile), filepath.Ext(genericSourceFile))
		s.Name = strings.ToLower(fmt.Sprintf("%s%s_%s%s.go", cfg.Prefix, s.Name, srcName, cfg.Suffix))
	})

	for _, source := range codes {
		writeSourceCodeToFile(source, outputDirectory)
	}

}

func printErrors(errors []error) {
	for _, e := range errors {
		fmt.Printf("Error: %s\n", e.Error())
	}
}

func writeSourceCodeToFile(source *SourceCode, outputDirectory string) {
	var b bytes.Buffer

	printer.Fprint(&b, token.NewFileSet(), source.File)
	outPath := filepath.Join(outputDirectory, source.Name)
	ioutil.WriteFile(outPath, b.Bytes(), 0644)
}
