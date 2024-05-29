package repository

import (
	"bytes"
	"fmt"
	"github.com/k0marov/gometa/lib/generator/gen"
	"github.com/k0marov/gometa/lib/helpers"
	"github.com/k0marov/gometa/lib/schema"
	"io"
	"text/template"
)

var entityTemplate = template.Must(template.New("").Parse(`
package {{ .PackageName }} 

import (
    "{{.EntityImport}}"
)

type {{ .Entity.Name }} struct {
	{{ range $field := .Entity.Fields }} 
	{{ $field.GoName }} {{ $field.Type.GolangType }} {{ $field.GetGormTags }} {{ end }}
}

func (e *{{ .Entity.Name }}) ToEntity() models.{{.Entity.Name}} {
	return models.{{ .Entity.Name }}{
		{{ range $field := .Entity.Fields }} 
		{{ $field.GoName }}: e.{{ $field.GoName }}, {{ end }}
	}
}

func MapEntity(e models.{{ .Entity.Name }}) {{ .Entity.Name }} {
	return {{ .Entity.Name }}{
		{{ range $field := .Entity.Fields }} 
		{{ $field.GoName }}: e.{{ $field.GoName }}, {{ end }}
	}
}
`))

// TODO: remove unneeded newline after models.{

func GenerateModel(out io.Writer, ent schema.Entity, genCtx gen.GenerationContext) error {
	templateData := struct {
		gen.GenerationContext
		Entity schema.Entity
	}{
		GenerationContext: genCtx,
		Entity:            ent,
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
