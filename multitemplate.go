package multitemplate

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin/render"
)

type Render struct {
	templates map[string]*template.Template
	Funcs     template.FuncMap
}

var HTMLMulti render.Render = Render{}

func New() Render {
	return Render{templates: make(map[string]*template.Template)}
}

func (r Render) Add(name string, tmpl *template.Template) {
	if tmpl == nil {
		panic("template can not be nil")
	}
	if len(name) == 0 {
		panic("template name cannot be empty")
	}
	r.templates[name] = tmpl
}

func (r Render) AddFromFiles(name string, files ...string) *template.Template {
	tmpl := template.Must(template.New(name).Funcs(r.Funcs).ParseFiles(files...))
	r.Add(name, tmpl)
	return tmpl
}

func (r Render) AddFromGlob(name, glob string) *template.Template {
	tmpl := template.Must(template.New(name).Funcs(r.Funcs).ParseGlob(glob))
	r.Add(name, tmpl)
	return tmpl
}

func (r *Render) AddFromString(name, templateString string) *template.Template {
	tmpl := template.Must(template.New(name).Funcs(r.Funcs).Parse(templateString))
	r.Add(name, tmpl)
	return tmpl
}

func (r Render) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	writeHeader(w, code, "text/html")
	name := data[0].(string)
	obj := data[1]

	tmpl, ok := r.templates[name]
	if !ok {
		panic("unknown template name: " + name)
	}

	return tmpl.Execute(w, obj)
}

func writeHeader(w http.ResponseWriter, code int, contentType string) {
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	w.WriteHeader(code)
}
