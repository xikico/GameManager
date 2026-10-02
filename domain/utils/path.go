package utils

import (
	"strings"
)

// NormalizePath removes whitespace and quotes accidentally copied around a path.
func NormalizePath(value string) string {
	return strings.Trim(value, " \t\r\n\"")
}
