package web

import (
	"context"
	"net/http"
	"path"
	"regexp"

	"github.com/robinlant/mywiki/internal/wiki/internal/store"
)

var validPath = regexp.MustCompile("^/(edit|view|save|styles)/([a-zA-z+0-9._-]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, context.Context, store.Store, string), s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, ctx, s, m[2])
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, s store.Store, title string) {
	p, ok, err := s.LoadPage(ctx, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if !ok {
		p = &store.Page{Title: title, Body: []byte("This wiki page is emtpy...")}
	}
	renderTemplate(w, "view", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, s store.Store, title string) {
	body := r.FormValue("body")
	p := &store.Page{Title: title, Body: []byte(body)}
	err := s.SavePage(ctx, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func editHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, s store.Store, title string) {
	p, ok, err := s.LoadPage(ctx, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if !ok {
		p = &store.Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func styleHandler(w http.ResponseWriter, r *http.Request, _ context.Context, _ store.Store, style string) {
	w.Header().Set("Content-Type", "text/css")
	http.ServeFile(w, r, path.Join(stylesDir, style))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "root", &store.Page{})
}
