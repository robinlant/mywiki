package web

import (
	"log"
	"net/http"

	"github.com/robinlant/mywiki/internal/store"
)

func Run(store store.Store, addr string) {
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/view/", makePageHandler(viewHandler, store))
	http.HandleFunc("/edit/", makePageHandler(editHandler, store))
	http.HandleFunc("/save/", makePageHandler(saveHandler, store))
	http.HandleFunc("/styles/", makePageHandler(styleHandler, store))
	http.HandleFunc("/search/", makeGenericHandler(searchHandler, store))
	http.HandleFunc("/", makeGenericHandler(rootHandler, store))

	log.Fatal(http.ListenAndServe(addr, nil))
}
