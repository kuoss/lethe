package testutil

import (
	"fmt"
	"runtime"
	"strings"
)

func TC() string {
	_, file, line, _ := runtime.Caller(1)
	if slash := strings.LastIndex(file, "/"); slash >= 0 {
		file = file[slash+1:]
	}
	return fmt.Sprintf("TESTCASE:%s:%d", file, line)
}
