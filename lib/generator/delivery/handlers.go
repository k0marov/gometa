package delivery

import (
	"bytes"
	"github.com/k0marov/gometa/lib/helpers"
	"github.com/k0marov/gometa/lib/schema"
	"io"
	"log"
	"text/template"
)

// TODO: add mappers

var handlersTemplate = template.Must(template.New("").Parse(`
package {{.PackageName}}

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
	"context"
    "{{ .EntityImport }}"
	"{{ .ModuleName }}/internal/clienterrs"
)

type Service interface {
    Create(ctx context.Context, entity *models.{{ .EntityName }}) (*models.{{ .EntityName }}, error)
    Get(ctx context.Context, id uint64) (*models.{{ .EntityName }}, error) 
    GetAll(ctx context.Context) ([]*models.{{ .EntityName }}, error) 
    Update(ctx context.Context, entity *models.{{ .EntityName }}) error 
    Delete(ctx context.Context, id uint64) error 
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
// @Param {{.EntityName}} body models.{{.EntityName}} true "info about new object"
// @Success 201 {object} models.{{.EntityName}} 
// @Router /api/v1/{{.PackageName}}s [post]
func (h *Handlers) Create(c *gin.Context) { 
    toCreate := new(models.{{ .EntityName }})
    if c.BindJSON(toCreate) != nil {
        return 
    }
    created, err := h.svc.Create(c.Request.Context(), toCreate)
    if err != nil {
		clienterrs.WriteErrorResponse(c.Writer, err) 
        return 
    }
	c.JSON(http.StatusCreated, created) 
}

// Get godoc 
// @Summary Get {{ .EntityName }} by id
// @Description Gets {{ .EntityName }} by id
// @ID get-{{.PackageName}}-by-id 
// @Tags {{.PackageName}} 
// @Accept json 
// @Produce json 
// @Success 200 {object} models.{{.EntityName}} 
// @Router /api/v1/{{.PackageName}}s/:id [get]
func (h *Handlers) Get(c *gin.Context) { 
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
		clienterrs.WriteErrorResponse(c.Writer, clienterrs.ErrIDNotInt) 
        return 
    }
    entity, err := h.svc.Get(c.Request.Context(), id)
    if err != nil {
		clienterrs.WriteErrorResponse(c.Writer, err) 
        return 
    }
    c.JSON(http.StatusOK, entity)
}

// GetAll godoc 
// @Summary Get all {{ .EntityName }}s
// @Description Gets all {{ .EntityName }} saved in db
// @ID get-all-{{.PackageName}}
// @Tags {{.PackageName}} 
// @Accept json 
// @Produce json 
// @Success 200 {object} []models.{{.EntityName}} 
// @Router /api/v1/{{.PackageName}}s [get]
func (h *Handlers) GetAll(c *gin.Context) { 
    entities, err := h.svc.GetAll(c.Request.Context())
    if err != nil {
		clienterrs.WriteErrorResponse(c.Writer, err) 
        return 
    }
    c.JSON(http.StatusOK, entities)
}

// Update godoc 
// @Summary Update {{ .EntityName }} 
// @Description Updates {{ .EntityName }}, returns an error if it does not exist
// @ID update-{{.PackageName}}
// @Tags {{.PackageName}} 
// @Accept json 
// @Produce json 
// @Param upd body models.{{.EntityName}} true "info about updating"
// @Success 200 
// @Router /api/v1/{{.PackageName}}s [put]
func (h *Handlers) Update(c *gin.Context) { 
    upd := new(models.{{ .EntityName }})
    if c.BindJSON(&upd) != nil {
        return 
    }
    err := h.svc.Update(c.Request.Context(), upd)
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
		clienterrs.WriteErrorResponse(c.Writer, clienterrs.ErrIDNotInt) 
        return 
    }
    err = h.svc.Delete(c.Request.Context(), id)
    if err != nil {
		clienterrs.WriteErrorResponse(c.Writer, err) 
        return 
    }
	c.Status(http.StatusNoContent) 
}
`))

func GenerateHandlers(ent schema.Entity, out io.Writer, moduleName, packageName, entityImport string) {
	templateData := struct {
		PackageName  string
		ModuleName   string
		EntityName   string
		EntityImport string
	}{
		ModuleName:   moduleName,
		PackageName:  packageName,
		EntityName:   ent.Name,
		EntityImport: entityImport,
	}
	var generated bytes.Buffer
	err := handlersTemplate.Execute(&generated, templateData)
	if err != nil {
		log.Fatalf("error while executing handlers template: %v", err)
	}
	if err := helpers.WriteFormatted(generated.Bytes(), out); err != nil {
		log.Fatalf("while formatting handlers file: %v", err)
	}
}
