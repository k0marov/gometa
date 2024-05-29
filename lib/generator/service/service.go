package service

import (
	"bytes"
	"fmt"
	"github.com/k0marov/gometa/lib/generator/gen"
	"github.com/k0marov/gometa/lib/helpers"
	"io"
	"text/template"
)

var serviceTemplate = template.Must(template.New("").Parse(`
package {{ .PackageName }}

import (
	"context"
	"{{ .ModuleName }}/pkg/logger"
    "{{ .EntityImport }}"
)

type Repo interface {
    Create(entity models.{{ .EntityName }}) (models.{{.EntityName}}, error)
    Get(id string) (models.{{ .EntityName }}, error) 
    GetAll() ([]models.{{ .EntityName }}, error) 
    Update(entity models.{{ .EntityName }}) error 
    Delete(id string) error 
}

type ServiceImpl struct {
    repo Repo
}

func NewServiceImpl(repo Repo) *ServiceImpl {
    return &ServiceImpl{repo: repo}
}

func (s *ServiceImpl) Create(ctx context.Context, entity models.{{ .EntityName }}) (models.{{ .EntityName }}, error) {
	logger.Debug("creating {{ .EntityName}}", "value", entity) 
    // TODO: add business logic to Service.Create
    return s.repo.Create(entity)
}

func (s *ServiceImpl) Get(ctx context.Context, id string) (models.{{ .EntityName }}, error) {
	logger.Debug("getting {{ .EntityName}}", "id", id) 
    entity, err := s.repo.Get(id)
    // TODO: add business logic to Service.Get
    return entity, err
}

func (s *ServiceImpl) GetAll(ctx context.Context) ([]models.{{ .EntityName }}, error) {
	logger.Debug("getting all {{ .EntityName}}s") 
    // TODO: add business logic 
    return s.repo.GetAll()
}

func (s *ServiceImpl) Update(ctx context.Context, entity models.{{ .EntityName }}) error {
	logger.Debug("updating {{.EntityName}}", "value", entity) 
    // TODO: add business logic to Service.Update
    return s.repo.Update(entity)
}

func (s *ServiceImpl) Delete(ctx context.Context, id string) error {
	logger.Debug("deleting {{.EntityName}}", "id", id) 
    // TODO: add business logic to Service.Delete
    return s.repo.Delete(id)
}
`))

func Generate(out io.Writer, genCtx gen.GenerationContext) error {
	templateData := struct {
		gen.GenerationContext
	}{
		GenerationContext: genCtx,
	}
	var generated bytes.Buffer
	err := serviceTemplate.Execute(&generated, templateData)
	if err != nil {
		return fmt.Errorf("error while executing service template: %w", err)
	}
	if err := helpers.WriteFormatted(generated.Bytes(), out); err != nil {
		return fmt.Errorf("while formatting service file: %w", err)
	}
	return nil
}
