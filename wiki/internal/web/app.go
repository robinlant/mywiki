package web

import (
	"log"
	"net/http"

	"github.com/robinlant/mywiki/wiki/internal/store"
)

func Run(store store.Store, addr string) {
	http.HandleFunc("/favicon.ico", faviconHandler) //TODO add other favico formats and android/ios suppoet
	http.HandleFunc("/view/", makePageHandler(viewHandler, store))
	http.HandleFunc("/edit/", makePageHandler(editHandler, store))
	http.HandleFunc("/save/", makePageHandler(saveHandler, store))
	http.HandleFunc("/goto/", makeGenericHandler(gotoHandler, store))
	http.HandleFunc("/styles/", makePageHandler(styleHandler, store))
	http.HandleFunc("/search/", makeGenericHandler(searchHandler, store))
	// root has to be the last as / matches all URIs
	http.HandleFunc("/", makeGenericHandler(rootHandler, store))

	log.Fatal(http.ListenAndServe(addr, nil))
}
