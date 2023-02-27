package testutil

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/spf13/cast"
)

func DirectoryExists(dir string) bool {
	d, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return d.IsDir()
}

func FileExists(file string) bool {
	d, err := os.Stat(file)
	if err != nil {
		return false
	}
	return !d.IsDir()
}

func CheckEqualJSON(t *testing.T, got interface{}, want string, extras ...interface{}) {
	// preprocess
	switch reflect.ValueOf(got).Type().String() {
	case "*errors.errorString":
		got = fmt.Sprintf("%s", got)
	case "[]uint8":
		got = cast.ToString(got)
	}
	// extraMessage
	extraMessage := ""
	for _, extra := range extras {
		if extra != nil {
			extraMessage += fmt.Sprintf("%v", extra)
		}
	}
	temp, err := json.Marshal(got)
	_, file, line, _ := runtime.Caller(1)
	t.Logf("%s:%d: %s\n", filepath.Base(file), line, extraMessage)
	if err != nil {
		t.Fatalf("%s:%d: %s\ncannot marshal to json from got=[%v]", filepath.Base(file), line, extraMessage, got)
	}
	gotJSONString := string(temp)
	if strings.Compare(gotJSONString, want) != 0 {
		t.Fatalf("%s:%d: %s\nwant == `%v`\ngot === `%s`", filepath.Base(file), line, extraMessage, want, gotJSONString)
	}
	t.Logf("want: %s\n             got: %s\n", want, gotJSONString)
}
