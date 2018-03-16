package util

import (
	"strings"
	"text/template"
)

func PrepareString(s string) string {
	r := template.HTMLEscapeString(s)
	r = strings.Replace(s, "'", "''", -1)
	r = "'" + r + "'"
	return r
}
