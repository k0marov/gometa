package schema

import (
	"encoding/json"
	"fmt"
	"github.com/k0marov/gometa/lib/helpers"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const PrimaryKeyName = "id"

func Parse(filePath string) (Entity, error) {
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
			field.GoName = "ID" // convention of gorm
			if field.Type == Int {
				field.Type = Uint64
			} else if field.Type != String {
				return Entity{}, fmt.Errorf("got primary key type %q, but only int and string are supported", field.Type)
			}
			hasPrimaryKey = true
		}
		ent.Fields = append(ent.Fields, field)
		sort.Slice(ent.Fields, func(i, j int) bool {
			return ent.Fields[i].GoName < ent.Fields[j].GoName
		})
	}

	if !hasPrimaryKey {
		return Entity{}, fmt.Errorf("entity at %q does not have a primary key. Please add %q field", filePath, PrimaryKeyName)
	}

	return ent, nil
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
