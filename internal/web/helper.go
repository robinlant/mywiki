package web

import (
	"bytes"
	"fmt"
	"log"
	"regexp"

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

// TODO rework this decoding mess it has to many reallocations
func decodeT(p *store.Page) string {
	return string(decodeTitle([]byte(p.Title)))
}

func encodeT(title string) string {
	return string(encodeTitle(([]byte(title))))
}

func getDisplay(p *store.Page) Display {
	return Display{
		Display:  decodeT(p),
		ViewHref: "/view/" + p.Title,
		Page:     p,
	}
}

func replaceChars(b []byte, old byte, new byte) []byte {
	r := make([]byte, len(b))
	for i, v := range b {
		if v == old {
			r[i] = new
		} else {
			r[i] = v
		}

	}
	return r
}

func encodeTitle(s []byte) []byte {
	s = bytes.TrimSpace(s)

	return replaceChars(s, ' ', '+')
}

func decodeTitle(s []byte) []byte {
	return replaceChars(s, '+', ' ')
}

func addTitleReferences(s []byte) []byte {
	const refClass = "reference"

	s = reference.ReplaceAllFunc(s, func(m []byte) []byte {
		t := m[2 : len(m)-2]
		et := encodeTitle(t)
		var b []byte
		b = fmt.Appendf(b, `<a href="/view/%s" class="%s">%s</a>`, et, refClass, t)
		return b
	})

	return s
}
