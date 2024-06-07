package main

import (
	"fmt"
	"gitlab.sch.ocrv.com.rzd/blockchain/platform/gometa.git/lib/generator"
	"log"
	"os"
	"slices"
)

const usage = `Usage: gometa <path-to-schema> [<path-to-project>]
  - path-to-schema: path to a json file with schema (see README.md) 
  - path-to-project: optional, path to Go project dir, defaults to '.'
`

func main() {
	if len(os.Args) != 2 && len(os.Args) != 3 || slices.Contains(os.Args, "-h") {
		fmt.Print(usage)
		return
	}
	schemaPath := os.Args[1]

	projectPath := "."
	if len(os.Args) > 2 {
		projectPath = os.Args[2]
	}

	if err := generator.Generate(schemaPath, projectPath); err != nil {
		log.Printf("got error: %v", err)
		os.Exit(1)
	}
}
