package letheql

import (
	"fmt"
	"testing"
	"time"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/testutil"
)

func Test_Query_Success(t *testing.T) {
	testutil.SetTestLogs()

	var query string
	var want string
	var got QueryData
	var err error

	query = `pod{namespace="namespace01"}`
	want = `{"resultType":"logs","logs":["2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]}`
	got, err = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s ", err, len(got.Logs), query))

	query = `pod{namespace="not-exists"}`
	want = `{"resultType":"logs"}`
	got, err = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s ", err, len(got.Logs), query))

	query = `pod{namespace="namespace01",pod="nginx"}`
	want = `{"resultType":"logs"}`
	got, err = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s ", err, len(got.Logs), query))

	// wildcard on pod
	query = `pod{namespace="namespace01",pod="nginx-*"}`
	want = `{"resultType":"logs","logs":["2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]}`
	got, err = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s ", err, len(got.Logs), query))

	query = `pod{namespace="namespace01",container="nginx"}`
	want = `{"resultType":"logs","logs":["2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]}`
	got, err = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s ", err, len(got.Logs), query))

	// wildcard on namespace
	query = `pod{namespace="namespace*",container="nginx"}`
	want = `{"resultType":"logs","logs":["2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]}`
	got, err = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s ", err, len(got.Logs), query))

	// duration
	query = `pod{namespace="namespace01",pod="nginx-*"}[3m]`
	want = `{"resultType":"logs","logs":["2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]}`
	got, err = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s ", err, len(got.Logs), query))

	// all pod logs
	query = `pod`
	want = `{"resultType":"logs","logs":["2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]}`
	got, err = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s ", err, len(got.Logs), query))

	// all pod logs
	query = `pod{}`
	want = `{"resultType":"logs","logs":["2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]}`
	got, err = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s ", err, len(got.Logs), query))

	// function
	query = `count_over_time(pod{namespace="namespace01",pod="nginx-*"}[3m])`
	want = `{"resultType":"scalar","scalar":6}`
	got, err = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s ", err, len(got.Logs), query))

	// scalar
	query = `1`
	want = `{"resultType":"scalar","scalar":1}`
	got, err = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s ", err, len(got.Logs), query))

	// // operator
	query = `count_over_time(pod{}[3m])`
	want = `{"resultType":"scalar","scalar":21}`
	got, err = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s ", err, len(got.Logs), query))

	query = `count_over_time(pod{}[3m]) > 10`
	want = `{"resultType":"scalar","scalar":1}`
	got, err = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s ", err, len(got.Logs), query))

	query = `count_over_time(pod{}[3m]) < 10`
	want = `{"resultType":"scalar"}`
	got, err = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s ", err, len(got.Logs), query))

	query = `count_over_time(pod{}[3m])`
	want = `{"resultType":"scalar","scalar":21}`
	got, err = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s ", err, len(got.Logs), query))

	query = `count_over_time(pod{}[3m]) == 21`
	want = `{"resultType":"scalar","scalar":1}`
	got, err = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s ", err, len(got.Logs), query))

	query = `count_over_time(pod{}[3m]) != 21`
	want = `{"resultType":"scalar"}`
	got, err = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s ", err, len(got.Logs), query))

	// query = `count_over_time(pod{namespace="namespace01",pod="nginx-*"}[3m]) > 0`
	// want = `{"resultType":"number","number":6}`
	// got, err = ProcQuery(query, TimeRange{})
	// testutil.CheckEqualJSON(t, got, want,fmt.Sprintf("err=%v lines=%d query=%s ", err, len(got.Logs), query))
}

func ago(m int) time.Time {
	return config.GetNow().Add(time.Duration(-m) * time.Minute)
}
func Test_QueryWithTimeRange(t *testing.T) {

	testutil.SetTestLogs()

	var query string
	var want string
	var got QueryData
	var err error
	var timeRange TimeRange

	now := config.GetNow()

	query = `pod{namespace="namespace01",pod="nginx-*"}[3m]`
	timeRange = TimeRange{}
	want = `{"resultType":"logs","logs":["2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]}`
	got, err = ProcQuery(query, timeRange)
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s timeRange=%v", err, len(got.Logs), query, timeRange))

	query = `pod{namespace="namespace01",pod="nginx-*"}`
	timeRange = TimeRange{Start: ago(999999), End: now}
	want = `{"resultType":"logs","logs":["2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]}`
	got, err = ProcQuery(query, timeRange)
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s timeRange=%v", err, len(got.Logs), query, timeRange))

	query = `pod{namespace="namespace01",pod="nginx-*"}`
	timeRange = TimeRange{Start: ago(1), End: now}
	want = `{"resultType":"logs","logs":["2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]}`
	got, err = ProcQuery(query, timeRange)
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s timeRange=%v", err, len(got.Logs), query, timeRange))

	query = `count_over_time(pod{namespace="*"})`
	timeRange = TimeRange{Start: ago(1), End: now}
	want = `{"resultType":"scalar","scalar":3}`
	got, err = ProcQuery(query, timeRange)
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s timeRange=%v", err, len(got.Logs), query, timeRange))

	query = `count_over_time(pod{namespace="*"})`
	timeRange = TimeRange{Start: ago(2), End: now}
	want = `{"resultType":"scalar","scalar":18}`
	got, err = ProcQuery(query, timeRange)
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s timeRange=%v", err, len(got.Logs), query, timeRange))

	query = `count_over_time(pod{namespace="*"})`
	timeRange = TimeRange{Start: ago(2), End: now}
	want = `{"resultType":"scalar","scalar":18}`
	got, err = ProcQuery(query, timeRange)
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s timeRange=%v", err, len(got.Logs), query, timeRange))

	query = `count_over_time(pod{namespace="*"})`
	want = `{"resultType":"scalar","scalar":21}`
	got, err = ProcQuery(query, TimeRange{Start: ago(3), End: time.Time{}})
	testutil.CheckEqualJSON(t, got, want, fmt.Sprintf("err=%v lines=%d query=%s timeRange=%v", err, len(got.Logs), query, timeRange))
}

func Test_QueryFail(t *testing.T) {
	var query string
	var want string
	var got error
	var timeRange TimeRange

	// explicitly empty
	query = `pod{namespace=""}`
	want = `"namespace value cannot be empty"`
	_, got = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, "query=", query, " err=", got)

	// unknown label
	query = `pod{foo=""}`
	want = `"unknown label foo"`
	_, got = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, "query=", query, " err=", got)

	query = `{namespace="hello"}`
	want = `"a log name must be specified"`
	_, got = ProcQuery(query, TimeRange{})
	testutil.CheckEqualJSON(t, got, want, "query=", query, " err=", got)

	query = `count_over_time(pod{namespace="*"})`
	timeRange = TimeRange{Start: ago(0), End: now}
	want = `"end time and start time are the same"`
	_, got = ProcQuery(query, timeRange)
	testutil.CheckEqualJSON(t, got, want, "query=", query, " err=", got)

}
