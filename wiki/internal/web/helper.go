package web

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/robinlant/mywiki/wiki/internal/quote"
	"github.com/robinlant/mywiki/wiki/internal/store"
)

var reference = regexp.MustCompile(`\[{2}(.*?)\]{2}`)

var SetDevMode, DevMode = func() (func(bool), func() bool) {
	var devMode bool

	return func(mode bool) {
			devMode = mode
			log.Printf("[INFO] development mode: %t", devMode)
		}, func() bool {
			return devMode
		}

}()

func getDisplay(p *store.Page) Display {
	return Display{
		Display:  decodeTitle(p.Title),
		ViewHref: "/view/" + p.Title,
		Page:     p,
	}
}

func replaceChars(b string, old rune, new rune) string {
	r := make([]rune, len(b))
	for i, v := range b {
		if v == old {
			r[i] = new
		} else {
			r[i] = v
		}

	}
	return string(r)
}

func encodeTitle(s string) string {
	s = strings.TrimSpace(s)

	return replaceChars(s, ' ', '+')
}

func decodeTitle(s string) string {
	return replaceChars(s, '+', ' ')
}

func addTitleReferences(s []byte) []byte {
	const refClass = "reference"

	s = reference.ReplaceAllFunc(s, func(m []byte) []byte {
		t := m[2 : len(m)-2]
		et := encodeTitle(string(t))
		var b []byte
		b = fmt.Appendf(b, `<a href="/view/%s" class="%s">%s</a>`, et, refClass, t)
		return b
	})

	return s
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

func getRandomQuoteOrWarn(qs *quote.Service) *quote.Quote {
	if qs.BaseUrl == "" {
		return nil
	}
	q, err := qs.GetRandomQuote()
	if err != nil {
		log.Printf("[WARN] quote will not be displayed as quote service returnded err: %s", err.Error())
	}
	return q
}
