package web

import (
	"bytes"
	"fmt"
	"regexp"
)

var reference = regexp.MustCompile(`\[{2}(.*?)\]{2}`)

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
