package letheql

import (
	"github.com/kuoss/lethe/logs/logStore"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kuoss/lethe/clock"
	_ "github.com/kuoss/lethe/storage/driver/filesystem"
	"github.com/kuoss/lethe/testutil"
)

func init() {
	testutil.Init()
	testutil.SetTestLogFiles()
}

func Test_Query_Success(t *testing.T) {

	tests := map[string]struct {
		query string
		want  QueryData
	}{
		"01": {query: `pod{namespace="namespace01"}`,
			want: QueryData{ResultType: "logs", Logs: []logStore.LogLine{logStore.PodLog{Name: "", Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
		"02": {query: `pod{namespace="not-exists"}`,
			want: QueryData{ResultType: ValueTypeLogs, Logs: []logStore.LogLine{}}},
		"03": {query: `pod{namespace="namespace01",pod="nginx"}`,
			want: QueryData{ResultType: "logs", Logs: []logStore.LogLine{}, Scalar: 0}},
		"04": {query: `pod{namespace="namespace01",pod="nginx-*"}`,
			want: QueryData{ResultType: "logs", Logs: []logStore.LogLine{logStore.PodLog{Name: "", Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
		"05": {query: `pod{namespace="namespace01",container="nginx"}`,
			want: QueryData{ResultType: "logs", Logs: []logStore.LogLine{logStore.PodLog{Name: "", Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
		"06": {query: `pod{namespace="namespace*",container="nginx"}`,
			want: QueryData{ResultType: "logs", Logs: []logStore.LogLine{logStore.PodLog{Name: "", Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
		"07": {query: `pod{namespace="namespace01",pod="nginx-*"}[3m]`,
			want: QueryData{ResultType: "logs", Logs: []logStore.LogLine{logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
		"08": {query: `pod`,
			want: QueryData{ResultType: "logs", Logs: []logStore.LogLine{logStore.PodLog{Name: "", Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
		"09": {query: `pod{}`,
			want: QueryData{ResultType: "logs", Logs: []logStore.LogLine{logStore.PodLog{Name: "", Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
		"10": {query: `count_over_time(pod{namespace="namespace01",pod="nginx-*"}[3m])`,
			want: QueryData{ResultType: "scalar", Logs: nil, Scalar: 5}},
		"11": {query: `1`,
			want: QueryData{ResultType: ValueTypeScalar, Logs: nil, Scalar: 1}},
		"12": {query: `count_over_time(pod{}[3m])`,
			want: QueryData{ResultType: "scalar", Logs: nil, Scalar: 20}},
		"13": {query: `count_over_time(pod{}[3m]) > 10`,
			want: QueryData{ResultType: ValueTypeScalar, Logs: nil, Scalar: 1}},
		"14": {query: `count_over_time(pod{}[3m]) < 10`,
			want: QueryData{ResultType: ValueTypeScalar, Logs: nil, Scalar: 0}},
		"15": {query: `count_over_time(pod{}[3m]) == 21`,
			want: QueryData{ResultType: ValueTypeScalar, Logs: nil, Scalar: 0}},

		// todo
		// scalar operator same with filter ( != )
		// "operator != count_over_time_with_duration": {query: `count_over_time(pod{}[3m]) != 21`, want: QueryData{
		// 	ResultType: ValueTypeScalar,
		// 	Logs:       nil,
		// 	Scalar:     0,
		// }},

		"16": {query: `pod{namespace="namespace01"} |= hello`,
			want: QueryData{ResultType: "logs", Logs: []logStore.LogLine{logStore.PodLog{Name: "", Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
		"17": {query: `pod{namespace="namespace01"} != hello`,
			want: QueryData{ResultType: "logs", Logs: []logStore.LogLine{logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}}}},
		"18": {query: `pod{namespace="namespace01"} |~ (.ro.*o)`,
			want: QueryData{ResultType: ValueTypeLogs, Logs: []logStore.LogLine{logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}}}},
		"19": {query: `pod{namespace="namespace01"} !~ (.le)`,
			want: QueryData{ResultType: "logs", Logs: []logStore.LogLine{logStore.PodLog{Name: "", Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
	}

	for name, tt := range tests {
		t.Run(name+" "+tt.query, func(subt *testing.T) {
			got, err := ProcQuery(tt.query, TimeRange{})
			if err != nil {
				subt.Fatalf("query: %s error: %s", name, tt.query)
			}
			assert.Equal(subt, tt.want, got)
		})
	}
}

func ago(m int) time.Time {
	return clock.GetNow().Add(time.Duration(-m) * time.Minute)
}
func Test_QueryWithTimeRange(t *testing.T) {

	now := clock.GetNow()

	tests := map[string]struct {
		query     string
		timeRange TimeRange
		want      QueryData
	}{
		// modify test case.. for supoorting sort by time
		"01": {
			query: `pod{namespace="namespace01",pod="nginx-*"}[3m]`, timeRange: TimeRange{},
			want: QueryData{ResultType: "logs", Logs: []logStore.LogLine{logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}},
		},
		"02": {
			query: `pod{namespace="namespace01",pod="nginx-*"}`, timeRange: TimeRange{Start: ago(999999), End: now},
			want: QueryData{ResultType: "logs", Logs: []logStore.LogLine{logStore.PodLog{Name: "", Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}},
		},
		"03": {
			query: `pod{namespace="namespace01",pod="nginx-*"}`, timeRange: TimeRange{Start: ago(1), End: now},
			want: QueryData{ResultType: "logs", Logs: []logStore.LogLine{logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logStore.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}},
		},
		"04": {query: `count_over_time(pod{namespace="*"})`, timeRange: TimeRange{Start: ago(1), End: now},
			want: QueryData{ResultType: "scalar", Logs: nil, Scalar: 2}},
		"05": {query: `count_over_time(pod{namespace="*"})`, timeRange: TimeRange{Start: ago(2), End: now},
			want: QueryData{ResultType: "scalar", Logs: nil, Scalar: 17}},
		"06": {query: `count_over_time(pod{namespace="*"})`, timeRange: TimeRange{Start: ago(3), End: now},
			want: QueryData{ResultType: ValueTypeScalar, Logs: nil, Scalar: 20}},
	}

	for name, tt := range tests {
		t.Run(name+"_"+tt.query, func(subt *testing.T) {
			got, err := ProcQuery(tt.query, tt.timeRange)
			if err != nil {
				subt.Fatalf("query: %s error: %s", name, tt.query)
			}
			assert.Equal(subt, tt.want, got)
		})
	}
}

func Test_QueryFail(t *testing.T) {

	var now = clock.GetNow()
	tests := map[string]struct {
		query     string
		timeRange TimeRange
		want      string
	}{
		// modify test case.. for supoorting sort by time
		"01": {query: `pod{namespace=""}`, timeRange: TimeRange{},
			want: "namespace value cannot be empty"},
		"02": {query: `pod{foo=""}`, timeRange: TimeRange{},
			want: "unknown label foo"},
		"03": {query: `{namespace="hello"}`, timeRange: TimeRange{},
			want: "a log name must be specified"},
		"04": {query: `count_over_time(pod{namespace="*"})`, timeRange: TimeRange{Start: ago(0), End: now},
			want: "end time and start time are the same"},
	}

	for name, tt := range tests {
		t.Run(name+"_"+tt.query, func(subt *testing.T) {
			_, err := ProcQuery(tt.query, tt.timeRange)
			if assert.Error(subt, err) {
				assert.Equal(subt, tt.want, err.Error())
			}
		})
	}
}

func TestNewQuery(t *testing.T) {

	e := &Engine{}

	tests := map[string]struct {
		query string
		want  *query
	}{
		`pod metric with namespace label matcher`: {
			query: `pod{namespace="namespace01"}`,
			want: &query{
				q:       `pod{namespace="namespace01"}`,
				filter:  nil,
				keyword: "",
				engine:  e,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name+"_"+tt.query, func(subt *testing.T) {
			query, err := e.newQuery(tt.query)
			if err != nil {
				subt.Fatalf("%s test failed. build new Query with err: %v", tt.query, err)
			}
			assert.Equal(subt, query, tt.want)
		})
	}
}
