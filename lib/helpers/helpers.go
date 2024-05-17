package helpers

import (
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func JsonNameToPascalCase(jsonName string) string {
	camelCased := ""
	for _, namePart := range strings.Split(jsonName, "_") {
		camelCased += strings.ToUpper(namePart[:1]) + namePart[1:]
	}
	return camelCased
}

func CreateFileRecursively(path string) (file *os.File) {
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		log.Fatalf("failed creating directory to put generated file in: %v", err)
	}
	file, err = os.Create(path)
	if err != nil {
		log.Fatalf("failed creating file for putting generated code in: %v", err)
	}
	return file
}

func WriteFormatted(source []byte, out io.Writer) error {
	formatted, err := format.Source(source)
	if err != nil {
		return fmt.Errorf("error when applying go fmt: %w", err)
	}
	_, err = out.Write(formatted)
	if err != nil {
		return fmt.Errorf("error when writing formatted code to output file: %w", err)
	}
	return nil
}
