package helpers

import (
	"bytes"
	"go/format"
	"io"
	"log"
	"os"
	"os/exec"
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

func WriteFormatted(source []byte, out io.Writer) {
	formatted, err := format.Source(source)
	if err != nil {
		log.Fatalf("error when applying go fmt: %v", err)
		return
	}
	_, err = out.Write(formatted)
	if err != nil {
		log.Fatalf("error when writing formatted code to output file: %v", err)
	}
}

func GetGoImportPath(filepath string) string {
	c := exec.Command("go", "list", filepath)
	var output bytes.Buffer
	c.Stdout = &output
	c.Stderr = &output
	err := c.Run()
	if err != nil {
		log.Fatalf("failed getting package name for %q: %v. Output: %q", filepath, err, output.String())
	}
	result := output.String()
	return result[:len(result)-1] // remove EOF from the end
}
