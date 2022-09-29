package util

import (
	"strconv"
	"strings"
)

func SubstrAfter(haystack string, needle string) string {
	pos := strings.Index(haystack, needle)
	if pos == -1 {
		return haystack
	}
	return haystack[pos+len(needle):]
}
func SubstrBefore(haystack string, needle string) string {
	pos := strings.Index(haystack, needle)
	if pos == -1 {
		return haystack
	}
	return haystack[:pos]
}

func CountNewlines(s string) string {
	n := strings.Count(s, "\n")
	if len(s) > 0 && !strings.HasSuffix(s, "\n") {
		n++
	}
	return strconv.Itoa(n)
}
