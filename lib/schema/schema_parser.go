package schema

import (
	"encoding/json"
	"github.com/k0marov/gometa/lib/helpers"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const PrimaryKeyName = "id"

func Parse(filePath string) Entity {
	jsonEntity := parseJsonFile(filePath)
	log.Printf("succesffuly unmarshalled entity scheme at %q", filePath)

	_, fileName := filepath.Split(filePath)
	jsonName := strings.ReplaceAll(fileName, ".json", "")
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
		}
		ent.Fields = append(ent.Fields, field)
	}

	if !hasPrimaryKey {
		log.Fatalf("entity at %q does not have a primary key. Please add %q field", filePath, PrimaryKeyName)
	}

	return ent
}

func parseJsonFile(filePath string) map[string]any {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("unable to open entity scheme file %q: %v", filePath, err)
	}
	entityScheme := map[string]any{}
	err = json.NewDecoder(file).Decode(&entityScheme)
	if err != nil {
		log.Fatalf("unable to unmarshal entity scheme at %q as json: %v", filePath, err)
	}
	return entityScheme
}
