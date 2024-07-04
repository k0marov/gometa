package setup

import (
	"fmt"
	"github.com/k0marov/gometa.git/lib/helpers"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"os"
	"path/filepath"
)

func addToContainerStruct(goPackageName, entityName, moduleName, containerPath string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, containerPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parsing container.go: %w", err)
	}
	serviceImport := filepath.Join(moduleName, "internal", "services", goPackageName)
	helpers.AddImport(f, serviceImport, "")
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
			Names:   []*ast.Ident{{Name: entityName + "Service"}},
			Type:    &ast.StarExpr{X: &ast.SelectorExpr{X: &ast.Ident{Name: goPackageName}, Sel: &ast.Ident{Name: "ServiceImpl"}}},
			Tag:     nil,
			Comment: nil,
		})
		return false
	})
	containerF, err := os.OpenFile(containerPath, os.O_WRONLY, 0644)
	defer containerF.Close()
	if err != nil {
		return fmt.Errorf("opening container file: %w", err)
	}
	if err := format.Node(containerF, fset, f); err != nil {
		return fmt.Errorf("writing updated container file: %w", err)
	}
	return nil
}
