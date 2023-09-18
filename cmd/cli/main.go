package main

import (
	"github.com/k0marov/gometa/lib/generator"
	"os"
)

func main() {
	generator.Generate("test/results/", os.Args[1])
}
