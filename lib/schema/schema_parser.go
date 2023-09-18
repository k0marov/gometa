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
const PrimaryKeyType = Int

func Parse(filePath string) Entity {
	jsonEntity := parseJsonFile(filePath)
	log.Printf("succesffuly unmarshalled entity scheme at %q", filePath)

	_, fileName := filepath.Split(filePath)
	ent := Entity{
		Name:   helpers.JsonNameToCamelCase(strings.ReplaceAll(fileName, ".json", "")),
		Fields: []Field{},
	}
	hasPrimaryKey := false
	for name, exampleValue := range jsonEntity {
		field := Field{
			JsonName: name,
			Type:     fieldTypeFromInterface(exampleValue),
		}
		if field.JsonName == PrimaryKeyName {
			hasPrimaryKey = true
			if field.Type != PrimaryKeyType {
				log.Fatalf("field name %q is reserved for primary key, it must be of type %q", PrimaryKeyName, PrimaryKeyType)
			}
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
