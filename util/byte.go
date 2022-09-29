package util

import (
	"fmt"
	"strconv"
)

func StringToKB(str string) (int, error) {
	unit := str[len(str)-1:]
	num, err := strconv.Atoi(str[:len(str)-1])
	if err != nil {
		return 0, err
	}
	switch unit {
	case "m":
		return num * 1024, nil
	case "g":
		return num * 1024 * 1024, nil
	}
	return 0, fmt.Errorf("cannot accept unit '%s' in '%s''. allowed units: [m, g]", unit, str)
}
