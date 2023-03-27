package filter

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
"with namespace01 and include hello keyword": {query: `pod{namespace="namespace01"} |= hello`, want: `{"resultType":"logs","logs":["2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]}`},
		"with namespace01 and exclude hello keyword": {query: `pod{namespace="namespace01"} != hello`, want: `{"resultType":"logs","logs":["2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum"]}`},
		"with namespace01 and includeRegex  keyword": {query: `pod{namespace="namespace01"} |~ (.ro.*o)`, want: `{"resultType":"logs","logs":["2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar"]}`},
		"with namespace01 and excludeRegex  keyword": {query: `pod{namespace="namespace01"} !~ (.le)`, want: `{"resultType":"logs","logs":["2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar","2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world","2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"]}`},

*/

func TestGetFilterFromQuery(t *testing.T) {
	tests := map[string]struct {
		query  string
		filter Filter
	}{
		"include filter": {query: `pod{namespace="namespace01"} |= hello`, filter: &includeFilter{keyword: "hello"}},
		"exclude filter": {query: `pod{namespace="namespace01"} != hello`, filter: &excludeFilter{keyword: "hello"}},

		"include regex filter": {query: `pod{namespace="namespace01"} |~ (.ro.*o)`, filter: &includeRegexFilter{
			regex:   regexp.MustCompile(`(.ro.*o)`),
			keyword: `(.ro.*o)`},
		},
		"exclude regex filter": {query: `pod{namespace="namespace01"} !~ (.ro.*o)`, filter: &excludeRegexFilter{
			regex:   regexp.MustCompile(`(.ro.*o)`),
			keyword: `(.ro.*o)`},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(subt *testing.T) {
			got, err := FromQuery(tt.query)
			if err != nil {
				subt.Fatalf("query: %s err: %s", name, err.Error())
			}

			assert.Equal(subt, tt.filter, got)
		})
	}
}
