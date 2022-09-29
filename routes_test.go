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

func Test_Routes_Ping(t *testing.T) {
	var got string
	var want string

	got = GET("/ping", nil)
	want = `"{\"message\":\"pong\"}"`
	testutil.CheckEqualJSON(t, got, want)
}

func Test_Routes_Query(t *testing.T) {
	var got string
	var want string

	got = GET("/api/v1/query", Params{"query": ``})
	want = `"{\"error\":\"empty query\",\"status\":\"error\"}"`
	testutil.CheckEqualJSON(t, got, want)

	got = GET("/api/v1/query", Params{"query": `hello`})
	want = `"{\"error\":\"unknown log name\",\"status\":\"error\"}"`
	testutil.CheckEqualJSON(t, got, want)

	got = GET("/api/v1/query", Params{"query": `pod{namespace="namespace01"}`})
	want = `"{\"data\":{\"result\":[\"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum\",\"2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\",\"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\"],\"resultType\":\"logs\"},\"status\":\"success\"}"`
	testutil.CheckEqualJSON(t, got, want)

	got = GET("/api/v1/query", Params{"query": `pod{namespace="namespace02"}`})
	want = `"{\"data\":{\"result\":[\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\"],\"resultType\":\"logs\"},\"status\":\"success\"}"`
	testutil.CheckEqualJSON(t, got, want)

	got = GET("/api/v1/query", Params{"query": `pod{namespace="namespace03"}`})
	want = `"{\"data\":{\"result\":[],\"resultType\":\"logs\"},\"status\":\"success\"}"`
	testutil.CheckEqualJSON(t, got, want)

	got = GET("/api/v1/query", Params{"query": `pod{namespace="namespace01"}[2m]`})
	want = `"{\"data\":{\"result\":[\"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\",\"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\"],\"resultType\":\"logs\"},\"status\":\"success\"}"`
	testutil.CheckEqualJSON(t, got, want)

	got = GET("/api/v1/query", Params{"query": `count_over_time(pod{namespace="namespace01"}[2m])`})
	want = `"{\"data\":{\"result\":[{\"value\":6}],\"resultType\":\"vector\"},\"status\":\"success\"}"`
	testutil.CheckEqualJSON(t, got, want)

	got = GET("/api/v1/query", Params{"query": `count_over_time(pod{namespace="namespace01"}[2m])>0`})
	want = `"{\"data\":{\"result\":[{\"value\":1}],\"resultType\":\"vector\"},\"status\":\"success\"}"`
	testutil.CheckEqualJSON(t, got, want)

	got = GET("/api/v1/query", Params{"query": `count_over_time(pod{namespace="namespace03"}[2m])>0`})
	want = `"{\"data\":{\"result\":[],\"resultType\":\"vector\"},\"status\":\"success\"}"`
	testutil.CheckEqualJSON(t, got, want)
}

func Test_Routes_QueryRange(t *testing.T) {
	var got string
	var want string

	got = GET("/api/v1/query_range", Params{"query": ``})
	want = `"{\"error\":\"empty query\",\"status\":\"error\"}"`
	testutil.CheckEqualJSON(t, got, want)

	got = GET("/api/v1/query_range", Params{"query": `hello`})
	want = `"{\"error\":\"unknown log name\",\"status\":\"error\"}"`
	testutil.CheckEqualJSON(t, got, want)

	got = GET("/api/v1/query_range", Params{"query": `pod{namespace="namespace01"}`})
	want = `"{\"data\":{\"result\":[\"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum\",\"2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\",\"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\"],\"resultType\":\"logs\"},\"status\":\"success\"}"`
	testutil.CheckEqualJSON(t, got, want)

	got = GET("/api/v1/query_range", Params{"query": `pod{namespace="namespace02"}`})
	want = `"{\"data\":{\"result\":[\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\"],\"resultType\":\"logs\"},\"status\":\"success\"}"`
	testutil.CheckEqualJSON(t, got, want)

	got = GET("/api/v1/query_range", Params{"query": `pod{namespace="namespace03"}`})
	want = `"{\"data\":{\"result\":[],\"resultType\":\"logs\"},\"status\":\"success\"}"`
	testutil.CheckEqualJSON(t, got, want)

	got = GET("/api/v1/query_range", Params{"query": `pod{namespace="namespace01"}[2m]`})
	want = `"{\"data\":{\"result\":[\"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar\",\"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\",\"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\",\"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\"],\"resultType\":\"logs\"},\"status\":\"success\"}"`
	testutil.CheckEqualJSON(t, got, want)

	got = GET("/api/v1/query_range", Params{"query": `count_over_time(pod{namespace="namespace01"}[2m])`})
	want = `"{\"data\":{\"result\":[{\"value\":6}],\"resultType\":\"vector\"},\"status\":\"success\"}"`
	testutil.CheckEqualJSON(t, got, want)

	got = GET("/api/v1/query_range", Params{"query": `count_over_time(pod{namespace="namespace01"}[2m])>0`})
	want = `"{\"data\":{\"result\":[{\"value\":1}],\"resultType\":\"vector\"},\"status\":\"success\"}"`
	testutil.CheckEqualJSON(t, got, want)

	got = GET("/api/v1/query_range", Params{"query": `count_over_time(pod{namespace="namespace03"}[2m])>0`})
	want = `"{\"data\":{\"result\":[],\"resultType\":\"vector\"},\"status\":\"success\"}"`
	testutil.CheckEqualJSON(t, got, want)
}
