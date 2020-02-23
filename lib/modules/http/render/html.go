package render

import (
	"html/template"
	"net/http"
)

type Delims struct {
	// Left delimiter, defaults to {{.
	Left string
	// Right delimiter, defaults to }}.
	Right string
}

type HTMLRender interface {
	// Instance returns an HTML instance.
	Instance(string, interface{}) Render
}

type HTMLProduction struct {
	Template *template.Template
	Delims   Delims
}

type HTMLDebug struct {
	Files   []string
	Glob    string
	Delims  Delims
	FuncMap template.FuncMap
}

var htmlContentType = []string{"text/html; charset=utf-8"}

type HTML struct {
	Template *template.Template
	Name     string
	Data     interface{}
}

func (r HTMLProduction) Instance(name string, data interface{}) Render {
	return HTML{
		Template: r.Template,
		Name:     name,
		Data:     data,
	}
}

func (r HTML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	if r.Name == "" {
		return r.Template.Execute(w, r.Data)
	}
	return r.Template.ExecuteTemplate(w, r.Name, r.Data)
}

// WriteContentType (HTML) writes HTML ContentType.
func (r HTML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, htmlContentType)
}
