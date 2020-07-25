package render

import (
	"io"
	"text/template"

	"github.com/labstack/echo/v4"
)

// Renderer contains the compiled templates
type Renderer struct {
	templates *template.Template
}

// New creates a new Renderer instance
func New(templates string) *Renderer {
	return &Renderer{
		templates: template.Must(template.ParseGlob(templates)),
	}
}

// Render implements the standard render interface
func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.templates.ExecuteTemplate(w, name, data)
}
