package main

import (
	"github.com/k0marov/gometa/lib/generator"
	"log"
	"os"
)

func main() {
	if err := generator.Generate(os.Args[1], os.Args[2]); err != nil {
		log.Printf("got error: %v", err)
		os.Exit(1)
	}
}
