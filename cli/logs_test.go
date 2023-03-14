package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	// cliutil "github.com/kuoss/lethe/cli/util"
	testutil "github.com/kuoss/lethe/testutil"
)

func init() {
	testutil.Init()
	// testutil.SetTestLogFiles()
	// rotator = logs.NewRotator()
	// config.SetWriter(cliutil.GetWriter())
}

func Test_logs(t *testing.T) {

	var got string
	var want string

	got = execute("logs")
	want = "error: logs command needs an flag: --query"
	assert.Equal(t, want, got)

	got = execute("logs", "--query")
	want = "Error: flag needs an argument: --query\nUsage:\n  lethetool logs [flags]\n\nFlags:\n  -h, --help           help for logs\n  -q, --query string   letheql"
	assert.Equal(t, want, got)

	got = execute("logs", "--query", `pod{namespace=""}`)
	want = "namespace value cannot be empty"
	assert.Equal(t, want, got)

	got = execute("logs", "--query", `pod{namespace="ns-not-exists"}`)
	want = ""
	assert.Equal(t, want, got)

	got = execute("logs", "--query", `pod{namespace="namespace01"}`)
	want = "2009-11-10T21:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\n2009-11-10T21:01:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\n2009-11-10T21:02:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\n2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\n2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\n2009-11-10T22:56:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\n2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\n2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\n2009-11-10T22:57:00.000000Z[namespace01|apache-75675f5897-7ci7o|httpd] hello from sidecar\n2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\n2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar\n2009-11-10T22:58:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\n2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum\n2009-11-10T22:59:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world"
	assert.Equal(t, want, got)

	got = execute("logs", "--query", `pod{namespace="namespace02"}`)
	want = "2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world\n2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] lerom ipsum\n2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|nginx] hello world\n2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\n2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] lerom from sidecar\n2009-11-10T22:58:00.000000Z[namespace02|nginx-deployment-75675f5897-7ci7o|sidecar] hello from sidecar\n2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\n2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\n2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\n2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\n2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar\n2009-11-10T22:58:00.000000Z[namespace02|apache-75675f5897-7ci7o|httpd] hello from sidecar"
	assert.Equal(t, want, got)

	time.Sleep(3000 * time.Millisecond)
}
