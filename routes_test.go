package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/lethe/testutil"
)

var r *gin.Engine
var w *httptest.ResponseRecorder

type Params map[string]string

func init() {
	testutil.Init()
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

	tests := map[string]struct {
		url   string
		param Params
		want  string
	}{
		// modify test case.. for supoorting sort by time
		"ping":              {url: "/ping", param: nil, want: `"{\"message\":\"pong\"}"`},
		"empty query":       {url: "/api/v1/query", param: Params{"query": ``}, want: `"{\"error\":\"empty query\",\"status\":\"error\"}"`},
		"unknown query":     {url: "/api/v1/query", param: Params{"query": `hello`}, want: `"{\"error\":\"unknown log name\",\"status\":\"error\"}"`},
		"query namespace01": {url: "/api/v1/query", param: Params{"query": `pod{namespace="namespace01"}`}, want: `"{\"data\":{\"result\":[\"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum\",\"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\",\"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\"],\"resultType\":\"logs\"},\"status\":\"success\"}"`},
		"query namespace02": {url: "/api/v1/query", param: Params{"query": `pod{namespace="namespace02"}`}, want: `"{\"data\":{\"result\":[\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\"],\"resultType\":\"logs\"},\"status\":\"success\"}"`},
		"query namespace03": {url: "/api/v1/query", param: Params{"query": `pod{namespace="namespace03"}`}, want: `"{\"data\":{\"result\":[],\"resultType\":\"logs\"},\"status\":\"success\"}"`},

		"query with duration":                                                 {url: "/api/v1/query", param: Params{"query": `pod{namespace="namespace01"}[2m]`}, want: `"{\"data\":{\"result\":[\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum\",\"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\",\"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\"],\"resultType\":\"logs\"},\"status\":\"success\"}"`},
		"query with count_over_time Function":                                 {url: "/api/v1/query", param: Params{"query": `count_over_time(pod{namespace="namespace01"}[2m])`}, want: `"{\"data\":{\"result\":[{\"value\":6}],\"resultType\":\"vector\"},\"status\":\"success\"}"`},
		"query with count_over_time Function and binary operator":             {url: "/api/v1/query", param: Params{"query": `count_over_time(pod{namespace="namespace01"}[2m])>0`}, want: `"{\"data\":{\"result\":[{\"value\":1}],\"resultType\":\"vector\"},\"status\":\"success\"}"`},
		"query namespace03 with count_over_time Function and binary operator": {url: "/api/v1/query", param: Params{"query": `count_over_time(pod{namespace="namespace03"}[2m])>0`}, want: `"{\"data\":{\"result\":[],\"resultType\":\"vector\"},\"status\":\"success\"}"`},
	}

	for name, tt := range tests {
		t.Run(name, func(subt *testing.T) {
			got := GET(tt.url, tt.param)
			testutil.CheckEqualJSON(subt, got, tt.want)
		})
	}
}

func Test_Routes_Registry_Sample(t *testing.T) {

	tests := map[string]struct {
		url   string
		param Params
		want  string
	}{
		// modify test case.. for supoorting sort by time
		"ping": {url: "/ping", param: nil, want: `"{\"message\":\"pong\"}"`},
		"query registry": {url: "/api/v1/query_range", param: Params{
			"query": `pod{namespace="registry", pod="registry-manager-.*"}`,
			"start": "1673222724.056",
			"end":   "1673223024.056",
		}, want: `"{\"data\":{\"result\":[\"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum\",\"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\",\"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\"],\"resultType\":\"logs\"},\"status\":\"success\"}"`},
	}

	for name, tt := range tests {
		t.Run(name, func(subt *testing.T) {
			got := GET(tt.url, tt.param)
			testutil.CheckEqualJSON(subt, got, tt.want)
		})
	}
}

//fail
func Test_Routes_QueryRange(t *testing.T) {
	tests := map[string]struct {
		url   string
		param Params
		want  string
	}{
		// modify test case.. for supoorting sort by time
		"empty query":       {url: "/api/v1/query_range", param: Params{"query": ``}, want: `"{\"error\":\"empty query\",\"status\":\"error\"}"`},
		"unknown query":     {url: "/api/v1/query_range", param: Params{"query": `hello`}, want: `"{\"error\":\"unknown log name\",\"status\":\"error\"}"`},
		"query namespace01": {url: "/api/v1/query_range", param: Params{"query": `pod{namespace="namespace01"}`}, want: `"{\"data\":{\"result\":[\"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum\",\"2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\",\"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\"],\"resultType\":\"logs\"},\"status\":\"success\"}"`},
		"query namespace02": {url: "/api/v1/query_range", param: Params{"query": `pod{namespace="namespace02"}`}, want: `"{\"data\":{\"result\":[\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\"],\"resultType\":\"logs\"},\"status\":\"success\"}"`},
		"query namespace03": {url: "/api/v1/query_range", param: Params{"query": `pod{namespace="namespace03"}`}, want: `"{\"data\":{\"result\":[],\"resultType\":\"logs\"},\"status\":\"success\"}"`},

		"query with duration":                                                 {url: "/api/v1/query_range", param: Params{"query": `pod{namespace="namespace01"}[2m]`}, want: `"{\"data\":{\"result\":[\"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\",\"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\"],\"resultType\":\"logs\"},\"status\":\"success\"}"`},
		"query with count_over_time Function":                                 {url: "/api/v1/query_range", param: Params{"query": `count_over_time(pod{namespace="namespace01"}[2m])`}, want: `"{\"data\":{\"result\":[{\"value\":6}],\"resultType\":\"vector\"},\"status\":\"success\"}"`},
		"query with count_over_time Function and binary operator":             {url: "/api/v1/query_range", param: Params{"query": `count_over_time(pod{namespace="namespace01"}[2m])>0`}, want: `"{\"data\":{\"result\":[{\"value\":1}],\"resultType\":\"vector\"},\"status\":\"success\"}"`},
		"query namespace03 with count_over_time Function and binary operator": {url: "/api/v1/query_range", param: Params{"query": `count_over_time(pod{namespace="namespace03"}[2m])>0`}, want: `"{\"data\":{\"result\":[],\"resultType\":\"vector\"},\"status\":\"success\"}"`},
	}

	for name, tt := range tests {
		t.Run(name, func(subt *testing.T) {
			got := GET(tt.url, tt.param)
			testutil.CheckEqualJSON(subt, got, tt.want)
		})
	}
}
