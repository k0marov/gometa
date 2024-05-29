package setup

import (
	"fmt"
	"github.com/k0marov/gometa/lib/helpers"
	"github.com/k0marov/gometa/lib/schema"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"os"
	"path/filepath"
)

func modifyContainer(ent schema.Entity, moduleName, containerPath string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, containerPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parsing container.go: %w", err)
	}
	serviceImport := filepath.Join(moduleName, "internal", "services", ent.JsonName)
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
		return fmt.Errorf("opening container file: %w", err)
	}
	if err := format.Node(containerF, fset, f); err != nil {
		return fmt.Errorf("writing updated container file: %w", err)
	}
	return nil
}

func modifyRouter(ent schema.Entity, moduleName, routerPath string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, routerPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parsing router.go: %w", err)
	}
	controllerImport := filepath.Join(moduleName, "internal", "web", "controllers", "apiv1", ent.JsonName)
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
					X:   &ast.Ident{Name: ent.JsonName},
					Sel: &ast.Ident{Name: "NewController"},
				},
				Args: []ast.Expr{&ast.SelectorExpr{
					X:   &ast.Ident{Name: "container"},
					Sel: &ast.Ident{Name: ent.Name + "Service"},
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

func modifyApplication(ent schema.Entity, moduleName, applicationPath string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, applicationPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parsing application.go: %w", err)
	}
	serviceImport := filepath.Join(moduleName, "internal", "services", ent.JsonName)
	helpers.AddImport(f, serviceImport, "")
	repoImport := filepath.Join(moduleName, "internal", "repository", ent.JsonName)
	repoImportAlias := "repo" + ent.JsonName
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
			Key: &ast.Ident{Name: ent.Name + "Service"},
			Value: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: ent.JsonName},
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

func AddToApplication(ent schema.Entity, projectDir string, moduleName string) error {
	if err := modifyContainer(ent, moduleName, filepath.Join(projectDir, "internal", "app", "dependencies", "container.go")); err != nil {
		return fmt.Errorf("while adding new entity service to container: %w", err)
	}
	if err := modifyRouter(ent, moduleName, filepath.Join(projectDir, "internal", "app", "initializers", "router.go")); err != nil {
		return fmt.Errorf("while adding new entity controller to router: %w", err)
	}
	if err := modifyApplication(ent, moduleName, filepath.Join(projectDir, "internal", "app", "application.go")); err != nil {
		return fmt.Errorf("while adding new entity to application: %w", err)
	}
	return nil
}
