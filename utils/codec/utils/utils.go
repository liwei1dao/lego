package utils

import (
	"unicode"
	"unicode/utf8"
)

///是否是内部字段
func IsExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}
