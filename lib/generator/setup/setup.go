package setup

import (
	"github.com/k0marov/gometa/lib/schema"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
)

func AddToApplication(ent schema.Entity, projectDir string, moduleName, entityImport string) {
	containerPath := filepath.Join(projectDir, "internal", "app", "dependencies", "container.go")
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, containerPath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("parsing container.go: %v", err)
	}
	ast.Inspect(file, func(node ast.Node) bool {
		ast.Print(fset, node)
		return true
	})
	//if err := helpers.WriteFormatted(generated.Bytes(), out); err != nil {
	//	log.Fatalf("while formatting service file: %v", err)
	//}
}
