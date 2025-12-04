package web

import (
	"log"
	"net/http"

	"github.com/robinlant/mywiki/wiki/internal/quote"
	"github.com/robinlant/mywiki/wiki/internal/store"
)

func Run(store store.Store, addr string, qs *quote.Service) {
	http.HandleFunc("/favicon.ico", faviconHandler) //TODO add other favico formats and android/ios suppoet
	http.HandleFunc("/view/", makePageHandler(viewHandler, store, qs))
	http.HandleFunc("/edit/", makePageHandler(editHandler, store, qs))
	http.HandleFunc("/save/", makePageHandler(saveHandler, store, qs))
	http.HandleFunc("/goto/", makeGenericHandler(gotoHandler, store, qs))
	http.HandleFunc("/styles/", makePageHandler(styleHandler, store, qs))
	http.HandleFunc("/search/", makeGenericHandler(searchHandler, store, qs))
	// root has to be the last as / matches all URIs
	http.HandleFunc("/", makeGenericHandler(rootHandler, store, qs))

	log.Fatal(http.ListenAndServe(addr, nil))
}
