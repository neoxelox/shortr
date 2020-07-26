package render

import (
	"html/template"
	"io"
	"shortr/config"

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

type context struct {
	AppPort   int
	AppHost   string
	AppScheme string
	Scope     interface{}
}

// Render implements the standard render interface
func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	scheme := "http"
	if config.GetEnvAsBool("APP_SSL_ENABLED", false) {
		scheme = "https"
	}
	ctx := &context{
		AppPort:   config.GetEnvAsInt("APP_PORT", 80),
		AppHost:   config.GetEnvAsString("APP_HOST", "localhost"),
		AppScheme: scheme,
		Scope:     data,
	}
	return r.templates.ExecuteTemplate(w, name, ctx)
}
