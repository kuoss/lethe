package letheql

import (
	"fmt"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPromqlParse(t *testing.T) {

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
	}
	for name, tt := range tests {
		t.Run(name+" "+tt.query, func(subt *testing.T) {
			expr, err := parser.ParseExpr(tt.query)
			if err != nil {
				return
			}
			fmt.Printf("%#v", expr)
			got, err := ProcQuery(tt.query, TimeRange{})
			if err != nil {
				subt.Fatalf("query: %s error: %s", name, tt.query)
			}
			assert.Equal(subt, tt.want, got)
		})
	}
}
