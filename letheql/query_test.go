package letheql

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/kuoss/lethe/config"
	_ "github.com/kuoss/lethe/storage/driver/filesystem"
	"github.com/kuoss/lethe/testutil"
)

func Test_Query_Success(t *testing.T) {

	testutil.SetTestLogs()

	tests := map[string]struct {
		query string
		want  QueryData
	}{
		// modify test case.. for supoorting sort by time
		"with namespace01": {query: `pod{namespace="namespace01"}`, want: QueryData{
			ResultType: ValueTypeLogs,
			Logs:       []string{"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"},
			Scalar:     0,
		}},
		"with not exist namespace": {query: `pod{namespace="not-exists"}`, want: QueryData{
			ResultType: ValueTypeLogs,
			Logs:       []string{},
		}},
		"with namespace01 and pod=nginx": {query: `pod{namespace="namespace01",pod="nginx"}`, want: QueryData{
			ResultType: ValueTypeLogs,
			Logs:       []string{},
		}},
		"with namespace01 and pod=nginx-*": {query: `pod{namespace="namespace01",pod="nginx-*"}`, want: QueryData{
			ResultType: ValueTypeLogs,
			Logs:       []string{"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"},
			Scalar:     0,
		}},
		"with namespace01 and container=nginx": {query: `pod{namespace="namespace01",container="nginx"}`, want: QueryData{
			ResultType: ValueTypeLogs,
			Logs:       []string{"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"},
			Scalar:     0,
		}},
		"with namespace* and container=nginx": {query: `pod{namespace="namespace*",container="nginx"}`, want: QueryData{
			ResultType: ValueTypeLogs,
			Logs:       []string{"2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"},
			Scalar:     0,
		}},
		"with duration": {query: `pod{namespace="namespace01",pod="nginx-*"}[3m]`, want: QueryData{
			ResultType: ValueTypeLogs,
			Logs:       []string{"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"},
			Scalar:     0,
		}},
		"all pod log": {query: `pod`, want: QueryData{
			ResultType: ValueTypeLogs,
			Logs:       []string{"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"},
			Scalar:     0,
		}},
		"all pod log pod{}": {query: `pod{}`, want: QueryData{
			ResultType: ValueTypeLogs,
			Logs:       []string{"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"},
			Scalar:     0,
		}},
		"count over time function": {query: `count_over_time(pod{namespace="namespace01",pod="nginx-*"}[3m])`, want: QueryData{
			ResultType: ValueTypeScalar,
			Logs:       nil,
			Scalar:     6,
		}},
		"scalar": {query: `1`, want: QueryData{
			ResultType: ValueTypeScalar,
			Logs:       nil,
			Scalar:     1,
		}},
		"count over time with duration": {query: `count_over_time(pod{}[3m])`, want: QueryData{
			ResultType: ValueTypeScalar,
			Logs:       nil,
			Scalar:     21,
		}},
		"operator > count over time with duration": {query: `count_over_time(pod{}[3m]) > 10`, want: QueryData{
			ResultType: ValueTypeScalar,
			Logs:       nil,
			Scalar:     1,
		}},
		"operator < count over time with duration": {query: `count_over_time(pod{}[3m]) < 10`, want: QueryData{
			ResultType: ValueTypeScalar,
			Logs:       nil,
			Scalar:     0,
		}},
		"operator = count over time with duration": {query: `count_over_time(pod{}[3m]) == 21`, want: QueryData{
			ResultType: ValueTypeScalar,
			Logs:       nil,
			Scalar:     1,
		}},

		// todo
		// scalar operator same with filter ( != )
		"operator != count over time with duration": {query: `count_over_time(pod{}[3m]) != 21`, want: QueryData{
			ResultType: ValueTypeScalar,
			Logs:       nil,
			Scalar:     0,
		}},

		"with namespace01 and include hello keyword": {query: `pod{namespace="namespace01"} |= hello`, want: QueryData{
			ResultType: ValueTypeLogs,
			Logs:       []string{"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"},
			Scalar:     0,
		}},
		"with namespace01 and exclude hello keyword": {query: `pod{namespace="namespace01"} != hello`, want: QueryData{
			ResultType: ValueTypeLogs,
			Logs:       []string{"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum"},
			Scalar:     0,
		}},
		"with namespace01 and includeRegex  keyword": {query: `pod{namespace="namespace01"} |~ (.ro.*o)`, want: QueryData{
			ResultType: ValueTypeLogs,
			Logs:       []string{"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar"},
			Scalar:     0,
		}},
		"with namespace01 and excludeRegex  keyword": {query: `pod{namespace="namespace01"} !~ (.le)`, want: QueryData{
			ResultType: ValueTypeLogs,
			Logs:       []string{"2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"},
			Scalar:     0,
		}},
	}

	for name, tt := range tests {
		t.Run(name, func(subt *testing.T) {
			got, err := ProcQuery(tt.query, TimeRange{})
			if err != nil {
				subt.Fatalf("query: %s error: %s", name, tt.query)
			}
			assert.Equal(subt, tt.want, got)
		})
	}
}

func ago(m int) time.Time {
	return config.GetNow().Add(time.Duration(-m) * time.Minute)
}
func Test_QueryWithTimeRange(t *testing.T) {

	testutil.SetTestLogs()

	now := config.GetNow()

	tests := map[string]struct {
		query     string
		timeRange TimeRange
		want      QueryData
	}{
		// modify test case.. for supoorting sort by time
		"without time range": {query: `pod{namespace="namespace01",pod="nginx-*"}[3m]`, timeRange: TimeRange{}, want: QueryData{
			ResultType: ValueTypeLogs,
			Logs:       []string{"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"},
			Scalar:     0,
		}},
		"with from 999999min ago until now": {query: `pod{namespace="namespace01",pod="nginx-*"}`, timeRange: TimeRange{Start: ago(999999), End: now}, want: QueryData{
			ResultType: ValueTypeLogs,
			Logs:       []string{"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"},
			Scalar:     0,
		}},
		"with from 1min ago until now": {query: `pod{namespace="namespace01",pod="nginx-*"}`, timeRange: TimeRange{Start: ago(1), End: now}, want: QueryData{
			ResultType: ValueTypeLogs,
			Logs:       []string{"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"},
			Scalar:     0,
		}},
		"count_over_time function from 1min until now": {query: `count_over_time(pod{namespace="*"})`, timeRange: TimeRange{Start: ago(1), End: now}, want: QueryData{
			ResultType: ValueTypeScalar,
			Logs:       nil,
			Scalar:     3,
		}},
		"count_over_time function from 2min until now": {query: `count_over_time(pod{namespace="*"})`, timeRange: TimeRange{Start: ago(2), End: now}, want: QueryData{
			ResultType: ValueTypeScalar,
			Logs:       nil,
			Scalar:     18,
		}},
		"count_over_time function from 3min until now": {query: `count_over_time(pod{namespace="*"})`, timeRange: TimeRange{Start: ago(3), End: time.Time{}}, want: QueryData{
			ResultType: ValueTypeScalar,
			Logs:       nil,
			Scalar:     21,
		}},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s query: %s", name, tt.query), func(subt *testing.T) {
			got, err := ProcQuery(tt.query, tt.timeRange)
			if err != nil {
				subt.Fatalf("query: %s error: %s", name, tt.query)
			}
			assert.Equal(subt, tt.want, got)
		})
	}
}

func Test_QueryFail(t *testing.T) {

	tests := map[string]struct {
		query     string
		timeRange TimeRange
		want      string
	}{
		// modify test case.. for supoorting sort by time
		"explicitly empty namespace":    {query: `pod{namespace=""}`, timeRange: TimeRange{}, want: "namespace value cannot be empty"},
		"unknown label":                 {query: `pod{foo=""}`, timeRange: TimeRange{}, want: "unknown label foo"},
		"explicitly empty log type":     {query: `{namespace="hello"}`, timeRange: TimeRange{}, want: "a log name must be specified"},
		"start time same with end time": {query: `count_over_time(pod{namespace="*"})`, timeRange: TimeRange{Start: ago(0), End: now}, want: "end time and start time are the same"},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s query: %s", name, tt.query), func(subt *testing.T) {
			_, err := ProcQuery(tt.query, tt.timeRange)
			if assert.Error(subt, err) {
				assert.Equal(subt, tt.want, err.Error())
			}
		})
	}
}
