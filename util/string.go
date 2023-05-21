package util

import (
	"fmt"
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

func StringToBytes(str string) (int, error) {
	unit := str[len(str)-1:]
	num, err := strconv.Atoi(str[:len(str)-1])
	if err != nil {
		return 0, err
	}
	switch unit {
	case "k":
		return num * 1024, nil
	case "m":
		return num * 1024 * 1024, nil
	case "g":
		return num * 1024 * 1024 * 1024, nil
	}
	return 0, fmt.Errorf("cannot accept unit '%s' in '%s''. allowed units: [k, m, g]", unit, str)
}
