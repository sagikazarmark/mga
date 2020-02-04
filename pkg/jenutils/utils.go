package jenutils

import (
	"unicode"
)

// Export makes sure that an identifier will be exported.
func Export(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}

	return ""
}

// Unexport makes sure that an identifier will be unexported.
func Unexport(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}

	return ""
}
