package web

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path"
	"sync"
)

var staticDir = "static"
var tmplDir = path.Join(staticDir, "templates")
var stylesDir = path.Join(staticDir, "styles")

var tmplCache = make(map[string]*template.Template)
var cacheMutex sync.RWMutex

var commonTemplates = []string{
	tmplDir + "/base.html",
	tmplDir + "/header.html",
	tmplDir + "/footer.html",
}

func loadTemplate(tmpl string) (*template.Template, error) {
	cacheMutex.RLock()
	if t, ok := tmplCache[tmpl]; ok {
		cacheMutex.RUnlock()
		return t, nil
	}
	cacheMutex.RUnlock()

	files := append(commonTemplates, tmplDir+"/"+tmpl+".html")

	t, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	cacheMutex.Lock()
	tmplCache[tmpl] = t
	cacheMutex.Unlock()

	return t, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	t, err := loadTemplate(tmpl)
	if err != nil {
		log.Printf("[ERROR] Failed to load template %s: %v", tmpl, err)
		http.Error(w, "Internal Server Error loading template", http.StatusInternalServerError)
		return
	}

	var b bytes.Buffer
	err = t.ExecuteTemplate(&b, tmpl, p)
	if err != nil {
		log.Printf("[ERROR] Failed to execute template %s: %v", tmpl, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	b.WriteTo(w)
}
