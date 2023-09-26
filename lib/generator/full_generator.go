package generator

import (
	"github.com/k0marov/gometa/lib/generator/base_main"
	"github.com/k0marov/gometa/lib/generator/client_errors"
	"github.com/k0marov/gometa/lib/generator/delivery"
	"github.com/k0marov/gometa/lib/generator/entity_struct"
	"github.com/k0marov/gometa/lib/generator/repository"
	"github.com/k0marov/gometa/lib/generator/service"
	"github.com/k0marov/gometa/lib/generator/setup"
	"github.com/k0marov/gometa/lib/helpers"
	"github.com/k0marov/gometa/lib/schema"
	"path/filepath"
	"strings"
)

func Generate(schemaPath string, withMain bool) {
	schemaPath, _ = filepath.Abs(schemaPath)

	ent := schema.Parse(schemaPath)
	crudDir := filepath.Dir(schemaPath)

	packageName := filepath.Base(crudDir)

	entityFile := helpers.CreateFileRecursively(filepath.Join(crudDir, "entity", "entity.go"))
	entity_struct.Generate(ent, entityFile, packageName)

	entityImportPath := helpers.GetGoImportPath(filepath.Join(crudDir, "entity"))
	basePackagePath := strings.TrimSuffix(entityImportPath, "/entity")

	clientErrorsFile := helpers.CreateFileRecursively(filepath.Join(crudDir, "client_errors", "client_errors.go"))
	client_errors.Generate(ent.Name, clientErrorsFile)

	repoFile := helpers.CreateFileRecursively(filepath.Join(crudDir, "repository", "repository.go"))
	repository.Generate(ent, repoFile, entityImportPath)

	serviceFile := helpers.CreateFileRecursively(filepath.Join(crudDir, "service", "service.go"))
	service.Generate(ent, serviceFile, entityImportPath)

	handlersFile := helpers.CreateFileRecursively(filepath.Join(crudDir, "delivery", "handlers.go"))
	delivery.GenerateHandlers(ent, handlersFile, entityImportPath)

	setupFile := helpers.CreateFileRecursively(filepath.Join(crudDir, "setup.go"))
	setup.Generate(ent, setupFile, packageName, basePackagePath)

	if withMain {
		baseMainFile := helpers.CreateFileRecursively(filepath.Join(filepath.Dir(crudDir), "main.go"))
		base_main.Generate(ent, baseMainFile, basePackagePath)
	}
}
