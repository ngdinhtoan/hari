package parser

import (
	"bytes"
	"strings"
)

var nameReplacer = strings.NewReplacer(
	"Id", "ID",
	"Url", "URL",
	"Http", "HTTP",
	"Uri", "URI",
)

// ToCamelCase convert string to CamelCase, the first char also in upper case.
func ToCamelCase(s string) string {
	chunks := alphaNumRegex.FindAll([]byte(s), -1)
	for i, val := range chunks {
		chunks[i] = bytes.Title(val)
	}

	return string(bytes.Join(chunks, nil))
}
