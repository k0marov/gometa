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
package {{ .PackageName }} 

import (
    "github.com/jinzhu/gorm"
)

type {{ .EntityName }}Repository struct {
    db *gorm.DB
}

func New{{ .EntityName }}Repository(db *gorm.DB) {{ .EntityName }}Repository {
    return {{ .EntityName }}Repository{ db: db}
}

func (r {{ .EntityName }}Repository) Get(id uint64) (*{{ .EntityName }}, error) {
    entity := new({{ .EntityName }})
    err := r.db.Limit(1).Where("id = ?", id).Find(entity).Error()
    return entity, err
}


func (r {{ .EntityName }}Repository) Create(entity *{{ .EntityName }}) error {
    return r.db.Create(entity).Error
}

func (r {{ .EntityName }}Repository) Update(entity *{{ .EntityName }}) error {
    return r.db.Model(entity).Update(entity).Error
}

func (r {{ .EntityName }}Repository) Delete(entity *{{ .EntityName }}) error {
    return r.db.Delete(entity).Error
}
`))

func Generate(ent schema.Entity, out io.Writer) {
	templateData := struct {
		PackageName string
		//EntityPackage string
		EntityName string
	}{
		PackageName: ent.JsonName,
		//EntityPackage: entityPackagePath,
		EntityName: ent.Name,
	}
	err := repoTemplate.Execute(out, templateData)
	if err != nil {
		log.Fatalf("error while executing repository template: %v", err)
	}
}
