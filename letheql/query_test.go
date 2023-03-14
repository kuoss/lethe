package letheql

import (
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
			want: QueryData{ResultType: "logs", Logs: []string{"2009-11-10T21:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T21:01:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T21:02:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"}, Scalar: 0}},
		"02": {query: `pod{namespace="not-exists"}`,
			want: QueryData{ResultType: ValueTypeLogs, Logs: []string{}}},
		"03": {query: `pod{namespace="namespace01",pod="nginx"}`,
			want: QueryData{ResultType: "logs", Logs: []string{}, Scalar: 0}},
		"04": {query: `pod{namespace="namespace01",pod="nginx-*"}`,
			want: QueryData{ResultType: "logs", Logs: []string{"2009-11-10T21:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T21:01:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T21:02:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"}, Scalar: 0}},
		"05": {query: `pod{namespace="namespace01",container="nginx"}`,
			want: QueryData{ResultType: "logs", Logs: []string{"2009-11-10T21:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T21:01:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T21:02:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"}, Scalar: 0}},
		"06": {query: `pod{namespace="namespace*",container="nginx"}`,
			want: QueryData{ResultType: "logs", Logs: []string{"2009-11-10T21:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T21:01:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T21:02:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"}, Scalar: 0}},
		"07": {query: `pod{namespace="namespace01",pod="nginx-*"}[3m]`,
			want: QueryData{ResultType: "logs", Logs: []string{"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"}, Scalar: 0}},
		"08": {query: `pod`,
			want: QueryData{ResultType: "logs", Logs: []string{"2009-11-10T21:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T21:01:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T21:02:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"}, Scalar: 0}},
		"09": {query: `pod{}`,
			want: QueryData{ResultType: "logs", Logs: []string{"2009-11-10T21:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T21:01:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T21:02:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"}, Scalar: 0}},
		"10": {query: `count_over_time(pod{namespace="namespace01",pod="nginx-*"}[3m])`,
			want: QueryData{ResultType: "scalar", Logs: []string(nil), Scalar: 5}},
		"11": {query: `1`,
			want: QueryData{ResultType: ValueTypeScalar, Logs: nil, Scalar: 1}},
		"12": {query: `count_over_time(pod{}[3m])`,
			want: QueryData{ResultType: "scalar", Logs: []string(nil), Scalar: 20}},
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
			want: QueryData{ResultType: "logs", Logs: []string{"2009-11-10T21:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T21:01:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T21:02:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"}, Scalar: 0}},
		"17": {query: `pod{namespace="namespace01"} != hello`,
			want: QueryData{ResultType: "logs", Logs: []string{"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum"}, Scalar: 0}},
		"18": {query: `pod{namespace="namespace01"} |~ (.ro.*o)`,
			want: QueryData{ResultType: ValueTypeLogs, Logs: []string{"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar"}, Scalar: 0}},
		"19": {query: `pod{namespace="namespace01"} !~ (.le)`,
			want: QueryData{ResultType: "logs", Logs: []string{"2009-11-10T21:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T21:01:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T21:02:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"}, Scalar: 0}},
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
		"01": {query: `pod{namespace="namespace01",pod="nginx-*"}[3m]`, timeRange: TimeRange{},
			want: QueryData{ResultType: "logs", Logs: []string{"2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"}, Scalar: 0}},
		"02": {query: `pod{namespace="namespace01",pod="nginx-*"}`, timeRange: TimeRange{Start: ago(999999), End: now},
			want: QueryData{ResultType: "logs", Logs: []string{"2009-11-10T21:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T21:01:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T21:02:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar", "2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"}, Scalar: 0}},
		"03": {query: `pod{namespace="namespace01",pod="nginx-*"}`, timeRange: TimeRange{Start: ago(1), End: now},
			want: QueryData{ResultType: "logs", Logs: []string{"2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum", "2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"}, Scalar: 0}},
		"04": {query: `count_over_time(pod{namespace="*"})`, timeRange: TimeRange{Start: ago(1), End: now},
			want: QueryData{ResultType: "scalar", Logs: []string(nil), Scalar: 2}},
		"05": {query: `count_over_time(pod{namespace="*"})`, timeRange: TimeRange{Start: ago(2), End: now},
			want: QueryData{ResultType: "scalar", Logs: []string(nil), Scalar: 17}},
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
