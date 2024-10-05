package router

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/kuoss/common/tester"
	"github.com/kuoss/lethe/clock"
	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	testCases := []struct {
		qs       string
		wantCode int
		wantBody string
	}{
		{
			`hello`,
			500, `{"error":"unknown logType: hello","errorType":"queryError","status":"error"}`,
		},
		{
			`pod`,
			500, `{"error":"getTargets err: target matcher err: not found label 'namespace' for logType 'pod'","errorType":"queryError","status":"error"}`,
		},
		{
			`pod{}`,
			500, `{"error":"getTargets err: target matcher err: not found label 'namespace' for logType 'pod'","errorType":"queryError","status":"error"}`,
		},
		{
			`pod{namespace="namespace01"}`,
			200, `{"data":{"result":[{"time":"2009-11-10T22:59:00.000000Z","namespace":"namespace01","pod":"nginx-deployment-75675f5897-7ci7o","container":"nginx","log":"lerom ipsum"},{"time":"2009-11-10T22:59:00.000000Z","namespace":"namespace01","pod":"nginx-deployment-75675f5897-7ci7o","container":"nginx","log":"hello world"}],"resultType":"logs"},"status":"success"}`,
		},
		{
			`pod{namespace="namespace01"} |= "ipsum"`,
			200, `{"data":{"result":[{"time":"2009-11-10T22:59:00.000000Z","namespace":"namespace01","pod":"nginx-deployment-75675f5897-7ci7o","container":"nginx","log":"lerom ipsum"}],"resultType":"logs"},"status":"success"}`,
		},
		{
			`node{node="node01",process!="kubelet"} |= "hello" != "sidecar"`,
			200, `{"data":{"result":[{"time":"2009-11-10T23:00:00.000000Z","node":"node01","process":"containerd","log":"hello world"}],"resultType":"logs"},"status":"success"}`,
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i, tc.qs), func(t *testing.T) {
			clock.SetPlaygroundMode(true)
			defer clock.SetPlaygroundMode(false)

			v := url.Values{}
			v.Add("query", tc.qs)
			code, body, cleanup := testGET(t, "/api/v1/query?"+v.Encode())
			defer cleanup()
			require.Equal(t, tc.wantCode, code)
			require.Equal(t, tc.wantBody, body)
		})
	}
}

func TestQueryRange(t *testing.T) {
	clock.SetPlaygroundMode(true)
	defer clock.SetPlaygroundMode(false)

	now := clock.Now()
	ago10d := now.Add(-240 * time.Hour)
	testCases := []struct {
		qs       string
		start    time.Time
		end      time.Time
		wantCode int
		wantBody string
	}{
		{
			`hello`, ago10d, now,
			500, `{"error":"unknown logType: hello","errorType":"queryError","status":"error"}`,
		},
		{
			`pod`, ago10d, now,
			500, `{"error":"getTargets err: target matcher err: not found label 'namespace' for logType 'pod'","errorType":"queryError","status":"error"}`,
		},
		{
			`pod{}`, ago10d, now,
			500, `{"error":"getTargets err: target matcher err: not found label 'namespace' for logType 'pod'","errorType":"queryError","status":"error"}`,
		},
		{
			`pod{namespace="namespace01"}`, ago10d, now,
			200, `{"data":{"result":[{"time":"2009-11-10T21:00:00.000000Z","namespace":"namespace01","pod":"nginx-deployment-75675f5897-7ci7o","container":"nginx","log":"hello world"},{"time":"2009-11-10T21:01:00.000000Z","namespace":"namespace01","pod":"nginx-deployment-75675f5897-7ci7o","container":"nginx","log":"hello world"},{"time":"2009-11-10T21:02:00.000000Z","namespace":"namespace01","pod":"nginx-deployment-75675f5897-7ci7o","container":"nginx","log":"hello world"},{"time":"2009-11-10T22:56:00.000000Z","namespace":"namespace01","pod":"apache-75675f5897-7ci7o","container":"httpd","log":"hello from sidecar"},{"time":"2009-11-10T22:56:00.000000Z","namespace":"namespace01","pod":"apache-75675f5897-7ci7o","container":"httpd","log":"hello from sidecar"},{"time":"2009-11-10T22:56:00.000000Z","namespace":"namespace01","pod":"apache-75675f5897-7ci7o","container":"httpd","log":"hello from sidecar"},{"time":"2009-11-10T22:59:00.000000Z","namespace":"namespace01","pod":"nginx-deployment-75675f5897-7ci7o","container":"nginx","log":"lerom ipsum"},{"time":"2009-11-10T22:57:00.000000Z","namespace":"namespace01","pod":"apache-75675f5897-7ci7o","container":"httpd","log":"hello from sidecar"},{"time":"2009-11-10T22:57:00.000000Z","namespace":"namespace01","pod":"apache-75675f5897-7ci7o","container":"httpd","log":"hello from sidecar"},{"time":"2009-11-10T22:57:00.000000Z","namespace":"namespace01","pod":"apache-75675f5897-7ci7o","container":"httpd","log":"hello from sidecar"},{"time":"2009-11-10T22:58:00.000000Z","namespace":"namespace01","pod":"nginx-deployment-75675f5897-7ci7o","container":"sidecar","log":"hello from sidecar"},{"time":"2009-11-10T22:58:00.000000Z","namespace":"namespace01","pod":"nginx-deployment-75675f5897-7ci7o","container":"sidecar","log":"lerom from sidecar"},{"time":"2009-11-10T22:58:00.000000Z","namespace":"namespace01","pod":"nginx-deployment-75675f5897-7ci7o","container":"sidecar","log":"hello from sidecar"},{"time":"2009-11-10T22:59:00.000000Z","namespace":"namespace01","pod":"nginx-deployment-75675f5897-7ci7o","container":"nginx","log":"hello world"}],"resultType":"logs"},"status":"success"}`,
		},
		{
			`pod{namespace="namespace01"} |= "ipsum"`, ago10d, now,
			200, `{"data":{"result":[{"time":"2009-11-10T22:59:00.000000Z","namespace":"namespace01","pod":"nginx-deployment-75675f5897-7ci7o","container":"nginx","log":"lerom ipsum"}],"resultType":"logs"},"status":"success"}`,
		},
		{
			`node{node="node01",process!="kubelet"} |= "hello" != "sidecar"`, ago10d, now,
			200, `{"data":{"result":[{"time":"2009-11-10T22:59:00.000000Z","node":"node01","process":"containerd","log":"hello world"},{"time":"2009-11-10T23:00:00.000000Z","node":"node01","process":"containerd","log":"hello world"}],"resultType":"logs"},"status":"success"}`,
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i), func(t *testing.T) {

			v := url.Values{}
			v.Add("query", tc.qs)
			v.Add("start", time2string(tc.start))
			v.Add("end", time2string(tc.end))
			code, body, cleanup := testGET(t, "/api/v1/query_range?"+v.Encode())
			defer cleanup()
			require.Equal(t, tc.wantCode, code)
			require.Equal(t, tc.wantBody, body)
		})
	}
}

func time2string(t time.Time) string {
	return fmt.Sprintf("%d", t.Unix())
}

func Test_FloatStringToTime(t *testing.T) {
	want := time.Date(2015, time.July, 1, 20, 10, 51, 780999898, time.UTC)
	got := floatStringToTime("1435781451.781")
	require.True(t, want.Equal(got))
}
