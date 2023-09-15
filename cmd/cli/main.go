package main

import (
	"github.com/k0marov/gometa/lib/schema"
	"log"
	"os"
)

func main() {
	ent := schema.Parse(os.Args[1])
	log.Printf("Got schema %#v", ent)
}
