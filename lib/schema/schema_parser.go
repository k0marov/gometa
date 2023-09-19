package schema

import (
	"encoding/json"
	"github.com/k0marov/gometa/lib/helpers"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const PrimaryKeyName = "id"

func Parse(filePath string) Entity {
	jsonEntity := parseJsonFile(filePath)

	_, fileName := filepath.Split(filePath)
	jsonName := strings.Split(fileName, ".")[0]
	ent := Entity{
		JsonName: jsonName,
		Name:     helpers.JsonNameToPascalCase(jsonName),
		Fields:   []Field{},
	}
	hasPrimaryKey := false
	for name, exampleValue := range jsonEntity {
		field := Field{
			JsonName: name,
			GoName:   helpers.JsonNameToPascalCase(name),
			Type:     fieldTypeFromInterface(exampleValue),
		}
		if field.JsonName == PrimaryKeyName {
			hasPrimaryKey = true
			field.GoName = "ID" // convention of gorm
			field.Type = Uint64
			log.Printf("Changed %q field to %q with type %q because it's a convention of gorm", field.JsonName, "ID", field.Type)
		}
		ent.Fields = append(ent.Fields, field)
		sort.Slice(ent.Fields, func(i, j int) bool {
			return ent.Fields[i].GoName < ent.Fields[j].GoName
		})
	}

	if !hasPrimaryKey {
		log.Fatalf("entity at %q does not have a primary key. Please add %q field", filePath, PrimaryKeyName)
	}

	return ent
}

func parseJsonFile(filePath string) map[string]any {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("unable to open entity scheme file: %v", err)
	}
	entityScheme := map[string]any{}
	err = json.NewDecoder(file).Decode(&entityScheme)
	if err != nil {
		log.Fatalf("unable to unmarshal entity scheme at %q as json: %v", filePath, err)
	}
	return entityScheme
}
