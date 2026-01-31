package templates

import (
	"bytes"
	"embed"
	"html/template"

	"github.com/gin-gonic/gin"
)

//go:embed *
var templateFS embed.FS

var templateCache = make(map[string]*template.Template)

func InitTemplates(e *gin.Engine) {
	t := template.Must(template.ParseFS(templateFS, "./*.html"))
	e.SetHTMLTemplate(t)
}

func RenderTemplate(templateName string, data interface{}) (string, error) {
	t := templateCache[templateName]
	if t == nil {
		loadedTemplate, err := template.ParseFS(templateFS, templateName)
		if err != nil {
			return "", err
		}

		t = loadedTemplate
		templateCache[templateName] = loadedTemplate
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
