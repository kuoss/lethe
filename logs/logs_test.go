package logs

import (
	"fmt"
	"github.com/kuoss/lethe/storage/driver/factory"
	_ "github.com/kuoss/lethe/storage/driver/filesystem"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/testutil"
	"github.com/kuoss/lethe/util"
)

const (
	Logs  = "logs"
	Count = "count"
)

func Test_Logs_Pod_Get(t *testing.T) {
	testutil.SetTestLogs()

	tests := map[string]struct {
		namespace   string
		pod         string
		container   string
		duration    int
		endTime     time.Time
		LogsOrCount string
		keywords    string
		want        string
	}{
		// modify test case.. for supoorting sort by time

		"all namespaces with duration 5min": {namespace: "*", pod: "", container: "", duration: 5 * 60, endTime: time.Time{}, LogsOrCount: Logs, keywords: "",
			want: `{"IsCounting":false,"Logs":["2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"],"Count":24}`},
		"keyword": {namespace: "namespace*", pod: "", container: "", duration: 1 * 60, endTime: time.Time{}, LogsOrCount: Logs, keywords: "hello",
			want: `{"IsCounting":false,"Logs":[["2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"],"Count":24}`},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(subt *testing.T) {
			got, err := podLogs(tt.namespace, tt.pod, tt.container, tt.duration, tt.endTime, tt.LogsOrCount, tt.keywords)
			testutil.CheckEqualJSON(subt, got, tt.want, fmt.Sprintf("err=%v lines=%d ", err, len(got.Logs)))
		})
	}

	/*
		got, _ = podLogs("*", "", "", 5*60, time.Time{}, Logs, "")
		want = `["2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]`
		testutil.CheckEqualJSON(t, got, want, "length=", len(got.Logs))

		got, _ = podLogs("*", "", "", 0, time.Time{}, Logs, "")
		want = `["2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]`
		testutil.CheckEqualJSON(t, got, want, "length=", len(got.Logs))

		got, _ = podLogs("namespace01", "nginx", "", 1*60, time.Time{}, Logs, "")
		want = `[]`
		testutil.CheckEqualJSON(t, got, want, "length=", len(got.Logs))

		got, _ = podLogs("namespace01", "nginx-*", "", 1*60, time.Time{}, Logs, "")
		want = `["2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]`
		testutil.CheckEqualJSON(t, got, want, "length=", len(got.Logs))

		got, _ = podLogs("namespace01", "", "nginx", 1*60, time.Time{}, Logs, "")
		want = `["2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]`
		testutil.CheckEqualJSON(t, got, want, "length=", len(got.Logs))

		got, _ = podLogs("namespace", "", "nginx", 1*60, time.Time{}, Logs, "")
		want = `[]`
		testutil.CheckEqualJSON(t, got, want, "length=", len(got.Logs))

		got, _ = podLogs("namespace*", "", "nginx", 1*60, time.Time{}, Logs, "")
		want = `["2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]`
		testutil.CheckEqualJSON(t, got, want, "length=", len(got.Logs))

		got, _ = podLogs("namespace*", "", "nginx", 1*60, time.Time{}, Logs, "hello")
		want = `["2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]`
		testutil.CheckEqualJSON(t, got, want, "length=", len(got.Logs))

		got, _ = podLogs("namespace*", "", "nginx", 2*60, time.Time{}, Logs, "hello")
		want = `["2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]`
		testutil.CheckEqualJSON(t, got, want, "length=", len(got.Logs))

		got, _ = podLogs("namespace*", "", "nginx", 3*60, time.Time{}, Logs, "hello")
		want = `["2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]`
		testutil.CheckEqualJSON(t, got, want, "length=", len(got.Logs))

		got, _ = podLogs("namespace*", "", "nginx", 5*60, time.Time{}, Logs, "hello")
		want = `["2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]`
		testutil.CheckEqualJSON(t, got, want, "length=", len(got.Logs))

	*/
}

func Test_Logs_Pod_Count(t *testing.T) {
	testutil.SetTestLogs()

	var got Result
	var want string

	got, _ = podLogs("*", "", "", 0, time.Time{}, Logs, "")
	want = "24"
	testutil.CheckEqualJSON(t, got, want, "length=", got)

	got, _ = podLogs("namespace01", "nginx", "", 1*60, time.Time{}, Count, "")
	want = "0"
	testutil.CheckEqualJSON(t, got, want, "length=", got)

	got, _ = podLogs("namespace01", "nginx-*", "", 2*60, time.Time{}, Count, "")
	want = "6"
	testutil.CheckEqualJSON(t, got, want, "length=", got)

	got, _ = podLogs("namespace01", "", "nginx", 2*60, time.Time{}, Count, "")
	want = "3"
	testutil.CheckEqualJSON(t, got, want, "length=", got)

	got, _ = podLogs("namespace", "", "nginx", 1*60, time.Time{}, Count, "")
	want = "0"
	testutil.CheckEqualJSON(t, got, want, "length=", got)

	got, _ = podLogs("namespace*", "", "nginx", 2*60, time.Time{}, Count, "")
	want = "6"
	testutil.CheckEqualJSON(t, got, want, "length=", got)

	got, _ = podLogs("namespace*", "", "nginx", 1*60, time.Time{}, Count, "hello")
	want = "2"
	testutil.CheckEqualJSON(t, got, want, "length=", got)

	got, _ = podLogs("namespace*", "", "nginx", 2*60, time.Time{}, Count, "hello")
	want = "4"
	testutil.CheckEqualJSON(t, got, want, "length=", got)

	got, _ = podLogs("namespace*", "", "nginx", 3*60, time.Time{}, Count, "hello")
	want = "4"
	testutil.CheckEqualJSON(t, got, want, "length=", got)

	got, _ = podLogs("namespace*", "", "nginx", 4*60, time.Time{}, Count, "")
	want = "4"
	testutil.CheckEqualJSON(t, got, want, "length=", got)
}

func Test_Logs_Node_Get(t *testing.T) {
	testutil.SetTestLogs()

	tests := map[string]struct {
		node        string
		process     string
		duration    int
		endTime     time.Time
		LogsOrCount string
		keywords    string
		want        string
	}{
		// modify test case.. for supoorting sort by time
		"all nodes with duration 5min": {node: "*", process: "", duration: 5 * 60, endTime: time.Time{}, LogsOrCount: Logs, keywords: "",
			want: `{"IsCounting":false,"Logs":["2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node01|dockerd] hello from sidecar","2009-11-10T22:58:00.000000Z[node01|dockerd] lerom from sidecar","2009-11-10T22:58:00.000000Z[node01|dockerd] hello from sidecar","2009-11-10T22:58:00.000000Z[node02|containerd] hello world","2009-11-10T22:58:00.000000Z[node02|containerd] lerom ipsum","2009-11-10T22:58:00.000000Z[node02|containerd] hello world","2009-11-10T22:58:00.000000Z[node02|dockerd] hello from sidecar","2009-11-10T22:58:00.000000Z[node02|dockerd] lerom from sidecar","2009-11-10T22:58:00.000000Z[node02|dockerd] hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:59:00.000000Z[node01|containerd] lerom ipsum","2009-11-10T22:59:00.000000Z[node01|containerd] hello world","2009-11-10T23:00:00.000000Z[node01|containerd] hello world"],"Count":24}`,
		},
		"nodes01 kubelet 5min": {node: "node01", process: "kubelet", duration: 5 * 60, endTime: time.Time{}, LogsOrCount: Logs, keywords: "",
			want: `{"IsCounting":false,"Logs":["2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"],"Count":6}`,
		},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(subt *testing.T) {
			got, err := nodeLogs(tt.node, tt.process, tt.duration, tt.endTime, tt.LogsOrCount, tt.keywords)
			testutil.CheckEqualJSON(subt, got, tt.want, fmt.Sprintf("err=%v lines=%d ", err, len(got.Logs)))
		})
	}

	/*
		got, _ = nodeLogs("*", "", 5*60, time.Time{}, Logs, "")
		want = `["2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node01|dockerd] hello from sidecar","2009-11-10T22:58:00.000000Z[node01|dockerd] hello from sidecar","2009-11-10T22:58:00.000000Z[node01|dockerd] lerom from sidecar","2009-11-10T22:58:00.000000Z[node02|containerd] hello world","2009-11-10T22:58:00.000000Z[node02|containerd] hello world","2009-11-10T22:58:00.000000Z[node02|containerd] lerom ipsum","2009-11-10T22:58:00.000000Z[node02|dockerd] hello from sidecar","2009-11-10T22:58:00.000000Z[node02|dockerd] hello from sidecar","2009-11-10T22:58:00.000000Z[node02|dockerd] lerom from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:59:00.000000Z[node01|containerd] hello world","2009-11-10T22:59:00.000000Z[node01|containerd] lerom ipsum","2009-11-10T23:00:00.000000Z[node01|containerd] hello world"]`
		testutil.CheckEqualJSON(t, got, want, " length=", len(got.Logs))

		got, _ = nodeLogs("node01", "kubelet", 5*60, time.Time{}, Logs, "")
		want = `["2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"]`
		testutil.CheckEqualJSON(t, got, want, " length=", len(got.Logs))

		got, _ = nodeLogs("node01", "*", 1*60, time.Time{}, Logs, "")
		want = `["2009-11-10T22:59:00.000000Z[node01|containerd] lerom ipsum","2009-11-10T22:59:00.000000Z[node01|containerd] hello world","2009-11-10T23:00:00.000000Z[node01|containerd] hello world"]`
		testutil.CheckEqualJSON(t, got, want, " length=", len(got.Logs))

		got, _ = nodeLogs("node01", "kube*", 5*60, time.Time{}, Logs, "")
		want = `["2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"]`
		testutil.CheckEqualJSON(t, got, want, " length=", len(got.Logs))

		got, _ = nodeLogs("node01", "kubelet", 4*60, time.Time{}, Logs, "")
		want = `["2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"]`
		testutil.CheckEqualJSON(t, got, want, " length=", len(got.Logs))

		got, _ = nodeLogs("node", "kubelet", 1*60, time.Time{}, Logs, "")
		want = `[]`
		testutil.CheckEqualJSON(t, got, want, " length=", len(got.Logs))

		got, _ = nodeLogs("node*", "kubelet", 5*60, time.Time{}, Logs, "")
		want = `["2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"]`
		testutil.CheckEqualJSON(t, got, want, " length=", len(got.Logs))

		got, _ = nodeLogs("node*", "kubelet", 5*60, time.Time{}, Logs, "hello")
		want = `["2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:56:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"]`
		testutil.CheckEqualJSON(t, got, want, " length=", len(got.Logs))

		got, _ = nodeLogs("node*", "kubelet", 1*60, time.Time{}, Logs, "hello")
		want = `[]`
		testutil.CheckEqualJSON(t, got, want, " length=", len(got.Logs))

		got, _ = nodeLogs("node*", "kubelet", 3*60, time.Time{}, Logs, "hello")
		want = `["2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"]`
		testutil.CheckEqualJSON(t, got, want, " length=", len(got.Logs))

		got, _ = nodeLogs("node*", "kubelet", 3*60, time.Time{}, Logs, "hello")
		want = `["2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:57:00.000000Z[node01|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar","2009-11-10T22:58:00.000000Z[node02|kubelet] I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"]`
		testutil.CheckEqualJSON(t, got, want, " length=", len(got.Logs))

	*/
}

func Test_Logs_Node_Count(t *testing.T) {
	testutil.SetTestLogs()

	var got Result
	var want string

	got, _ = nodeLogs("*", "", 5*60, time.Time{}, Count, "")
	want = "24"
	testutil.CheckEqualJSON(t, got, want)

	got, _ = nodeLogs("node01", "kubelet", 5*60, time.Time{}, Count, "")
	want = "6"
	testutil.CheckEqualJSON(t, got, want)

	got, _ = nodeLogs("node01", "kube*", 5*60, time.Time{}, Count, "")
	want = "6"
	testutil.CheckEqualJSON(t, got, want)

	got, _ = nodeLogs("node01", "", 5*60, time.Time{}, Count, "kubelet")
	want = "6"
	testutil.CheckEqualJSON(t, got, want)

	got, _ = nodeLogs("node", "", 1*60, time.Time{}, Count, "kubelet")
	want = "0"
	testutil.CheckEqualJSON(t, got, want)

	got, _ = nodeLogs("node*", "", 5*60, time.Time{}, Count, "kubelet")
	want = "12"
	testutil.CheckEqualJSON(t, got, want)

	got, _ = nodeLogs("node*", "kubelet", 5*60, time.Time{}, Count, "hello")
	want = "12"
	testutil.CheckEqualJSON(t, got, want)

	got, _ = nodeLogs("node*", "kubelet", 1*60, time.Time{}, Count, "hello")
	want = "0"
	testutil.CheckEqualJSON(t, got, want)

	got, _ = nodeLogs("node*", "kubelet", 2*60, time.Time{}, Count, "hello")
	want = "6"
	testutil.CheckEqualJSON(t, got, want)

	got, _ = nodeLogs("node*", "kubelet", 3*60, time.Time{}, Count, "hello")
	want = "9"
	testutil.CheckEqualJSON(t, got, want)

	got, _ = nodeLogs("node*", "kubelet", 4*60, time.Time{}, Count, "hello")
	want = "12"
	testutil.CheckEqualJSON(t, got, want)
}

func getLogFiles(command string) []string {
	out, _, _ := util.Execute(command)
	return strings.Split(strings.TrimRight(out, "\n"), "\n")
}
func Test_getLogFilesSearchs(t *testing.T) {
	testutil.SetTestLogs()

	var got []logFileSearch
	var want string

	// list files
	got = matchedTimePatternFiles(getLogFiles("ls /tmp/log/*/*/*.log"), 0, time.Time{})
	want = `[{"File":"2009-11-10_23.log","Targets":["node01","namespace01"],"TimePatterns":["2009-11-10T23:00:00"]},{"File":"2009-11-10_22.log","Targets":["node01","node02","namespace01","namespace02"],"TimePatterns":[]}]`
	testutil.CheckEqualJSON(t, got, want, got)

	got = matchedTimePatternFiles(getLogFiles("ls /tmp/log/*/*/*.log"), 5, time.Time{})
	want = `[{"File":"2009-11-10_23.log","Targets":["node01","namespace01"],"TimePatterns":["2009-11-10T23:00:00"]},{"File":"2009-11-10_22.log","Targets":["node01","node02","namespace01","namespace02"],"TimePatterns":["2009-11-10T22:59:55","2009-11-10T22:59:56","2009-11-10T22:59:57","2009-11-10T22:59:58","2009-11-10T22:59:59"]}]`
	testutil.CheckEqualJSON(t, got, want, got)

	got = matchedTimePatternFiles(getLogFiles("ls /tmp/log/*/*/*.log"), 10, time.Time{})
	want = `[{"File":"2009-11-10_23.log","Targets":["node01","namespace01"],"TimePatterns":["2009-11-10T23:00:00"]},{"File":"2009-11-10_22.log","Targets":["node01","node02","namespace01","namespace02"],"TimePatterns":["2009-11-10T22:59:5"]}]`
	testutil.CheckEqualJSON(t, got, want, got)

	got = matchedTimePatternFiles(getLogFiles("ls /tmp/log/*/*/*.log"), 19, time.Time{})
	want = `[{"File":"2009-11-10_23.log","Targets":["node01","namespace01"],"TimePatterns":["2009-11-10T23:00:00"]},{"File":"2009-11-10_22.log","Targets":["node01","node02","namespace01","namespace02"],"TimePatterns":["2009-11-10T22:59:41","2009-11-10T22:59:42","2009-11-10T22:59:43","2009-11-10T22:59:44","2009-11-10T22:59:45","2009-11-10T22:59:46","2009-11-10T22:59:47","2009-11-10T22:59:48","2009-11-10T22:59:49","2009-11-10T22:59:5"]}]`
	testutil.CheckEqualJSON(t, got, want, got)

	got = matchedTimePatternFiles(getLogFiles("ls /tmp/log/*/*/*.log"), 20, time.Time{})
	want = `[{"File":"2009-11-10_23.log","Targets":["node01","namespace01"],"TimePatterns":["2009-11-10T23:00:00"]},{"File":"2009-11-10_22.log","Targets":["node01","node02","namespace01","namespace02"],"TimePatterns":["2009-11-10T22:59:4","2009-11-10T22:59:5"]}]`
	testutil.CheckEqualJSON(t, got, want, got)

	got = matchedTimePatternFiles(getLogFiles("ls /tmp/log/*/*/*.log"), 99, time.Time{})
	want = `[{"File":"2009-11-10_23.log","Targets":["node01","namespace01"],"TimePatterns":["2009-11-10T23:00:00"]},{"File":"2009-11-10_22.log","Targets":["node01","node02","namespace01","namespace02"],"TimePatterns":["2009-11-10T22:58:21","2009-11-10T22:58:22","2009-11-10T22:58:23","2009-11-10T22:58:24","2009-11-10T22:58:25","2009-11-10T22:58:26","2009-11-10T22:58:27","2009-11-10T22:58:28","2009-11-10T22:58:29","2009-11-10T22:59","2009-11-10T22:58:3","2009-11-10T22:58:4","2009-11-10T22:58:5"]}]`
	testutil.CheckEqualJSON(t, got, want, got)

	got = matchedTimePatternFiles(getLogFiles("ls /tmp/log/*/*/*.log"), 119, time.Time{})
	want = `[{"File":"2009-11-10_23.log","Targets":["node01","namespace01"],"TimePatterns":["2009-11-10T23:00:00"]},{"File":"2009-11-10_22.log","Targets":["node01","node02","namespace01","namespace02"],"TimePatterns":["2009-11-10T22:58:01","2009-11-10T22:58:02","2009-11-10T22:58:03","2009-11-10T22:58:04","2009-11-10T22:58:05","2009-11-10T22:58:06","2009-11-10T22:58:07","2009-11-10T22:58:08","2009-11-10T22:58:09","2009-11-10T22:59","2009-11-10T22:58:1","2009-11-10T22:58:2","2009-11-10T22:58:3","2009-11-10T22:58:4","2009-11-10T22:58:5"]}]`
	testutil.CheckEqualJSON(t, got, want, got)

	got = matchedTimePatternFiles(getLogFiles("ls /tmp/log/*/*/*.log"), 120, time.Time{})
	want = `[{"File":"2009-11-10_23.log","Targets":["node01","namespace01"],"TimePatterns":["2009-11-10T23:00:00"]},{"File":"2009-11-10_22.log","Targets":["node01","node02","namespace01","namespace02"],"TimePatterns":["2009-11-10T22:58","2009-11-10T22:59"]}]`
	testutil.CheckEqualJSON(t, got, want, got)

	got = matchedTimePatternFiles(getLogFiles("ls /tmp/log/*/*/*.log"), 999, time.Time{})
	want = `[{"File":"2009-11-10_23.log","Targets":["node01","namespace01"],"TimePatterns":["2009-11-10T23:00:00"]},{"File":"2009-11-10_22.log","Targets":["node01","node02","namespace01","namespace02"],"TimePatterns":["2009-11-10T22:43:21","2009-11-10T22:43:22","2009-11-10T22:43:23","2009-11-10T22:43:24","2009-11-10T22:43:25","2009-11-10T22:43:26","2009-11-10T22:43:27","2009-11-10T22:43:28","2009-11-10T22:43:29","2009-11-10T22:5","2009-11-10T22:44","2009-11-10T22:45","2009-11-10T22:46","2009-11-10T22:47","2009-11-10T22:48","2009-11-10T22:49","2009-11-10T22:43:3","2009-11-10T22:43:4","2009-11-10T22:43:5"]}]`
	testutil.CheckEqualJSON(t, got, want, got)

	got = matchedTimePatternFiles(getLogFiles("ls /tmp/log/*/*/*.log"), 1003, time.Time{})
	want = `[{"File":"2009-11-10_23.log","Targets":["node01","namespace01"],"TimePatterns":["2009-11-10T23:00:00"]},{"File":"2009-11-10_22.log","Targets":["node01","node02","namespace01","namespace02"],"TimePatterns":["2009-11-10T22:43:17","2009-11-10T22:43:18","2009-11-10T22:43:19","2009-11-10T22:5","2009-11-10T22:44","2009-11-10T22:45","2009-11-10T22:46","2009-11-10T22:47","2009-11-10T22:48","2009-11-10T22:49","2009-11-10T22:43:2","2009-11-10T22:43:3","2009-11-10T22:43:4","2009-11-10T22:43:5"]}]`
	testutil.CheckEqualJSON(t, got, want, got)

	got = matchedTimePatternFiles(getLogFiles("ls /tmp/log/*/*/*.log"), 0, config.GetNow().Add(time.Duration(-3)*time.Second))
	want = `[{"File":"2009-11-10_22.log","Targets":["node01","node02","namespace01","namespace02"],"TimePatterns":["2009-11-10T22:59:50","2009-11-10T22:59:51","2009-11-10T22:59:52","2009-11-10T22:59:53","2009-11-10T22:59:54","2009-11-10T22:59:55","2009-11-10T22:59:56","2009-11-10T22:59:57","2009-11-10T22:0","2009-11-10T22:1","2009-11-10T22:2","2009-11-10T22:3","2009-11-10T22:4","2009-11-10T22:50","2009-11-10T22:51","2009-11-10T22:52","2009-11-10T22:53","2009-11-10T22:54","2009-11-10T22:55","2009-11-10T22:56","2009-11-10T22:57","2009-11-10T22:58","2009-11-10T22:59:0","2009-11-10T22:59:1","2009-11-10T22:59:2","2009-11-10T22:59:3","2009-11-10T22:59:4"]}]`
	testutil.CheckEqualJSON(t, got, want, got)

	got = matchedTimePatternFiles(getLogFiles("ls /tmp/log/*/*/*.log"), 6, config.GetNow().Add(time.Duration(-3)*time.Second))
	want = `[{"File":"2009-11-10_22.log","Targets":["node01","node02","namespace01","namespace02"],"TimePatterns":["2009-11-10T22:59:51","2009-11-10T22:59:52","2009-11-10T22:59:53","2009-11-10T22:59:54","2009-11-10T22:59:55","2009-11-10T22:59:56","2009-11-10T22:59:57"]}]`
	testutil.CheckEqualJSON(t, got, want, got)

	got = matchedTimePatternFiles(getLogFiles("ls /tmp/log/*/*/*.log"), 61*60, config.GetNow().Add(time.Duration(-31)*time.Second))
	want = `[{"File":"2009-11-10_22.log","Targets":["node01","node02","namespace01","namespace02"],"TimePatterns":["2009-11-10T22:0","2009-11-10T22:1","2009-11-10T22:2","2009-11-10T22:3","2009-11-10T22:4","2009-11-10T22:50","2009-11-10T22:51","2009-11-10T22:52","2009-11-10T22:53","2009-11-10T22:54","2009-11-10T22:55","2009-11-10T22:56","2009-11-10T22:57","2009-11-10T22:58","2009-11-10T22:59:0","2009-11-10T22:59:1","2009-11-10T22:59:2"]}]`
	testutil.CheckEqualJSON(t, got, want, got)
}

func Test_reduceTimePatterns(t *testing.T) {
	testutil.SetTestLogs()

	var got []string
	var want string

	got = reduceTimePatterns([]string{"2009-11-10T23:00:00"})
	want = `["2009-11-10T23:00:00"]`
	testutil.CheckEqualJSON(t, got, want)

	got = reduceTimePatterns([]string{"2009-11-10T23:00:00", "2009-11-10T22:59:59", "2009-11-10T22:59:58", "2009-11-10T22:59:57"})
	want = `["2009-11-10T23:00:00","2009-11-10T22:59:59","2009-11-10T22:59:58","2009-11-10T22:59:57"]`
	testutil.CheckEqualJSON(t, got, want)

	got = reduceTimePatterns([]string{"2009-11-10T22:59:59", "2009-11-10T22:59:58", "2009-11-10T22:59:57", "2009-11-10T22:59:56", "2009-11-10T22:59:55", "2009-11-10T22:59:54", "2009-11-10T22:59:53", "2009-11-10T22:59:52", "2009-11-10T22:59:51", "2009-11-10T22:59:50"})
	want = `["2009-11-10T22:59:5"]`
	testutil.CheckEqualJSON(t, got, want)
}

func podLogs(namespace, pod, container string, durationSeconds int, endTime time.Time, typ string, keyword string) (Result, error) {
	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp", "log")})
	logStore := LogStore{driver: d}

	result, err := logStore.GetLogs(LogSearch{
		LogType:       PodLog{Name: POD_TYPE},
		TargetPattern: PatternedString(namespace),
		PodSearchParams: PodSearchParams{
			Pod:       PatternedString(pod),
			Container: PatternedString(container),
		},
		Keyword:         keyword,
		DurationSeconds: durationSeconds,
		EndTime:         endTime,
		IsCounting:      typ == Count,
	})
	if err != nil {
		return Result{}, err
	}
	return result, nil
}

func nodeLogs(node, process string, durationSeconds int, endTime time.Time, typ string, keyword string) (Result, error) {
	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp", "log")})
	logStore := LogStore{driver: d}

	result, err := logStore.GetLogs(LogSearch{
		LogType:       NodeLog{Name: NODE_TYPE},
		TargetPattern: PatternedString(node),
		NodeSearchParams: NodeSearchParams{
			Process: PatternedString(process),
		},
		Keyword:         keyword,
		DurationSeconds: durationSeconds,
		EndTime:         endTime,
		IsCounting:      typ == Count,
	})
	if err != nil {
		return Result{}, err
	}
	return result, nil
}
