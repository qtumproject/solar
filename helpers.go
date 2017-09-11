package solar

import (
	"path"
	"strings"
	"unicode"
	"unicode/utf8"
)

func basenameNoExt(filepath string) string {
	basename := path.Base(filepath)
	return strings.TrimSuffix(basename, path.Ext(basename))
}

func stringLowerFirstRune(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}
