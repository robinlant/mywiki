package web

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/robinlant/mywiki/internal/store"
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
