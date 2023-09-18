package helpers

import (
	"strings"
)

func JsonNameToPascalCase(jsonName string) string {
	camelCased := ""
	for _, namePart := range strings.Split(jsonName, "_") {
		camelCased += strings.ToUpper(namePart[:1]) + namePart[1:]
	}
	return camelCased
}
