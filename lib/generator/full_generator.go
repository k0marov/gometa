package generator

import (
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

func Generate(schemaPath string) {
	schemaPath, _ = filepath.Abs(schemaPath)

	ent := schema.Parse(schemaPath)
	crudDir := filepath.Dir(schemaPath)

	packageName := filepath.Base(crudDir)

	entityFile := helpers.CreateFileRecursively(filepath.Join(crudDir, "entity", "entity.go"))
	entity_struct.Generate(ent, entityFile, packageName)

	entityImportPath := helpers.GetGoImportPath(filepath.Join(crudDir, "entity"))
	basePackagePath := strings.TrimSuffix(entityImportPath, "/entity")

	repoFile := helpers.CreateFileRecursively(filepath.Join(crudDir, "repository", "repository.go"))
	repository.Generate(ent, repoFile, entityImportPath)

	serviceFile := helpers.CreateFileRecursively(filepath.Join(crudDir, "service", "service.go"))
	service.Generate(ent, serviceFile, entityImportPath)

	handlersFile := helpers.CreateFileRecursively(filepath.Join(crudDir, "delivery", "handlers.go"))
	delivery.GenerateHandlers(ent, handlersFile, entityImportPath)

	endpointsFile := helpers.CreateFileRecursively(filepath.Join(crudDir, "delivery", "endpoints.go"))
	delivery.GenerateEndpoints(ent, endpointsFile)

	setupFile := helpers.CreateFileRecursively(filepath.Join(crudDir, "setup.go"))
	setup.Generate(ent, setupFile, packageName, basePackagePath)
}
