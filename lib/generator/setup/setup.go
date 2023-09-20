package setup

import (
	"bytes"
	"github.com/k0marov/gometa/lib/helpers"
	"github.com/k0marov/gometa/lib/schema"
	"io"
	"log"
	"text/template"
)

var setupTemplate = template.Must(template.New("").Parse(`
package {{ .PackageName }} 

import (
    "github.com/jinzhu/gorm"
    "{{ .BasePackage }}/delivery"
    "{{ .BasePackage }}/service"
    "{{ .BasePackage }}/repository"
)

func SetupHandlers(db *gorm.DB) *delivery.{{ .EntityName }}Handlers {
    return delivery.New{{ .EntityName }}Handlers(
		service.New{{.EntityName}}ServiceImpl(
			repository.New{{.EntityName}}RepositoryImpl(db),
		),
	)
}
`))

// TODO: add returning id from Create

func Generate(ent schema.Entity, out io.Writer, packageName, basePackage string) {
	templateData := struct {
		EntityName  string
		PackageName string
		BasePackage string
	}{
		EntityName:  ent.Name,
		PackageName: packageName,
		BasePackage: basePackage,
	}
	var generated bytes.Buffer
	err := setupTemplate.Execute(&generated, templateData)
	if err != nil {
		log.Fatalf("error while executing di setup template: %v", err)
	}
	helpers.WriteFormatted(generated.Bytes(), out)
}
