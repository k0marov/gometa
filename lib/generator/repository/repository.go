package repository

import (
	"github.com/k0marov/gometa/lib/schema"
	"io"
	"log"
	"text/template"
)

// TODO: add logs and error wraps in repo template
// TODO: add creating SQL table in repo
var repoTemplate = template.Must(template.New("").Parse(`
package repository

import (
    "github.com/jinzhu/gorm"
    . {{ .EntityPackage }} 
)

type {{ .EntityName }}Repository struct {
    db *gorm.DB
}

func New{{ .EntityName }}Repository(db *gorm.DB) {{ .EntityName }}Repository {
    return {{ .EntityName }}Repository{ db: db}
}

func (r {{ .EntityName }}Repository) Get({{ .PrimaryName }} {{ .PrimaryType}}) (*{{ .EntityName }}, error) {
    entity := new({{ .EntityName }})
    err := r.db.Limit(1).Where("{{ .PrimarySQLName }} = ?", {{ .PrimaryName }}).Find(entity).Error()
    return entity, err
}


func (r {{ .EntityName }}Repository) Create(entity *{{ .EntityName }}) error {
    return r.db.Create(entity).Error
}

func (r {{ .EntityName }}Repository) Update(entity *{{ .EntityName }}) error {
    return r.db.Model(entity).Update.Error
}

func (r {{ .EntityName }}Repository) Update(entity *{{ .EntityName }}) error {
    return r.db.Model(entity).Update.Error
}

func (r {{ .EntityName }}Repository) Delete(entity *{{ .EntityName }}) error {
    return r.db.Delete(entity).Error
}
`))

func Generate(ent schema.Entity, entityPackagePath string, out io.Writer) {
	log.Printf("Generating repository code...")
	templateData := struct {
		EntityPackage  string
		EntityName     string
		PrimaryType    string
		PrimaryName    string
		PrimarySQLName string
	}{
		EntityPackage:  entityPackagePath,
		EntityName:     ent.Name,
		PrimaryType:    schema.PrimaryKeyType.GolangType(),
		PrimaryName:    schema.PrimaryKeyName,
		PrimarySQLName: schema.PrimaryKeyName,
	}
	err := repoTemplate.Execute(out, templateData)
	if err != nil {
		log.Fatalf("error while executing repository template: %v", err)
	}
	log.Printf("Successfully generated repository code")
}
