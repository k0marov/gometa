package generator

import (
	"github.com/k0marov/gometa/lib/generator/entity_struct"
	"github.com/k0marov/gometa/lib/generator/repository"
	"github.com/k0marov/gometa/lib/helpers"
	"github.com/k0marov/gometa/lib/schema"
	"path/filepath"
)

func Generate(outDirPath, schemaPath string) {
	ent := schema.Parse(schemaPath)
	crudDir := filepath.Join(outDirPath, ent.JsonName)

	entityFile := helpers.CreateFileRecursively(filepath.Join(crudDir, "entity.go"))
	entity_struct.Generate(ent, entityFile)

	repoFile := helpers.CreateFileRecursively(filepath.Join(crudDir, "repository.go"))
	repository.Generate(ent, repoFile)
}
