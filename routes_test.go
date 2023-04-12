package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/lethe/testutil"
	"github.com/stretchr/testify/assert"
)

var r *gin.Engine
var w *httptest.ResponseRecorder

type Params map[string]string

func init() {
	testutil.Init()
	testutil.SetTestLogFiles()
	r = NewRouter()
	w = httptest.NewRecorder()
}

func GET(url string, params Params) string {
	req, _ := http.NewRequest("GET", url, nil)

	if params != nil {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	r.ServeHTTP(w, req)
	body, _ := ioutil.ReadAll(w.Body)
	return string(body)
}

func Test_Routes_Query(t *testing.T) {
	prefix := "01_"
	tests := map[string]struct {
		url   string
		param Params
		want  string
	}{
		// modify test case.. for supoorting sort by time
		"01_ping": {url: "/ping", param: nil,
			want: "{\"message\":\"pong\"}"},
		"02_empty_query": {url: "/api/v1/query", param: Params{"query": ``},
			want: "{\"error\":\"empty query\",\"status\":\"error\"}"},
		"03_unknown_query": {url: "/api/v1/query", param: Params{"query": `hello`},
			want: "{\"error\":\"unknown log name\",\"status\":\"error\"}"},
		"04_query_namespace01": {url: "/api/v1/query", param: Params{"query": `pod{namespace="namespace01"}`},
			want: "{\"data\":{\"result\":[\"2009-11-10T21:00:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx]  hello world\",\"2009-11-10T21:01:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx]  hello world\",\"2009-11-10T21:02:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx]  hello world\",\"2009-11-10T22:56:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:56:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:56:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:57:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:57:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:57:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  lerom from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:59:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx]  lerom ipsum\",\"2009-11-10T22:59:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx]  hello world\"],\"resultType\":\"logs\"},\"status\":\"success\"}"},
		"05_query_namespace02": {url: "/api/v1/query", param: Params{"query": `pod{namespace="namespace02"}`},
			want: "{\"data\":{\"result\":[\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx]  hello world\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx]  lerom ipsum\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx]  hello world\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar]  lerom from sidecar\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\"],\"resultType\":\"logs\"},\"status\":\"success\"}"},
		"06_query_namespace03": {url: "/api/v1/query", param: Params{"query": `pod{namespace="namespace03"}`},
			want: "{\"data\":{\"result\":null,\"resultType\":\"logs\"},\"status\":\"success\"}"},
		"07_query_with_duration": {url: "/api/v1/query", param: Params{"query": `pod{namespace="namespace01"}[2m]`},
			want: "{\"data\":{\"result\":[\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  lerom from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:59:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx]  lerom ipsum\",\"2009-11-10T22:59:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx]  hello world\"],\"resultType\":\"logs\"},\"status\":\"success\"}"},
		"08_query_with_count_over_time_function": {url: "/api/v1/query", param: Params{"query": `count_over_time(pod{namespace="namespace01"}[2m])`},
			want: "{\"data\":{\"result\":[{\"value\":5}],\"resultType\":\"vector\"},\"status\":\"success\"}"},
		"09_query_with_count_over_time_function_and_binary_operator": {url: "/api/v1/query", param: Params{"query": `count_over_time(pod{namespace="namespace02"}[2m])>0`},
			want: "{\"data\":{\"result\":[{\"value\":1}],\"resultType\":\"vector\"},\"status\":\"success\"}"},
		"10_query_namespace03_with_count_over_time_function_and_binary_operator": {url: "/api/v1/query", param: Params{"query": `count_over_time(pod{namespace="namespace03"}[2m])>0`},
			want: "{\"data\":{\"result\":[],\"resultType\":\"vector\"},\"status\":\"success\"}"},
	}

	for name, tt := range tests {
		t.Run(prefix+name, func(subt *testing.T) {
			got := GET(tt.url, tt.param)
			assert.Equal(subt, tt.want, got)
		})
	}
}

func Test_Routes_QueryRange_Without_StartEnd(t *testing.T) {
	prefix := "02_"
	tests := map[string]struct {
		url   string
		param Params
		want  string
	}{
		// modify test case.. for supoorting sort by time
		"01_empty_query": {url: "/api/v1/query_range", param: Params{"query": ``},
			want: "{\"error\":\"empty query\",\"status\":\"error\"}"},
		"02_unknown_query": {url: "/api/v1/query_range", param: Params{"query": `hello`},
			want: "{\"error\":\"empty query\",\"status\":\"error\"}"},
		"03_query_namespace01": {url: "/api/v1/query_range", param: Params{"query": `pod{namespace="namespace01"}`},
			want: "{\"error\":\"empty query\",\"status\":\"error\"}"},
		"04_query_namespace02": {url: "/api/v1/query_range", param: Params{"query": `pod{namespace="namespace02"}`},
			want: "{\"error\":\"empty query\",\"status\":\"error\"}"},
		"05_query_namespace03": {url: "/api/v1/query_range", param: Params{"query": `pod{namespace="namespace03"}`},
			want: "{\"error\":\"empty query\",\"status\":\"error\"}"},

		"06_query_with_duration": {url: "/api/v1/query_range", param: Params{"query": `pod{namespace="namespace02"}[2m]`},
			want: "{\"error\":\"empty query\",\"status\":\"error\"}"},
		"07_query_with_count_over_time_function": {url: "/api/v1/query_range", param: Params{"query": `count_over_time(pod{namespace="namespace01"}[2m])`},
			want: "{\"error\":\"empty query\",\"status\":\"error\"}"},
		"08_query_with_count_over_time_function_and_binary_operator": {url: "/api/v1/query_range", param: Params{"query": `count_over_time(pod{namespace="namespace01"}[2m])>0`},
			want: "{\"error\":\"empty query\",\"status\":\"error\"}"},
		"09_query_namespace03_with_count_over_time_function_and_binary_operator": {url: "/api/v1/query_range", param: Params{"query": `count_over_time(pod{namespace="namespace03"}[2m])>0`},
			want: "{\"error\":\"empty query\",\"status\":\"error\"}"},
	}

	for name, tt := range tests {
		t.Run(prefix+name, func(subt *testing.T) {
			t.Log(tt.url, tt.param)
			got := GET(tt.url, tt.param)
			assert.Equal(subt, tt.want, got)
		})
	}
}

func Test_Routes_QueryRange_With_StartEnd(t *testing.T) {
	startTime, _ := time.Parse(time.RFC3339, "2009-11-10T22:57:00Z")
	endTime, _ := time.Parse(time.RFC3339, "2009-11-10T22:58:00Z")
	start := fmt.Sprint(startTime.Unix())
	end := fmt.Sprint(endTime.Unix())

	prefix := "03_"
	tests := map[string]struct {
		url   string
		param Params
		want  string
	}{
		// modify test case.. for supoorting sort by time
		"01_empty_query": {url: "/api/v1/query_range", param: Params{"start": start, "end": end, "query": ``},
			want: "{\"error\":\"empty query\",\"status\":\"error\"}"},
		"02_unknown_query": {url: "/api/v1/query_range", param: Params{"start": start, "end": end, "query": `hello`},
			want: "{\"error\":\"unknown log name\",\"status\":\"error\"}"},
		"03_query_namespace01": {url: "/api/v1/query_range", param: Params{"start": start, "end": end, "query": `pod{namespace="namespace01"}`},
			want: "{\"data\":{\"result\":[\"2009-11-10T22:57:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:57:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:57:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  lerom from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\"],\"resultType\":\"logs\"},\"status\":\"success\"}"},
		"04_query_namespace02": {url: "/api/v1/query_range", param: Params{"start": start, "end": end, "query": `pod{namespace="namespace02"}`},
			want: "{\"data\":{\"result\":[\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx]  hello world\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx]  lerom ipsum\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx]  hello world\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar]  lerom from sidecar\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\"],\"resultType\":\"logs\"},\"status\":\"success\"}"},
		"05_query_namespace03": {url: "/api/v1/query_range", param: Params{"start": start, "end": end, "query": `pod{namespace="namespace03"}`},
			want: "{\"data\":{\"result\":null,\"resultType\":\"logs\"},\"status\":\"success\"}"},

		"06_query_with_duration": {url: "/api/v1/query_range", param: Params{"start": start, "end": end, "query": `pod{namespace="namespace01"}[2m]`},
			want: "{\"data\":{\"result\":[\"2009-11-10T22:57:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:57:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:57:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  lerom from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\"],\"resultType\":\"logs\"},\"status\":\"success\"}"},
		"07_query_with_count_over_time_function": {url: "/api/v1/query_range", param: Params{"start": start, "end": end, "query": `count_over_time(pod{namespace="namespace01"}[2m])`},
			want: "{\"data\":{\"result\":[{\"value\":6}],\"resultType\":\"vector\"},\"status\":\"success\"}"},
		"08_query_with_count_over_time_function_and_binary_operator": {url: "/api/v1/query_range", param: Params{"start": start, "end": end, "query": `count_over_time(pod{namespace="namespace01"}[2m])>0`},
			want: "{\"data\":{\"result\":[{\"value\":1}],\"resultType\":\"vector\"},\"status\":\"success\"}"},
		"09_query_namespace03_with_count_over_time_function_and_binary_operator": {url: "/api/v1/query_range", param: Params{"start": start, "end": end, "query": `count_over_time(pod{namespace="namespace03"}[2m])>0`},
			want: "{\"data\":{\"result\":[],\"resultType\":\"vector\"},\"status\":\"success\"}"},
	}

	for name, tt := range tests {
		t.Run(prefix+name, func(subt *testing.T) {
			t.Log(tt.url, tt.param)
			got := GET(tt.url, tt.param)
			assert.Equal(subt, tt.want, got)
		})
	}
}
