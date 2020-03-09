package emoji

import (
	"strings"
	"unicode/utf8"
)

func Remove(content string) string {
	result := strings.Builder{}
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			result.WriteRune(value)
		}
	}
	return result.String()
}
