package generator

import (
	"github.com/k0marov/gometa/lib/generator/entity_struct"
	"github.com/k0marov/gometa/lib/generator/repository"
	"github.com/k0marov/gometa/lib/schema"
	"log"
	"os"
	"path/filepath"
)

func Generate(outDirPath, schemaPath string) {
	ent := schema.Parse(schemaPath)
	crudDir := filepath.Join(outDirPath, ent.JsonName)
	err := os.Mkdir(crudDir, 0777)
	if err != nil {
		log.Fatalf("failed creating package directory at %q: %v", crudDir, err)
	}
	entityFile, err := os.Create(filepath.Join(crudDir, "entity.go"))
	if err != nil {
		log.Fatalf("failed creating entity file at %q: %v", crudDir, err)
	}
	entity_struct.Generate(ent, entityFile)

	repoFile, err := os.Create(filepath.Join(crudDir, "repository.go"))
	if err != nil {
		log.Fatalf("failed creating repo file at %q: %v", crudDir, err)
	}
	repository.Generate(ent, repoFile)
}
