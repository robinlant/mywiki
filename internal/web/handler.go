package web

import (
	"net/http"
	"regexp"

	"github.com/robinlant/mywiki/internal/wiki/internal/store"
)

var validPath = regexp.MustCompile("^/(edit|view|save)/([a-zA-z+0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := store.LoadPage(title)
	if err != nil {
		p = &store.Page{Title: title, Body: []byte("This wiki page is emtpy...")}
	}
	renderTemplate(w, "view", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &store.Page{Title: title, Body: []byte(body)}
	err := p.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := store.LoadPage(title)
	if err != nil {
		p = &store.Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func rootHandler(w http.ResponseWriter, tmpl string, p *store.Page) {

}
