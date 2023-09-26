package base_main

import (
	"bytes"
	"github.com/k0marov/gometa/lib/helpers"
	"github.com/k0marov/gometa/lib/schema"
	"io"
	"log"
	"path/filepath"
	"text/template"
)

var baseMainTemplate = template.Must(template.New("").Parse(`
package main 

import (
	"{{.PackageImport}}"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
    "github.com/gin-gonic/gin"
)

func main() {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	r := gin.Default()
	apiv1 := r.Group("/api/v1")

	handlers := {{.PackageName}}.SetupHandlers(db)
	handlers.DefineRoutes(apiv1.Group("{{.APIGroup}}"))

	r.Run()
}
`))

func Generate(ent schema.Entity, out io.Writer, packageImport string) {
	templateData := struct {
		PackageImport string
		PackageName   string
		APIGroup      string
		Entity        schema.Entity
	}{
		PackageImport: packageImport,
		PackageName:   filepath.Base(packageImport),
		APIGroup:      ent.JsonName + "s",
		Entity:        ent,
	}
	var generated bytes.Buffer
	err := baseMainTemplate.Execute(&generated, templateData)
	if err != nil {
		log.Fatalf("error while executing base main template: %v", err)
	}
	if err := helpers.WriteFormatted(generated.Bytes(), out); err != nil {
		log.Fatalf("while formatting base main file: %v", err)
	}
}
