package setup

import (
	"fmt"
	"gitlab.sch.ocrv.com.rzd/blockchain/platform/gometa.git/lib/helpers"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"os"
	"path/filepath"
)

func addCtrlToRouter(goPackageName, entityName, moduleName, routerPath string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, routerPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parsing router.go: %w", err)
	}
	controllerImport := filepath.Join(moduleName, "internal", "web", "controllers", "apiv1", goPackageName)
	helpers.AddImport(f, controllerImport, "")
	astutil.Apply(f, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		if funcDecl.Name.Name != "buildControllers" {
			return true
		}
		// TODO: add line break before new controller
		astutil.Apply(funcDecl, nil, func(c *astutil.Cursor) bool {
			n = c.Node()
			compLit, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}
			arrayType, ok := compLit.Type.(*ast.ArrayType)
			if !ok {
				return true
			}
			arrayTypeExpr, ok := arrayType.Elt.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if arrayTypeExpr.Sel.Name != "Controller" {
				return true
			}
			compLit.Elts = append(compLit.Elts, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: goPackageName},
					Sel: &ast.Ident{Name: "NewController"},
				},
				Args: []ast.Expr{&ast.SelectorExpr{
					X:   &ast.Ident{Name: "container"},
					Sel: &ast.Ident{Name: entityName + "Service"},
				}},
			})
			return false
		})
		return false
	})
	containerF, err := os.OpenFile(routerPath, os.O_WRONLY, 0644)
	defer containerF.Close()
	if err != nil {
		return fmt.Errorf("opening router file: %w", err)
	}
	if err := format.Node(containerF, fset, f); err != nil {
		return fmt.Errorf("writing updated router file: %w", err)
	}
	return nil
}
