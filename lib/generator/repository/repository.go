package repository

import (
	"bytes"
	"fmt"
	"gitlab.sch.ocrv.com.rzd/blockchain/platform/gometa.git/lib/generator/gen"
	"gitlab.sch.ocrv.com.rzd/blockchain/platform/gometa.git/lib/helpers"
	"io"
	"text/template"
)

// TODO: add logs and error wraps in repo template
var repoTemplate = template.Must(template.New("").Parse(`
package {{.PackageName}} 

import (
	"context"
    "gorm.io/gorm"
	"log"
	"fmt"
	"errors"
    "{{.EntityImport}}"
	"{{.ModuleName}}/internal/clienterrs"
)

type RepositoryImpl struct {
    db *gorm.DB
}

func NewRepositoryImpl(db *gorm.DB) *RepositoryImpl {
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
    	log.Fatalf("failed to exec CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\", error: %v", err)
	}
    if err := db.AutoMigrate(&{{.EntityName}}{}); err != nil {
		log.Fatalf("failed migrating for {{.EntityName}}: %v", err) 
	}
    return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) Create(ctx context.Context, dto models.Create{{ .EntityName }}DTO) (models.{{.EntityName}}, error) {
	dbModel := MapCreateDTO(dto) 
    if err := r.db.WithContext(ctx).Create(&dbModel).Error; err != nil {
		return models.{{.EntityName}}{}, fmt.Errorf("creating entity in repo: %w", err)
	}
	return dbModel.ToEntity(), nil
}

func (r *RepositoryImpl) Get(ctx context.Context, id string) (models.{{ .EntityName }}, error) {
    dbModel := new({{ .EntityName }})
    err := r.db.WithContext(ctx).Where("id = ?", id).First(dbModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.{{ .EntityName }}{}, clienterrs.ErrNotFound
	}
    return dbModel.ToEntity(), err
}


func (r *RepositoryImpl) GetAll(ctx context.Context, limit, offset int, filters map[string]string) ([]models.{{ .EntityName }}, error) {
	// for filtration see https://gorm.io/docs/query.html#specify_search_fields

	var dbModels []*{{ .EntityName }} 
	if err := r.db.WithContext(ctx).Where(filters).Limit(limit).Offset(offset).Find(&dbModels).Error; err != nil {
		return nil, fmt.Errorf("getting rows from sql: %w", err) 	
	}
	entities := make([]models.{{ .EntityName }}, len(dbModels)) 
	for i := range dbModels {
		entities[i] = dbModels[i].ToEntity()
	}
	return entities, nil
}

func (r *RepositoryImpl) Update(ctx context.Context, entity models.{{ .EntityName }}) error {
    err := r.db.WithContext(ctx).Model(&{{.EntityName}}{}).Updates(MapEntity(entity)).Error	
	if err != nil {
		return fmt.Errorf("updating in sql: %w", err)
	}
	return nil 
}

func (r *RepositoryImpl) Delete(ctx context.Context, id string) error {
	var cnt int64
    if err := r.db.Model(&{{.EntityName}}{}).Where("id = ?", id).Count(&cnt).Error; err != nil {
		return fmt.Errorf("checking for {{.EntityName}} with id %q existence in db: %w", err) 
	}
	if cnt == 0 {
		return clienterrs.ErrNotFound 
	}
    return r.db.Where("id = ?", id).Delete(&{{.EntityName}}{}).Error
}
`))

func Generate(out io.Writer, genCtx gen.GenerationContext) error {
	templateData := struct {
		gen.GenerationContext
	}{
		GenerationContext: genCtx,
	}
	var generated bytes.Buffer
	err := repoTemplate.Execute(&generated, templateData)
	if err != nil {
		return fmt.Errorf("error while executing repository template: %w", err)
	}
	if err := helpers.WriteFormatted(generated.Bytes(), out); err != nil {
		return fmt.Errorf("while formatting repository file: %w", err)
	}
	return nil
}
