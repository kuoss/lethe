package letheql

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus/prometheus/storage"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.NotNil(t, engine1)
}

func TestNewInstantQuery_smoke(t *testing.T) {
	queryable := storage.QueryableFunc(func(ctx context.Context, mint, maxt int64) (storage.Querier, error) {
		return nil, nil
	})
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	qry, err := engine1.NewInstantQuery(ctx, queryable, "pod", time.Unix(1, 0))
	assert.NoError(t, err)
	assert.NotEmpty(t, qry)
}

// func TestNewInstantQuery(t *testing.T) {
// 	qry, err := engine1.NewInstantQuery(context.TODO(), "pod", clock.Now())
// 	assert.NoError(t, err)
// 	result := qry.Exec(context.TODO())
// 	assert.Equal(t, nil, result)
// }

// func TestProcQuery(t *testing.T) {
// 	testCases := []struct {
// 		query     string
// 		want      QueryData
// 		wantError string
// 	}{
// 		{
// 			"",
// 			QueryData{},
// 			"parseExpr err: 1:1: parse error: no expression found in input",
// 		},
// 		{
// 			`pod{namespace="namespace01"}`,
// 			QueryData{ResultType: "", Logs: []log.LogLine(nil), Scalar: 0},
// 			"",
// 		},
// 	}
// 	for i, tc := range testCases {
// 		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
// 			got, err := service.ProcQuery(tc.query, TimeRange{})
// 			if tc.wantError == "" {
// 				assert.NoError(t, err)
// 			} else {
// 				assert.EqualError(t, err, tc.wantError)
// 			}
// 			assert.Equal(t, tc.want, got)
// 		})
// 	}
// 	t.Fatal("")
// }

// func ago(m int) time.Time {
// 	return clock.Now().Add(time.Duration(-m) * time.Minute)
// }

// func Test_Query_Success(t *testing.T) {
// 	testCases := []struct {
// 		query string
// 		want  QueryData
// 	}{
// 		{`pod{namespace="namespace01"}`,
// 			QueryData{ResultType: "logs", Logs: []logService.LogLine{logService.PodLog{Name: "", Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
// 		{`pod{namespace="not-exists"}`,
// 			QueryData{ResultType: ValueTypeLogs, Logs: []logService.LogLine{}}},
// 		{`pod{namespace="namespace01",pod="nginx"}`,
// 			QueryData{ResultType: "logs", Logs: []logService.LogLine{}, Scalar: 0}},
// 		{`pod{namespace="namespace01",pod="nginx-*"}`,
// 			QueryData{ResultType: "logs", Logs: []logService.LogLine{logService.PodLog{Name: "", Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
// 		{`pod{namespace="namespace01",container="nginx"}`,
// 			QueryData{ResultType: "logs", Logs: []logService.LogLine{logService.PodLog{Name: "", Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
// 		{`pod{namespace="namespace*",container="nginx"}`,
// 			QueryData{ResultType: "logs", Logs: []logService.LogLine{logService.PodLog{Name: "", Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
// 		{`pod{namespace="namespace01",pod="nginx-*"}[3m]`,
// 			QueryData{ResultType: "logs", Logs: []logService.LogLine{logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
// 		{`pod`,
// 			QueryData{ResultType: "logs", Logs: []logService.LogLine{logService.PodLog{Name: "", Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
// 		{`pod{}`,
// 			QueryData{ResultType: "logs", Logs: []logService.LogLine{logService.PodLog{Name: "", Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
// 		{`count_over_time(pod{namespace="namespace01",pod="nginx-*"}[3m])`,
// 			QueryData{ResultType: "scalar", Logs: nil, Scalar: 5}},
// 		{`1`,
// 			QueryData{ResultType: ValueTypeScalar, Logs: nil, Scalar: 1}},
// 		{`count_over_time(pod{}[3m])`,
// 			QueryData{ResultType: "scalar", Logs: nil, Scalar: 20}},
// 		{`count_over_time(pod{}[3m]) > 10`,
// 			QueryData{ResultType: ValueTypeScalar, Logs: nil, Scalar: 1}},
// 		{`count_over_time(pod{}[3m]) < 10`,
// 			QueryData{ResultType: ValueTypeScalar, Logs: nil, Scalar: 0}},
// 		{`count_over_time(pod{}[3m]) == 21`,
// 			QueryData{ResultType: ValueTypeScalar, Logs: nil, Scalar: 0}},
// 		{`pod{namespace="namespace01"} |= hello`,
// 			QueryData{ResultType: "logs", Logs: []logService.LogLine{logService.PodLog{Name: "", Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
// 		{`pod{namespace="namespace01"} != hello`,
// 			QueryData{ResultType: "logs", Logs: []logService.LogLine{logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}}}},
// 		{`pod{namespace="namespace01"} |~ (.ro.*o)`,
// 			QueryData{ResultType: ValueTypeLogs, Logs: []logService.LogLine{logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}}}},
// 		{`pod{namespace="namespace01"} !~ (.le)`,
// 			QueryData{ResultType: "logs", Logs: []logService.LogLine{logService.PodLog{Name: "", Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:56:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:57:00Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
// 	}

// 	for i, tc := range testCases {
// 		t.Run(fmt.Sprintf("#%d_%s", i, tc.query), func(t *testing.T) {
// 			got, err := ProcQuery(tc.query, TimeRange{})
// 			assert.NoError(t, err)
// 			assert.Equal(t, tc.want, got)
// 		})
// 	}
// }
// func Test_QueryWithTimeRange(t *testing.T) {

// 	now := clock.Now()

// 	testCases := []struct {
// 		query     string
// 		timeRange TimeRange
// 		want      QueryData
// 	}{
// 		// modify test case.. for supoorting sort by time
// 		{`pod{namespace="namespace01",pod="nginx-*"}[3m]`,
// 			TimeRange{},
// 			QueryData{ResultType: "logs", Logs: []logService.LogLine{logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
// 		{`pod{namespace="namespace01",pod="nginx-*"}`,
// 			TimeRange{Start: ago(999999), End: now},
// 			QueryData{ResultType: "logs", Logs: []logService.LogLine{logService.PodLog{Name: "", Time: "2009-11-10T21:00:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T21:01:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T21:02:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " lerom from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:58:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: " hello from sidecar"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
// 		{`pod{namespace="namespace01",pod="nginx-*"}`,
// 			TimeRange{Start: ago(1), End: now},
// 			QueryData{ResultType: "logs", Logs: []logService.LogLine{logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " lerom ipsum"}, logService.PodLog{Name: "", Time: "2009-11-10T22:59:00Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: " hello world"}}}},
// 		{`count_over_time(pod{namespace="*"})`,
// 			TimeRange{Start: ago(1), End: now},
// 			QueryData{ResultType: "scalar", Logs: nil, Scalar: 2}},
// 		{`count_over_time(pod{namespace="*"})`,
// 			TimeRange{Start: ago(2), End: now},
// 			QueryData{ResultType: "scalar", Logs: nil, Scalar: 17}},
// 		{`count_over_time(pod{namespace="*"})`,
// 			TimeRange{Start: ago(3), End: now},
// 			QueryData{ResultType: ValueTypeScalar, Logs: nil, Scalar: 20}},
// 	}

// 	for i, tc := range testCases {
// 		t.Run(fmt.Sprintf("#%d_%s", i, tc.query), func(t *testing.T) {
// 			got, err := ProcQuery(tc.query, tc.timeRange)
// 			assert.NoError(t, err)
// 			assert.Equal(t, tc.want, got)
// 		})
// 	}
// }

// func Test_QueryFail(t *testing.T) {
// 	now := clock.Now()
// 	testCases := []struct {
// 		query     string
// 		timeRange TimeRange
// 		wantError string
// 	}{
// 		// modify test case.. for supoorting sort by time
// 		{`pod{namespace=""}`,
// 			TimeRange{},
// 			"namespace value cannot be empty"},
// 		{`pod{foo=""}`,
// 			TimeRange{},
// 			"unknown label foo"},
// 		{`{namespace="hello"}`,
// 			TimeRange{},
// 			"a log name must be specified"},
// 		{`count_over_time(pod{namespace="*"})`,
// 			TimeRange{Start: ago(0), End: now},
// 			"end time and start time are the same"},
// 	}

// 	for i, tc := range testCases {
// 		t.Run(fmt.Sprintf("#%d_%s", i, tc.query), func(t *testing.T) {
// 			_, err := ProcQuery(tc.query, tc.timeRange)
// 			assert.EqualError(t, err, tc.wantError)
// 		})
// 	}
// }

// func TestNewQuery(t *testing.T) {
// 	e := &Engine{}
// 	testCases := []struct {
// 		name  string
// 		query string
// 		want  *query
// 	}{
// 		{`pod metric with namespace label matcher`,
// 			`pod{namespace="namespace01"}`,
// 			&query{q: `pod{namespace="namespace01"}`, filter: nil, keyword: "", engine: e}},
// 	}
// 	for i, tc := range testCases {
// 		t.Run(fmt.Sprintf("#%d_%s", i, tc.query), func(t *testing.T) {
// 			query, err := e.newQuery(tc.query)
// 			assert.NoError(t, err)
// 			assert.Equal(t, query, tc.want)
// 		})
// 	}
// }
