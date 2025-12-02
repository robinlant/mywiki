package web

import (
	"log"
	"net/http"

	"github.com/robinlant/mywiki/internal/wiki/internal/store"
)

func Run(store store.Store, addr string) {
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler, store))
	http.HandleFunc("/edit/", makeHandler(editHandler, store))
	http.HandleFunc("/save/", makeHandler(saveHandler, store))
	http.HandleFunc("/styles/", makeHandler(styleHandler, store))
	http.HandleFunc("/", makeRootHandler(store))

	log.Fatal(http.ListenAndServe(addr, nil))
}
