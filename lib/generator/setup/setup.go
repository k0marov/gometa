package setup

import (
	"github.com/k0marov/gometa/lib/schema"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func AddToApplication(ent schema.Entity, projectDir string, moduleName, entityImport string) {
	containerPath := filepath.Join(projectDir, "internal", "app", "dependencies", "container.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, containerPath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("parsing container.go: %v", err)
	}
	serviceImport := filepath.Join(moduleName, "internal", "services", ent.JsonName)
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
				iSpec := &ast.ImportSpec{Path: &ast.BasicLit{Value: strconv.Quote(serviceImport)}}
				dd.Specs = append(dd.Specs, iSpec)
			}
		}
	}

	astutil.Apply(f, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		if typeSpec.Name.Name != "Container" {
			return true
		}
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}
		structType.Fields.List = append(structType.Fields.List, &ast.Field{
			Doc:     nil,
			Names:   []*ast.Ident{{Name: ent.Name + "Service"}},
			Type:    &ast.StarExpr{X: &ast.SelectorExpr{X: &ast.Ident{Name: ent.JsonName}, Sel: &ast.Ident{Name: "ServiceImpl"}}},
			Tag:     nil,
			Comment: nil,
		})
		return false
	})
	containerF, err := os.OpenFile(containerPath, os.O_WRONLY, 0644)
	defer containerF.Close()
	if err != nil {
		log.Fatalf("opening container file: %v", err)
	}
	if err := format.Node(containerF, fset, f); err != nil {
		log.Fatalf("writing updated container file: %v", err)
	}
}
