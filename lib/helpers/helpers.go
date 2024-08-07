package helpers

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func JsonNameToPascalCase(jsonName string) string {
	camelCased := ""
	for _, namePart := range strings.Split(jsonName, "_") {
		camelCased += strings.ToUpper(namePart[:1]) + namePart[1:]
	}
	return camelCased
}

func CreateFileRecursively(path string) (file *os.File, err error) {
	err = os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return nil, fmt.Errorf("failed creating directory to put generated file in: %w", err)
	}
	file, err = os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed creating file for putting generated code in: %w", err)
	}
	return file, nil
}

func WriteFormatted(source []byte, out io.Writer) error {
	formatted, err := format.Source(source)
	if err != nil {
		return fmt.Errorf("error when applying go fmt: %w", err)
	}
	_, err = out.Write(formatted)
	if err != nil {
		return fmt.Errorf("error when writing formatted code to output file: %w", err)
	}
	return nil
}

func AddImport(f *ast.File, importPath, alias string) {
	for i := 0; i < len(f.Decls); i++ {
		d := f.Decls[i]

		switch d.(type) {
		case *ast.FuncDecl:
			// No action
		case *ast.GenDecl:
			dd := d.(*ast.GenDecl)

			// IMPORT Declarations
			if dd.Tok == token.IMPORT {
				// Add the new import
				iSpec := &ast.ImportSpec{Path: &ast.BasicLit{Value: strconv.Quote(importPath)}}
				if alias != "" {
					iSpec.Name = &ast.Ident{Name: alias}
				}
				dd.Specs = append(dd.Specs, iSpec)
			}
		}
	}
}
