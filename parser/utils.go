package parser

import "strings"

func destringify(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return r == '\'' || r == '"'
	})
}
