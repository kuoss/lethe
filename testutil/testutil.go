package testutil

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/util"
	"github.com/spf13/cast"
)

var now = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local)

func Init() {
	config.LoadConfig()
	config.GetConfig().Set("retention.time", "3h")
	config.GetConfig().Set("retention.size", "10m")
	config.SetNow(now)
	config.SetLimit(1000)
	config.SetLogRoot("/tmp/log")
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
}

func SetTestLogFiles() {
	util.Execute("rm -rf /tmp/log")
	dirs := []string{"node/node01", "node/node02", "pod/namspace01", "pod/namespace02"}
	for _, dir := range dirs {
		util.Execute("mkdir -p /tmp/log/" + dir)
		for i := 5; i >= 0; i-- {
			yyyy_mm_dd_hh := now.Add(-time.Duration(i) * time.Hour).UTC().String()[0:13]
			name := strings.Replace(yyyy_mm_dd_hh, " ", "_", -1)
			util.Execute("fallocate -l 1M /tmp/log/" + dir + "/" + name + ".log")
			// t.Log("create file /tmp/log/" + dir + "/" + name + ".log")
		}
	}
}

func SetTestLogs() {
	util.Execute("rm -rf /tmp/log")
	util.Execute(`mkdir -p /tmp/log/pod/namespace01/`)
	util.Execute(`mkdir -p /tmp/log/pod/namespace02/`)
	util.Execute(`mkdir -p /tmp/log/node/node01/`)
	util.Execute(`mkdir -p /tmp/log/node/node02/`)

	// 11 lines
	util.Execute(`cat <<EOF > /tmp/log/pod/namespace01/2009-11-10_22.log
2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar
2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar
2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar
2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum
2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar
2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar
2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar
2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar
2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar
2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar
2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world
EOF`)
	// 3 lines
	util.Execute(`cat <<EOF > /tmp/log/pod/namespace01/2009-11-10_23.log
2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world
2009-11-10T23:01:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world
2009-11-10T23:02:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world
EOF`)
	// 12 lines
	util.Execute(`cat <<EOF > /tmp/log/pod/namespace02/2009-11-10_22.log
2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world
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
EOF`)
	// 10 lines
	util.Execute(`cat <<EOF > /tmp/log/node/node01/2009-11-10_22.log
2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] "RemoveContainer" hello from sidecar
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
EOF`)
	util.Execute(`cat <<EOF > /tmp/log/node/node01/2009-11-10_23.log
2009-11-10T23:00:00.000000Z[node01|containerd] hello world
2009-11-10T23:01:00.000000Z[node01|containerd] hello world
2009-11-10T23:02:00.000000Z[node01|containerd] hello world
EOF`)
	util.Execute(`cat <<EOF > /tmp/log/node/node02/2009-11-10_22.log
2009-11-10T22:58:00.000000Z[node02|containerd] hello world
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
EOF`)
}
