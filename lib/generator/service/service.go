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
	"fmt"
	"{{ .ModuleName }}/pkg/logger"
    "{{ .EntityImport }}"
)

type Repo interface {
    Create(ctx context.Context, dto models.Create{{ .EntityName }}DTO) (models.{{.EntityName}}, error)
    Get(ctx context.Context, id string) (models.{{ .EntityName }}, error) 
    GetAll(ctx context.Context, limit, offset int) ([]models.{{ .EntityName }}, error) 
    Update(ctx context.Context, entity models.{{ .EntityName }}) error 
    Delete(ctx context.Context, id string) error 
}

type ServiceImpl struct {
    repo Repo
}

func NewServiceImpl(repo Repo) *ServiceImpl {
    return &ServiceImpl{repo: repo}
}

func (s *ServiceImpl) Create(ctx context.Context, dto models.Create{{ .EntityName }}DTO) (models.{{ .EntityName }}, error) {
	logger.Debug("creating {{ .EntityName}}", "value", dto) 
    // TODO: add business logic to Service.Create
    return s.repo.Create(ctx, dto)
}

func (s *ServiceImpl) Get(ctx context.Context, id string) (models.{{ .EntityName }}, error) {
	logger.Debug("getting {{ .EntityName}}", "id", id) 
    // TODO: add business logic to Service.Get
    entity, err := s.repo.Get(ctx, id)
	if err != nil {
		return models.{{ .EntityName }}{}, fmt.Errorf("getting from repo: %w", err)
	}
    return entity, nil
}

func (s *ServiceImpl) GetAll(ctx context.Context, page int, pageSize int) ([]models.{{ .EntityName }}, error) {
    // TODO: add business logic 
	logger.Debug("getting all {{ .EntityName}}s", "page", page, "page_size", pageSize) 
	if page <= 0 {
		page = 1 
	} 
	limit := pageSize 
	offset := (page-1) * pageSize 
    entities, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("getting {{.EntityName}} entities from repo: %w", err) 
	}
	return entities, nil 
}

func (s *ServiceImpl) Update(ctx context.Context, entity models.{{ .EntityName }}) error {
	logger.Debug("updating {{.EntityName}}", "value", entity) 
    // TODO: add business logic to Service.Update
    err := s.repo.Update(ctx, entity)
	if err != nil {
		return fmt.Errorf("updating in repo: %w", err)
	}
	return nil	
}

func (s *ServiceImpl) Delete(ctx context.Context, id string) error {
	logger.Debug("deleting {{.EntityName}}", "id", id) 
    // TODO: add business logic to Service.Delete
    err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("deleting in repo: %w", err)
	}
	return nil 
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

// TODO: implement pagination
