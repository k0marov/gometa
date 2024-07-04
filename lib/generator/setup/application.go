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

func addToApplication(goPackageName, entityName, moduleName, applicationPath string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, applicationPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parsing application.go: %w", err)
	}
	serviceImport := filepath.Join(moduleName, "internal", "services", goPackageName)
	helpers.AddImport(f, serviceImport, "")
	repoImport := filepath.Join(moduleName, "internal", "repository", goPackageName)
	repoImportAlias := "repo" + goPackageName
	helpers.AddImport(f, repoImport, repoImportAlias)

	astutil.Apply(f, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		containerLit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		containerLitType, ok := containerLit.Type.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if containerLitType.Sel.Name != "Container" {
			return true
		}
		// TODO: add line break before new service
		containerLit.Elts = append(containerLit.Elts, &ast.KeyValueExpr{
			Key: &ast.Ident{Name: entityName + "Service"},
			Value: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: goPackageName},
					Sel: &ast.Ident{Name: "NewServiceImpl"},
				},
				Args: []ast.Expr{&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{Name: repoImportAlias},
						Sel: &ast.Ident{Name: "NewRepositoryImpl"},
					},
					Args: []ast.Expr{&ast.Ident{Name: "db"}},
				}},
			},
		})
		return false
	})
	applicationF, err := os.OpenFile(applicationPath, os.O_WRONLY, 0644)
	defer applicationF.Close()
	if err != nil {
		return fmt.Errorf("opening application file: %w", err)
	}
	if err := format.Node(applicationF, fset, f); err != nil {
		return fmt.Errorf("writing updated application file: %w", err)
	}
	return nil
}
