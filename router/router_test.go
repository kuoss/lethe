package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kuoss/lethe/storage/logservice/logmodel"
	"github.com/stretchr/testify/assert"
)

var (
	w *httptest.ResponseRecorder = httptest.NewRecorder()
)

type Params map[string]string

func TestNew(t *testing.T) {
	assert.NotEmpty(t, router1)
}

func GET(url string, params Params) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	if params != nil {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	router1.ServeHTTP(w, req)
	body, _ := io.ReadAll(w.Body)
	return string(body)
}

func Test_Routes_Query(t *testing.T) {
	testCases := []struct {
		url   string
		param Params
		want  string
	}{
		{
			"/ping", nil,
			"404 page not found",
		},
		{
			"/api/v1/query", Params{"query": ``},
			"{\"error\":\"empty query\",\"status\":\"error\"}",
		},
		{
			"/api/v1/query", Params{"query": `hello`},
			"{\"error\":\"unknown log name\",\"status\":\"error\"}",
		},
		{
			"/api/v1/query", Params{"query": `pod{namespace="namespace01"}`},
			"{\"data\":{\"result\":[\"2009-11-10T21:00:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx]  hello world\",\"2009-11-10T21:01:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx]  hello world\",\"2009-11-10T21:02:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx]  hello world\",\"2009-11-10T22:56:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:56:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:56:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:57:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:57:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:57:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  lerom from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:59:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx]  lerom ipsum\",\"2009-11-10T22:59:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx]  hello world\"],\"resultType\":\"logs\"},\"status\":\"success\"}",
		},
		{
			"/api/v1/query", Params{"query": `pod{namespace="namespace02"}`},
			"{\"data\":{\"result\":[\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx]  hello world\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx]  lerom ipsum\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx]  hello world\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar]  lerom from sidecar\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\"],\"resultType\":\"logs\"},\"status\":\"success\"}",
		},
		{
			"/api/v1/query", Params{"query": `pod{namespace="namespace03"}`},
			"{\"data\":{\"result\":null,\"resultType\":\"logs\"},\"status\":\"success\"}",
		},
		{
			"/api/v1/query", Params{"query": `pod{namespace="namespace01"}[2m]`},
			"{\"data\":{\"result\":[\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  lerom from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:59:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx]  lerom ipsum\",\"2009-11-10T22:59:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx]  hello world\"],\"resultType\":\"logs\"},\"status\":\"success\"}",
		},
		{
			"/api/v1/query", Params{"query": `count_over_time(pod{namespace="namespace01"}[2m])`},
			"{\"data\":{\"result\":[{\"value\":5}],\"resultType\":\"vector\"},\"status\":\"success\"}",
		},
		{
			"/api/v1/query", Params{"query": `count_over_time(pod{namespace="namespace02"}[2m])>0`},
			"{\"data\":{\"result\":[{\"value\":1}],\"resultType\":\"vector\"},\"status\":\"success\"}",
		},
		{
			"/api/v1/query", Params{"query": `count_over_time(pod{namespace="namespace03"}[2m])>0`},
			"{\"data\":{\"result\":[],\"resultType\":\"vector\"},\"status\":\"success\"}",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d_%s_%s", i, tc.url, tc.param), func(t *testing.T) {
			got := GET(tc.url, tc.param)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_Routes_QueryRange_Without_StartEnd(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		param Params
		want  string
	}{
		{"01_empty_query", "/api/v1/query_range",
			Params{"query": ``},
			"{\"error\":\"empty query\",\"status\":\"error\"}"},
		{"02_unknown_query", "/api/v1/query_range",
			Params{"query": `hello`},
			"{\"error\":\"empty query\",\"status\":\"error\"}"},
		{"03_query_namespace01", "/api/v1/query_range",
			Params{"query": `pod{namespace="namespace01"}`},
			"{\"error\":\"empty query\",\"status\":\"error\"}"},
		{"04_query_namespace02", "/api/v1/query_range",
			Params{"query": `pod{namespace="namespace02"}`},
			"{\"error\":\"empty query\",\"status\":\"error\"}"},
		{"05_query_namespace03", "/api/v1/query_range",
			Params{"query": `pod{namespace="namespace03"}`},
			"{\"error\":\"empty query\",\"status\":\"error\"}"},

		{"06_query_with_duration", "/api/v1/query_range",
			Params{"query": `pod{namespace="namespace02"}[2m]`},
			"{\"error\":\"empty query\",\"status\":\"error\"}"},
		{"07_query_with_count_over_time_function", "/api/v1/query_range",
			Params{"query": `count_over_time(pod{namespace="namespace01"}[2m])`},
			"{\"error\":\"empty query\",\"status\":\"error\"}"},
		{"08_query_with_count_over_time_function_and_binary_operator", "/api/v1/query_range",
			Params{"query": `count_over_time(pod{namespace="namespace01"}[2m])>0`},
			"{\"error\":\"empty query\",\"status\":\"error\"}"},
		{"09_query_namespace03_with_count_over_time_function_and_binary_operator", "/api/v1/query_range",
			Params{"query": `count_over_time(pod{namespace="namespace03"}[2m])>0`},
			"{\"error\":\"empty query\",\"status\":\"error\"}"},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d_%s", i, tc.name), func(t *testing.T) {
			t.Log(tc.url, tc.param)
			got := GET(tc.url, tc.param)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_Routes_QueryRange_With_StartEnd(t *testing.T) {
	startTime, _ := time.Parse(time.RFC3339, "2009-11-10T22:57:00Z")
	endTime, _ := time.Parse(time.RFC3339, "2009-11-10T22:58:00Z")
	start := fmt.Sprint(startTime.Unix())
	end := fmt.Sprint(endTime.Unix())

	testCases := []struct {
		name  string
		url   string
		param Params
		want  string
	}{
		// modify test case.. for supoorting sort by time
		{"01_empty_query", "/api/v1/query_range",
			Params{"start": start, "end": end, "query": ``},
			"{\"error\":\"empty query\",\"status\":\"error\"}"},
		{"02_unknown_query", "/api/v1/query_range",
			Params{"start": start, "end": end, "query": `hello`},
			"{\"error\":\"unknown log name\",\"status\":\"error\"}"},
		{"03_query_namespace01", "/api/v1/query_range",
			Params{"start": start, "end": end, "query": `pod{namespace="namespace01"}`},
			"{\"data\":{\"result\":[\"2009-11-10T22:57:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:57:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:57:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  lerom from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\"],\"resultType\":\"logs\"},\"status\":\"success\"}"},
		{"04_query_namespace02", "/api/v1/query_range",
			Params{"start": start, "end": end, "query": `pod{namespace="namespace02"}`},
			"{\"data\":{\"result\":[\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx]  hello world\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx]  lerom ipsum\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx]  hello world\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar]  lerom from sidecar\",\"2009-11-10T22:58:00Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace02|apache-75675f5897-7ci7o|httpd]  hello from sidecar\"],\"resultType\":\"logs\"},\"status\":\"success\"}"},
		{"05_query_namespace03", "/api/v1/query_range",
			Params{"start": start, "end": end, "query": `pod{namespace="namespace03"}`},
			"{\"data\":{\"result\":null,\"resultType\":\"logs\"},\"status\":\"success\"}"},

		{"06_query_with_duration", "/api/v1/query_range",
			Params{"start": start, "end": end, "query": `pod{namespace="namespace01"}[2m]`},
			"{\"data\":{\"result\":[\"2009-11-10T22:57:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:57:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:57:00Z[namespace01|apache-75675f5897-7ci7o|httpd]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  lerom from sidecar\",\"2009-11-10T22:58:00Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar]  hello from sidecar\"],\"resultType\":\"logs\"},\"status\":\"success\"}"},
		{"07_query_with_count_over_time_function", "/api/v1/query_range",
			Params{"start": start, "end": end, "query": `count_over_time(pod{namespace="namespace01"}[2m])`},
			"{\"data\":{\"result\":[{\"value\":6}],\"resultType\":\"vector\"},\"status\":\"success\"}"},
		{"08_query_with_count_over_time_function_and_binary_operator", "/api/v1/query_range",
			Params{"start": start, "end": end, "query": `count_over_time(pod{namespace="namespace01"}[2m])>0`},
			"{\"data\":{\"result\":[{\"value\":1}],\"resultType\":\"vector\"},\"status\":\"success\"}"},
		{"09_query_namespace03_with_count_over_time_function_and_binary_operator", "/api/v1/query_range",
			Params{"start": start, "end": end, "query": `count_over_time(pod{namespace="namespace03"}[2m])>0`},
			"{\"data\":{\"result\":[],\"resultType\":\"vector\"},\"status\":\"success\"}"},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d_%s", i, tc.name), func(t *testing.T) {
			t.Log(tc.url, tc.param)
			got := GET(tc.url, tc.param)
			assert.Equal(t, tc.want, got)
		})
	}
}

type TempDTO struct {
	Status string `json:"status,omitempty"`
	Data   struct {
		ResultType string            `json:"resultType,omitempty"`
		Result     []logmodel.PodLog `json:"result,omitempty"`
	} `json:"data"`
}

func Test_Routes_Query_json(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		param Params
		want  TempDTO
	}{
		// modify test case.. for supoorting sort by time
		{"04_query_namespace01", "/api/v1/query",
			Params{"query": `pod{namespace="namespace01"}`, "logFormat": "json"},
			TempDTO{Status: "success", Data: struct {
				ResultType string            "json:\"resultType,omitempty\""
				Result     []logmodel.PodLog "json:\"result,omitempty\""
			}{ResultType: "logs", Result: []logmodel.PodLog{{Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, {Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, {Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, {Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, {Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, {Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, {Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, {Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, {Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, {Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, {Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, {Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, {Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, {Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}}},
		{"05_query_namespace02", "/api/v1/query",
			Params{"query": `pod{namespace="namespace02"}`, "logFormat": "json"},
			TempDTO{Status: "success", Data: struct {
				ResultType string            "json:\"resultType,omitempty\""
				Result     []logmodel.PodLog "json:\"result,omitempty\""
			}{ResultType: "logs", Result: []logmodel.PodLog{{Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, {Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, {Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, {Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, {Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, {Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, {Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, {Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, {Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, {Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, {Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, {Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}}}}},
		{"06_query_namespace03", "/api/v1/query",
			Params{"query": `pod{namespace="namespace03"}`, "logFormat": "json"},
			TempDTO{Status: "success", Data: struct {
				ResultType string            "json:\"resultType,omitempty\""
				Result     []logmodel.PodLog "json:\"result,omitempty\""
			}{ResultType: "logs", Result: []logmodel.PodLog{}}}},
		{"07_query_with_duration", "/api/v1/query",
			Params{"query": `pod{namespace="namespace01"}[2m]`, "logFormat": "json"},
			TempDTO{Status: "success", Data: struct {
				ResultType string            "json:\"resultType,omitempty\""
				Result     []logmodel.PodLog "json:\"result,omitempty\""
			}{ResultType: "logs", Result: []logmodel.PodLog{{Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, {Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, {Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, {Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, {Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}}},
		{"08_query_with_count_over_time_function", "/api/v1/query",
			Params{"query": `count_over_time(pod{namespace="namespace01"}[2m])`, "logFormat": "json"},
			TempDTO{Status: "success", Data: struct {
				ResultType string            "json:\"resultType,omitempty\""
				Result     []logmodel.PodLog "json:\"result,omitempty\""
			}{ResultType: "vector", Result: []logmodel.PodLog{{Time: "", Namespace: "", Pod: "", Container: "", Log: ""}}}}},
		{"09_query_with_count_over_time_function_and_binary_operator", "/api/v1/query",
			Params{"query": `count_over_time(pod{namespace="namespace02"}[2m])>0`, "logFormat": "json"},
			TempDTO{Status: "success", Data: struct {
				ResultType string            "json:\"resultType,omitempty\""
				Result     []logmodel.PodLog "json:\"result,omitempty\""
			}{ResultType: "vector", Result: []logmodel.PodLog{{Time: "", Namespace: "", Pod: "", Container: "", Log: ""}}}}},
		{"10_query_namespace03_with_count_over_time_function_and_binary_operator", "/api/v1/query",
			Params{"query": `count_over_time(pod{namespace="namespace03"}[2m])>0`, "logFormat": "json"},
			TempDTO{Status: "success", Data: struct {
				ResultType string            "json:\"resultType,omitempty\""
				Result     []logmodel.PodLog "json:\"result,omitempty\""
			}{ResultType: "vector", Result: []logmodel.PodLog{}}},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d_%s", i, tc.name), func(t *testing.T) {
			got := GET(tc.url, tc.param)
			var gotDTO TempDTO
			err := json.Unmarshal([]byte(got), &gotDTO)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, gotDTO)
		})
	}
}
