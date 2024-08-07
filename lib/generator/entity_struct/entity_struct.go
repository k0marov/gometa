package entity_struct

import (
	"bytes"
	"fmt"
	"github.com/k0marov/gometa/lib/helpers"
	"github.com/k0marov/gometa/lib/schema"
	"io"
	"text/template"
)

var entityTemplate = template.Must(template.New("").Parse(`
package {{ .PackageName }} 

{{ if .Entity.HasTimeField }} 
import "time"
{{ end }}

type {{ .Entity.Name }} struct {
	{{ range $field := .Entity.Fields }} 
	{{ $field.GoName }} {{ $field.Type.GolangType }} {{ end }}
}

type Create{{ .Entity.Name }}DTO struct {
	{{ range $field := .Entity.Fields }} 
		{{- if $field.IsPrimaryKey }} {{ continue }} {{ end }}
		{{- $field.GoName }} {{ $field.Type.GolangType }} 
	{{ end }}
}
`))

func Generate(out io.Writer, ent schema.Entity, packageName string) error {
	templateData := struct {
		PackageName string
		Entity      schema.Entity
	}{
		PackageName: packageName,
		Entity:      ent,
	}
	var generated bytes.Buffer
	err := entityTemplate.Execute(&generated, templateData)
	if err != nil {
		return fmt.Errorf("error while executing entity struct template: %w", err)
	}
	if err := helpers.WriteFormatted(generated.Bytes(), out); err != nil {
		return fmt.Errorf("while formatting entity struct file: %w", err)
	}
	return nil
}
