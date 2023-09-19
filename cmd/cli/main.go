package main

import (
	"github.com/k0marov/gometa/lib/generator"
	"os"
)

func main() {
	generator.Generate(os.Args[1])
}
