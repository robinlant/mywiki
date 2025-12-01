package web

import (
	"html/template"
	"net/http"

	"github.com/robinlant/mywiki/internal/wiki/internal/store"
)

const tmplDir = "internal/web/templates"

var templates = template.Must(template.ParseGlob(tmplDir + "/*.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *store.Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
