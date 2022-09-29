package main

import (
	"testing"

	"github.com/kuoss/lethe/testutil"
)

func Test_logs(t *testing.T) {
	testutil.SetTestLogs()
	var got string
	var want string

	got = execute("logs")
	want = `"error: logs command needs an flag: --query"`
	testutil.CheckEqualJSON(t, got, want)

	got = execute("logs", "--query")
	want = `"Error: flag needs an argument: --query\nUsage:\n  lethetool logs [flags]\n\nFlags:\n  -h, --help           help for logs\n  -q, --query string   letheql"`
	testutil.CheckEqualJSON(t, got, want)

	got = execute("logs", "--query", `pod{namespace=""}`)
	want = `"namespace value cannot be empty"`
	testutil.CheckEqualJSON(t, got, want)

	got = execute("logs", "--query", `pod{namespace="ns-not-exists"}`)
	// want = `"not found matched namespace directories: err=exit status 2"`
	want = `"{logs [] 0}"`
	testutil.CheckEqualJSON(t, got, want)

	got = execute("logs", "--query", `pod{namespace="namespace01"}`)
	want = `"{logs [2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar 2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar 2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar 2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum 2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar 2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar 2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar 2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar 2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar 2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar 2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world 2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world 2009-11-10T23:01:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world 2009-11-10T23:02:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world] 0}"`
	testutil.CheckEqualJSON(t, got, want)

	got = execute("logs", "--query", `pod{namespace="namespace02"}`)
	want = `"{logs [2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world 2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum 2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world 2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar 2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar 2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar 2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar 2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar 2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar 2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar 2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar 2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar] 0}"`
	testutil.CheckEqualJSON(t, got, want)
}
