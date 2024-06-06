package delivery

import (
	"bytes"
	"fmt"
	"github.com/k0marov/gometa/lib/generator/gen"
	"github.com/k0marov/gometa/lib/helpers"
	"github.com/k0marov/gometa/lib/schema"
	"io"
	"text/template"
)

var handlersTemplate = template.Must(template.New("").Parse(`
package {{.PackageName}}

import (
    "net/http"
    "github.com/gin-gonic/gin"
	"context"
	"strconv"
    "{{ .EntityImport }}"
	"{{ .ModuleName }}/internal/clienterrs"
)

type Service interface {
    Create(ctx context.Context, dto models.Create{{ .EntityName }}DTO) (models.{{ .EntityName }}, error)
    Get(ctx context.Context, id string) (models.{{ .EntityName }}, error) 
    GetAll(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.{{ .EntityName }}, error) 
    Update(ctx context.Context, entity models.{{ .EntityName }}) error 
    Delete(ctx context.Context, id string) error 
}

type Handlers struct {
    svc Service
}

func NewController(svc Service) *Handlers {
    return &Handlers{svc: svc}
}

func (h *Handlers) DefineRoutes(r gin.IRouter) {
    r.POST("/api/v1/{{.PackageName}}s/", h.Create)
	r.GET("/api/v1/{{.PackageName}}s/:id", h.Get) 
	r.GET("/api/v1/{{.PackageName}}s", h.GetAll) 
    r.PUT("/api/v1/{{.PackageName}}s/", h.Update) 
	r.DELETE("/api/v1/{{.PackageName}}s/:id", h.Delete)
}

// Create godoc 
// @Summary Create {{ .EntityName }} 
// @Description Creates a new {{ .EntityName }}, returns an error if it already exists 
// @ID create-{{.PackageName}}
// @Tags {{.PackageName}} 
// @Accept json 
// @Produce json 
// @Param {{.EntityName}} body Create{{.EntityName}}Req true "info about new object"
// @Success 201 {object} {{.EntityName}} 
// @Router /api/v1/{{.PackageName}}s [post]
func (h *Handlers) Create(c *gin.Context) { 
    toCreate := new(Create{{ .EntityName }}Req)
    if c.BindJSON(toCreate) != nil {
        return 
    }
    created, err := h.svc.Create(c.Request.Context(), toCreate.ToDTO())
    if err != nil {
		clienterrs.WriteErrorResponse(c.Writer, err) 
        return 
    }
	c.JSON(http.StatusCreated, MapEntity(created)) 
}

// Get godoc 
// @Summary Get {{ .EntityName }} by id
// @Description Gets {{ .EntityName }} by id
// @ID get-{{.PackageName}}-by-id 
// @Tags {{.PackageName}} 
// @Accept json 
// @Produce json 
// @Success 200 {object} {{.EntityName}} 
// @Router /api/v1/{{.PackageName}}s/:id [get]
func (h *Handlers) Get(c *gin.Context) { 
	id := c.Param("id")
    entity, err := h.svc.Get(c.Request.Context(), id)
    if err != nil {
		clienterrs.WriteErrorResponse(c.Writer, err) 
        return 
    }
    c.JSON(http.StatusOK, MapEntity(entity))
}

// GetAll godoc 
// @Summary Get all {{ .EntityName }}s
// @Description Gets all {{ .EntityName }} entities saved in db with filters and pagination
// @ID get-all-{{.PackageName}}
// @Tags {{.PackageName}} 
// @Param page query int true "page number, starting from 1"
// @Param pageSize query int true "page size"
{{ range $field := .Entity.Fields -}} 
// @Param {{$field.JsonName}} query {{$field.Type.JsonType}} false "filter for {{$field.JsonName}}" 
{{ end -}}
// @Accept json 
// @Produce json 
// @Success 200 {object} []{{.EntityName}} 
// @Router /api/v1/{{.PackageName}}s [get]
func (h *Handlers) GetAll(c *gin.Context) { 
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		clienterrs.WriteErrorResponse(c.Writer, clienterrs.ErrInvalidPageParam)
		return
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		clienterrs.WriteErrorResponse(c.Writer, clienterrs.ErrInvalidPageSizeParam)
		return
	}
	
	filters := getFilters(c.Request) 

    entities, err := h.svc.GetAll(c.Request.Context(), page, pageSize, filters)
    if err != nil {
		clienterrs.WriteErrorResponse(c.Writer, err) 
        return 
    }
	mapped := make([]{{.EntityName}}, len(entities)) 
	for i := range entities {
		mapped[i] = MapEntity(entities[i]) 
	}
    c.JSON(http.StatusOK, mapped)
}

// Update godoc 
// @Summary Update {{ .EntityName }} 
// @Description Updates {{ .EntityName }}, returns an error if it does not exist
// @ID update-{{.PackageName}}
// @Tags {{.PackageName}} 
// @Accept json 
// @Produce json 
// @Param upd body {{.EntityName}} true "info about updating"
// @Success 200 
// @Router /api/v1/{{.PackageName}}s [put]
func (h *Handlers) Update(c *gin.Context) { 
    upd := new({{ .EntityName }})
    if c.BindJSON(&upd) != nil {
        return 
    }
    err := h.svc.Update(c.Request.Context(), upd.ToEntity())
    if err != nil {
		clienterrs.WriteErrorResponse(c.Writer, err) 
        return 
    }
	c.Status(http.StatusOK) 
}

// Delete godoc 
// @Summary Delete {{ .EntityName }} 
// @Description Deletes {{ .EntityName }}, returns an error if it does not exist
// @ID delete-{{.PackageName}}
// @Tags {{.PackageName}}
// @Accept json 
// @Produce json 
// @Success 204 
// @Router /api/v1/{{.PackageName}}s/:id [delete]
func (h *Handlers) Delete(c *gin.Context) { 
	id := c.Param("id")
    err := h.svc.Delete(c.Request.Context(), id)
    if err != nil {
		clienterrs.WriteErrorResponse(c.Writer, err) 
        return 
    }
	c.Status(http.StatusNoContent) 
}

func getFilters(r *http.Request) map[string]string {
	q := r.URL.Query()
	res := make(map[string]string) 
	{{ range $field := .Entity.Fields }} 
		if q.Has("{{$field.JsonName}}") {
			res["{{ $field.JsonName }}"] = q.Get("{{ $field.JsonName }}")
		} {{ end }}	
	return res
}

`))

func GenerateHandlers(out io.Writer, ent schema.Entity, genCtx gen.GenerationContext) error {
	templateData := struct {
		Entity schema.Entity
		gen.GenerationContext
	}{
		Entity:            ent,
		GenerationContext: genCtx,
	}
	var generated bytes.Buffer
	err := handlersTemplate.Execute(&generated, templateData)
	if err != nil {
		return fmt.Errorf("error while executing handlers template: %w", err)
	}
	if err := helpers.WriteFormatted(generated.Bytes(), out); err != nil {
		return fmt.Errorf("while formatting handlers file: %w", err)
	}
	return nil
}
