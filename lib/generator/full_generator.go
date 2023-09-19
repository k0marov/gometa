package generator

import (
	"github.com/k0marov/gometa/lib/generator/entity_struct"
	"github.com/k0marov/gometa/lib/generator/repository"
	"github.com/k0marov/gometa/lib/generator/service"
	"github.com/k0marov/gometa/lib/helpers"
	"github.com/k0marov/gometa/lib/schema"
	"path/filepath"
)

func Generate(schemaPath string) {
	ent := schema.Parse(schemaPath)
	crudDir := filepath.Dir(schemaPath)

	entityFile := helpers.CreateFileRecursively(filepath.Join(crudDir, "entity.go"))
	entity_struct.Generate(ent, entityFile)

	repoFile := helpers.CreateFileRecursively(filepath.Join(crudDir, "repository.go"))
	repository.Generate(ent, repoFile)

	serviceFile := helpers.CreateFileRecursively(filepath.Join(crudDir, "service.go"))
	service.Generate(ent, serviceFile)
}
