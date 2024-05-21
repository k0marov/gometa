package setup

import (
	"github.com/k0marov/gometa/lib/helpers"
	"github.com/k0marov/gometa/lib/schema"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"log"
	"os"
	"path/filepath"
)

func modifyContainer(ent schema.Entity, moduleName, containerPath string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, containerPath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("parsing container.go: %v", err)
	}
	serviceImport := filepath.Join(moduleName, "internal", "services", ent.JsonName)
	helpers.AddImport(f, serviceImport)
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

func modifyRouter(ent schema.Entity, moduleName, routerPath string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, routerPath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("parsing router.go: %v", err)
	}
	controllerImport := filepath.Join(moduleName, "internal", "web", "controllers", "apiv1", ent.JsonName)
	helpers.AddImport(f, controllerImport)
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
		log.Fatalf("opening router file: %v", err)
	}
	if err := format.Node(containerF, fset, f); err != nil {
		log.Fatalf("writing updated router file: %v", err)
	}
}

func AddToApplication(ent schema.Entity, projectDir string, moduleName string) {
	modifyContainer(ent, moduleName, filepath.Join(projectDir, "internal", "app", "dependencies", "container.go"))
	modifyRouter(ent, moduleName, filepath.Join(projectDir, "internal", "app", "initializers", "router.go"))
}
