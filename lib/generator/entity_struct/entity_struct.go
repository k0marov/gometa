package entity_struct

import (
	"github.com/k0marov/gometa/lib/schema"
	"io"
	"log"
	"text/template"
)

var entityTemplate = template.Must(template.New("").Parse(`
package {{ .PackageName }} 

type {{ .Entity.Name }} struct {
	{{ range $field := .Entity.Fields }}
        {{ $field.GoName }} {{ $field.Type.GolangType }} ` + "`json:\"{{$field.JsonName}}\"`" + `
    {{ end }}
}
`))

func Generate(ent schema.Entity, out io.Writer) {
	templateData := struct {
		PackageName string
		Entity      schema.Entity
	}{
		PackageName: ent.JsonName,
		Entity:      ent,
	}
	err := entityTemplate.Execute(out, templateData)
	if err != nil {
		log.Fatalf("error while executing entity struct template: %v", err)
	}
}
