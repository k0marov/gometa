package generator

import (
	"fmt"
	"github.com/k0marov/gometa/lib/generator/delivery"
	"github.com/k0marov/gometa/lib/generator/entity_struct"
	"github.com/k0marov/gometa/lib/generator/repository"
	"github.com/k0marov/gometa/lib/generator/service"
	"github.com/k0marov/gometa/lib/helpers"
	"github.com/k0marov/gometa/lib/schema"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetGoModuleName(goModContents io.Reader) (string, error) {
	contents, err := io.ReadAll(goModContents)
	if err != nil {
		return "", fmt.Errorf("while reading from go.mod: %w", err)
	}
	lines := strings.Split(string(contents), "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("got empty contents")
	}
	firstLine := strings.Split(lines[0], " ")
	if len(firstLine) != 2 {
		return "", fmt.Errorf("invalid first go.mod line: %q", lines[0])
	}
	return firstLine[1], nil
}

func Generate(schemaPath, projectDir string, withMain bool) {
	goMod, err := os.Open(filepath.Join(projectDir, "go.mod"))
	if err != nil {
		log.Fatalf("while opening project go.mod: %v", err)
	}
	defer goMod.Close()
	moduleName, err := GetGoModuleName(goMod)
	if err != nil {
		log.Fatalf("getting go module name: %v", err)
	}

	schemaPath, _ = filepath.Abs(schemaPath)

	ent := schema.Parse(schemaPath)
	internalDir := filepath.Join(projectDir, "internal")

	entityFile := helpers.CreateFileRecursively(filepath.Join(internalDir, "models", fmt.Sprintf("%s.go", ent.JsonName)))
	entity_struct.Generate(ent, entityFile, "models")

	entityImportPath := fmt.Sprintf("%s/internal/models", moduleName)

	//entityImportPath := helpers.GetGoImportPath(filepath.Join(internalDir, "entity"))
	//basePackagePath := strings.TrimSuffix(entityImportPath, "/entity")

	repoFile := helpers.CreateFileRecursively(filepath.Join(internalDir, "repository", ent.JsonName, "repository.go"))
	repository.Generate(ent, repoFile, ent.JsonName, entityImportPath)

	serviceFile := helpers.CreateFileRecursively(filepath.Join(internalDir, "services", ent.JsonName, "service.go"))
	service.Generate(ent, serviceFile, entityImportPath)

	handlersFile := helpers.CreateFileRecursively(filepath.Join(internalDir, "web", "controllers", "apiv1", ent.JsonName, "controller.go"))
	delivery.GenerateHandlers(ent, handlersFile, moduleName, ent.JsonName, entityImportPath)

	//setupFile := helpers.CreateFileRecursively(filepath.Join(internalDir, "setup.go"))
	//setup.Generate(ent, setupFile, packageName, basePackagePath)

	//if withMain {
	//	baseMainFile := helpers.CreateFileRecursively(filepath.Join(filepath.Dir(internalDir), "main.go"))
	//	base_main.Generate(ent, baseMainFile, basePackagePath)
	//}
}
