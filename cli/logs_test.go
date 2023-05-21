package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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

	tests := map[string]struct {
		args []string
		want string
	}{
		"logs --query": {
			args: []string{"logs", "--query"},
			want: "Error: flag needs an argument: --query\nUsage:\n  lethetool logs [flags]\n\nFlags:\n  -h, --help           help for logs\n  -q, --query string   letheql",
		},
		`logs --query pod{namespace=""}`: {
			args: []string{"logs", "--query", `pod{namespace=""}`},
			want: "namespace value cannot be empty",
		},
		`logs --query pod{namespace="ns-not-exists"}`: {
			args: []string{"logs", "--query", `pod{namespace="ns-not-exists"}`},
			want: "",
		},
		`logs --query pod{namespace="namespace01"}`: {
			args: []string{"logs", "--query", `pod{namespace="namespace01"}`},
			want: "{ 2009-11-10T21:00:00Z namespace01 nginx-deployment-75675f5897-7ci7o nginx  hello world}\n{ 2009-11-10T21:01:00Z namespace01 nginx-deployment-75675f5897-7ci7o nginx  hello world}\n{ 2009-11-10T21:02:00Z namespace01 nginx-deployment-75675f5897-7ci7o nginx  hello world}\n{ 2009-11-10T22:56:00Z namespace01 apache-75675f5897-7ci7o httpd  hello from sidecar}\n{ 2009-11-10T22:56:00Z namespace01 apache-75675f5897-7ci7o httpd  hello from sidecar}\n{ 2009-11-10T22:56:00Z namespace01 apache-75675f5897-7ci7o httpd  hello from sidecar}\n{ 2009-11-10T22:57:00Z namespace01 apache-75675f5897-7ci7o httpd  hello from sidecar}\n{ 2009-11-10T22:57:00Z namespace01 apache-75675f5897-7ci7o httpd  hello from sidecar}\n{ 2009-11-10T22:57:00Z namespace01 apache-75675f5897-7ci7o httpd  hello from sidecar}\n{ 2009-11-10T22:58:00Z namespace01 nginx-deployment-75675f5897-7ci7o sidecar  hello from sidecar}\n{ 2009-11-10T22:58:00Z namespace01 nginx-deployment-75675f5897-7ci7o sidecar  lerom from sidecar}\n{ 2009-11-10T22:58:00Z namespace01 nginx-deployment-75675f5897-7ci7o sidecar  hello from sidecar}\n{ 2009-11-10T22:59:00Z namespace01 nginx-deployment-75675f5897-7ci7o nginx  lerom ipsum}\n{ 2009-11-10T22:59:00Z namespace01 nginx-deployment-75675f5897-7ci7o nginx  hello world}",
		},
		`logs --query pod{namespace="namespace02"}`: {
			args: []string{"logs", "--query", `pod{namespace="namespace02"}`},
			want: "{ 2009-11-10T22:58:00Z namespace02 nginx-deployment-75675f5897-7ci7o nginx  hello world}\n{ 2009-11-10T22:58:00Z namespace02 nginx-deployment-75675f5897-7ci7o nginx  lerom ipsum}\n{ 2009-11-10T22:58:00Z namespace02 nginx-deployment-75675f5897-7ci7o nginx  hello world}\n{ 2009-11-10T22:58:00Z namespace02 nginx-deployment-75675f5897-7ci7o sidecar  hello from sidecar}\n{ 2009-11-10T22:58:00Z namespace02 nginx-deployment-75675f5897-7ci7o sidecar  lerom from sidecar}\n{ 2009-11-10T22:58:00Z namespace02 nginx-deployment-75675f5897-7ci7o sidecar  hello from sidecar}\n{ 2009-11-10T22:58:00Z namespace02 apache-75675f5897-7ci7o httpd  hello from sidecar}\n{ 2009-11-10T22:58:00Z namespace02 apache-75675f5897-7ci7o httpd  hello from sidecar}\n{ 2009-11-10T22:58:00Z namespace02 apache-75675f5897-7ci7o httpd  hello from sidecar}\n{ 2009-11-10T22:58:00Z namespace02 apache-75675f5897-7ci7o httpd  hello from sidecar}\n{ 2009-11-10T22:58:00Z namespace02 apache-75675f5897-7ci7o httpd  hello from sidecar}\n{ 2009-11-10T22:58:00Z namespace02 apache-75675f5897-7ci7o httpd  hello from sidecar}",
		},
	}

	for name, test := range tests {
		t.Run(name, func(subt *testing.T) {
			got := execute(test.args...)
			assert.Equal(subt, test.want, got)
		})
	}
	time.Sleep(3000 * time.Millisecond)
}
