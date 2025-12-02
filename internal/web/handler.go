package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strconv"

	"github.com/robinlant/mywiki/internal/wiki/internal/store"
)

var validPath = regexp.MustCompile("^/(edit|view|save|styles|search)/([a-zA-z+0-9._-]+)$")

type RootPageData struct {
	Title string
	Posts []*store.Page
}

type handlerFunc func(http.ResponseWriter, *http.Request, context.Context, store.Store, string)

func makePageHandler(fn handlerFunc, s store.Store) http.HandlerFunc {
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
		return
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
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func editHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, s store.Store, title string) {
	p, ok, err := s.LoadPage(ctx, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		p = &store.Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// TODO add erorr handling
func styleHandler(w http.ResponseWriter, r *http.Request, _ context.Context, _ store.Store, style string) {
	w.Header().Set("Content-Type", "text/css")
	http.ServeFile(w, r, path.Join(stylesDir, style))
}

type handlee func(http.ResponseWriter, *http.Request, context.Context, store.Store)

func makeGenericHandler(fn handlee, s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		fn(w, r, ctx, s)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, s store.Store) {
	var q = store.OrderQuery{
		Limit: 10,
		Field: "updatedat",
		Desc:  true,
	}

	ps, err := s.LoadPages(ctx, q)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := RootPageData{
		Title: "Home",
		Posts: ps,
	}

	renderTemplate(w, "root", data)
}

func queryParamOrDefault[T any](r *http.Request, key string, def T, conv func(string) (T, error)) T {
	s := r.URL.Query().Get(key)
	if s == "" {
		return def
	}

	v, err := conv(s)
	if err != nil {
		return def
	}
	return v
}

func searchHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, s store.Store) {
	strToUint := func(s string) (uint, error) {
		i, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(i), nil
	}

	q := store.SearchQuery{
		Search: r.URL.Query().Get("search"),
		Limit:  queryParamOrDefault(r, "limit", 10, strToUint),
		Page:   queryParamOrDefault(r, "page", 1, strToUint),
	}

	p, err := s.SearchPages(ctx, q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var b bytes.Buffer
	fmt.Fprintf(&b, "Search -> '%s'\nLimit -> %d\nPage -> %d\n\n", q.Search, q.Limit, q.Page)

	data, _ := json.MarshalIndent(p, "", "  ")
	b.Write(data)
	b.WriteTo(w)
}

// TODO finish or delete
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	http.ServeFile(w, r, path.Join(staticDir, "favicon.ico"))
}
