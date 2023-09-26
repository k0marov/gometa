package client_errors

import (
	"bytes"
	"github.com/k0marov/gometa/lib/helpers"
	"io"
	"log"
	"strings"
	"text/template"
)

var clientErrorsTemplate = template.Must(template.New("").Parse(`
	package client_errors

	import (
		"errors"
		"net/http"
		"github.com/gin-gonic/gin"
	)
	
	var Err{{.EntityName}}NotFound = errors.New("{{.EntityNameLower}} not found")

	func ErrorHandlingMiddleware() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Next() 
			lastErr := c.Errors.Last() 
			if lastErr == nil {
				return 
			}
			
			// Add other error checks here
			if errors.Is(lastErr, Err{{.EntityName}}NotFound) {
				c.JSON(http.StatusNotFound, Err{{.EntityName}}NotFound.Error())
			} else {
				c.Status(http.StatusInternalServerError)
			}
		}
	}
`))

func Generate(entityName string, out io.Writer) {
	templateData := struct {
		EntityName      string
		EntityNameLower string
	}{
		EntityName:      entityName,
		EntityNameLower: strings.ToLower(entityName),
	}
	var generated bytes.Buffer
	err := clientErrorsTemplate.Execute(&generated, templateData)
	if err != nil {
		log.Fatalf("error while executing client errors template: %v", err)
	}
	if err := helpers.WriteFormatted(generated.Bytes(), out); err != nil {
		log.Fatalf("while formatting client errors file: %v", err)
	}
}
