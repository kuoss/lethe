package testutil

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"

	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/kuoss/lethe/config"
	"github.com/spf13/cast"
)

var now = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

const (
	POD         = "pod"
	NODE        = "node"
	namespace01 = "namespace01"
	namespace02 = "namespace02"
	node01      = "node01"
	node02      = "node02"
)

func Init() {
	config.LoadConfig()
	config.GetConfig().Set("retention.time", "3h")
	config.GetConfig().Set("retention.size", "10m")
	config.SetNow(now)
	config.SetLimit(1000)
	//userHome, _ := os.UserHomeDir()
	//config.SetLogRoot(filepath.Join(userHome, "tmp", "log"))
	SetTestLogs()
}

func GetNow() time.Time {
	return now
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
	t.Logf("want: %s\ngot: %s\n", want, gotJSONString)
}

func SetTestLogFiles() {

	RootDirectory := filepath.Join("tmp2", "log")

	if runtime.GOOS == "windows" {
		userHomeDir, _ := os.UserHomeDir()
		RootDirectory = filepath.Join(userHomeDir, RootDirectory)
		fmt.Printf("Running on windows os\n")
	}

	err := os.RemoveAll(RootDirectory)
	if err != nil {
		fmt.Printf("%+v", err)
	}

	dirs := []string{filepath.Join(RootDirectory, POD, namespace01), filepath.Join(RootDirectory, POD, namespace02), filepath.Join(RootDirectory, NODE, node01), filepath.Join(RootDirectory, NODE, node02)}

	for _, dir := range dirs {
		os.MkdirAll(dir, 0777)
		for i := 5; i >= 0; i-- {
			yyyy_mm_dd_hh := now.Add(-time.Duration(i) * time.Hour).UTC().String()[0:13]
			name := strings.Replace(yyyy_mm_dd_hh, " ", "_", -1)

			f, _ := os.OpenFile(filepath.Join(dir, name)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			defer f.Close()
		}
	}
}

func SetTestLogs() {

	var RootDirectory string
	RootDirectory = filepath.Join("tmp", "log")

	if runtime.GOOS == "windows" {
		userHomeDir, _ := os.UserHomeDir()
		RootDirectory = filepath.Join(userHomeDir, RootDirectory)
		fmt.Printf("Running on windows os\n")
	}

	fmt.Printf("Root Directory: %s\n", RootDirectory)

	// flush
	err := os.RemoveAll(RootDirectory)
	if err != nil {
		fmt.Printf("%+v", err)
	}

	os.MkdirAll(filepath.Join(RootDirectory, POD, namespace01), 0777)
	os.MkdirAll(filepath.Join(RootDirectory, POD, namespace02), 0777)
	os.MkdirAll(filepath.Join(RootDirectory, NODE, node01), 0777)
	os.MkdirAll(filepath.Join(RootDirectory, NODE, node02), 0777)
	os.MkdirAll(filepath.Join(RootDirectory, NODE, node02), 0777)

	// 11 lines
	os.WriteFile(filepath.Join(RootDirectory, POD, namespace01, "2009-11-10_22.log"),
		[]byte(`2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar
2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar
2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar
2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum
2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar
2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar
2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar
2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar
2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar
2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar
2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world`), 0777)

	//3 lines
	os.WriteFile(path.Join(RootDirectory, POD, namespace01, "2009-11-10_23.log"),
		[]byte(`2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world
2009-11-10T23:01:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world
2009-11-10T23:02:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world
`), 0777)

	//12 lines
	os.WriteFile(path.Join(RootDirectory, POD, namespace02, "2009-11-10_22.log"),
		[]byte(`2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world
2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum
2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world
2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar
2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar
2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar
2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar
2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar
2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar
2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar
2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar
2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar
`), 0777)

	// 10 lines
	os.WriteFile(path.Join(RootDirectory, NODE, node01, "2009-11-10_22.log"),
		[]byte(`2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] "RemoveContainer" hello from sidecar
2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] "RemoveContainer" hello from sidecar
2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] "RemoveContainer" hello from sidecar
2009-11-10T22:59:00.000000Z[node01|containerd] lerom ipsum
2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] "RemoveContainer" hello from sidecar
2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] "RemoveContainer" hello from sidecar
2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] "RemoveContainer" hello from sidecar
2009-11-10T22:58:00.000000Z[node01|dockerd] hello from sidecar
2009-11-10T22:58:00.000000Z[node01|dockerd] lerom from sidecar
2009-11-10T22:58:00.000000Z[node01|dockerd] hello from sidecar
2009-11-10T22:59:00.000000Z[node01|containerd] hello world
`), 0777)

	// 3 lines
	os.WriteFile(path.Join(RootDirectory, NODE, node01, "2009-11-10_23.log"),
		[]byte(`2009-11-10T23:00:00.000000Z[node01|containerd] hello world
2009-11-10T23:01:00.000000Z[node01|containerd] hello world
2009-11-10T23:02:00.000000Z[node01|containerd] hello world
`), 0777)

	os.WriteFile(path.Join(RootDirectory, NODE, node02, "2009-11-10_22.log"),
		[]byte(`2009-11-10T22:58:00.000000Z[node02|containerd] hello world
2009-11-10T22:58:00.000000Z[node02|containerd] lerom ipsum
2009-11-10T22:58:00.000000Z[node02|containerd] hello world
2009-11-10T22:58:00.000000Z[node02|dockerd] hello from sidecar
2009-11-10T22:58:00.000000Z[node02|dockerd] lerom from sidecar
2009-11-10T22:58:00.000000Z[node02|dockerd] hello from sidecar
2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] "RemoveContainer" hello from sidecar
2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] "RemoveContainer" hello from sidecar
2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] "RemoveContainer" hello from sidecar
2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] "RemoveContainer" hello from sidecar
2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] "RemoveContainer" hello from sidecar
2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] "RemoveContainer" hello from sidecar
`), 0777)

}
