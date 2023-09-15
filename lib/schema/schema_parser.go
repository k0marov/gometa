package schema

import (
	"encoding/json"
	"log"
	"os"
)

const primaryKeyName = "id"
const primaryKeyType = Int

func Parse(filePath string) Entity {
	jsonEntity := parseJsonFile(filePath)
	log.Printf("succesffuly unmarshalled entity scheme at %q", filePath)

	ent := Entity{Fields: []Field{}}
	hasPrimaryKey := false
	for name, exampleValue := range jsonEntity {
		field := Field{
			Name: name,
			Type: fieldTypeFromInterface(exampleValue),
		}
		if field.Name == primaryKeyName {
			hasPrimaryKey = true
			if field.Type != primaryKeyType {
				log.Fatalf("field name %q is reserved for primary key, it must be of type %q", primaryKeyName, primaryKeyType)
			}
		}
		ent.Fields = append(ent.Fields, field)
	}

	if !hasPrimaryKey {
		log.Fatalf("entity at %q does not have a primary key. Please add %q field", filePath, primaryKeyName)
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
