package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/robinlant/mywiki/internal/wiki/internal/store"
)

const tmplDir = "static/templates"
const stylesDir = "static/styles"

var templates = template.Must(template.ParseGlob(tmplDir + "/*.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *store.Page) {
	log.Printf("[INFO] Looking up template %s", tmpl)
	err := templates.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
