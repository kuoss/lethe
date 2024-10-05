package letheql

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kuoss/common/tester"
	"github.com/kuoss/lethe/clock"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/letheql/model"
	"github.com/kuoss/lethe/letheql/parser"
	"github.com/kuoss/lethe/storage/fileservice"
	"github.com/kuoss/lethe/storage/logservice"
	"github.com/kuoss/lethe/storage/logservice/logmodel"
	"github.com/stretchr/testify/require"
)

func newLogService(t *testing.T) *logservice.LogService {
	cfg, err := config.New("test")
	require.NoError(t, err)

	fileService, err := fileservice.New(cfg)
	require.NoError(t, err)

	return logservice.New(cfg, fileService)
}

func TestNewEvaluator(t *testing.T) {
	logService := newLogService(t)
	testCases := []struct {
		startTimestamp int64
		endTimestamp   int64
		interval       int64
		want           evaluator
	}{
		{
			0, 0, 0,
			evaluator{
				logService,
				context.TODO(),
				0,
				0,
				0,
				time.Time{},
				time.Time{},
			},
		},
		{
			0, 0, 0,
			evaluator{logService, context.TODO(), 0, 0, 0, time.Time{}, time.Time{}},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got := evaluator{logService, context.TODO(), tc.startTimestamp, tc.endTimestamp, tc.interval, time.Time{}, time.Time{}}
			require.Equal(t, tc.want, got)
		})
	}
}

func TestEval(t *testing.T) {
	testCases := []struct {
		input          string
		wantParseError bool
		wantEvalError  bool
		wantWarnings   model.Warnings
		want           parser.Value
	}{
		// error
		{
			input:          ``,
			wantParseError: true,
		},
		{
			input:          `pod{namespace|="namespace01"}`,
			wantParseError: true,
		},
		{
			input:          `pod{namespace="namespace01",pod~="nginx-.*"}`,
			wantParseError: true,
		},
		{
			input:         `pod`,
			wantEvalError: true,
		},
		{
			input:         `node`,
			wantEvalError: true,
		},
		{
			input:         `node{namespace="namespace01"}`,
			wantEvalError: true,
		},
		{
			input:         `node{node="node01",process!="kubelet"} != "sidecar" |~ "*" `,
			wantEvalError: true,
		},
		{
			input: `"hello"`,
			want:  String{T: 0, V: "hello"},
		},
		{
			input: `pod{namespace="hello"}`,
			want:  model.Log{Name: "pod", Lines: []model.LogLine{}},
		},
		{
			input: `pod{namespace="namespace01",pod="nginx-.*"}`,
			want:  model.Log{Name: "pod", Lines: []model.LogLine{}},
		},
		{
			input: `pod{namespace="namespace01"}`,
			want: model.Log{Name: "pod", Lines: []model.LogLine{
				logmodel.PodLog{Time: "2009-11-10T21:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"},
				logmodel.PodLog{Time: "2009-11-10T21:01:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"},
				logmodel.PodLog{Time: "2009-11-10T21:02:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"},
				logmodel.PodLog{Time: "2009-11-10T22:56:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:56:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:56:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:59:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "lerom ipsum"},
				logmodel.PodLog{Time: "2009-11-10T22:57:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:57:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:57:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "lerom from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:59:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"}}},
		},
		{
			input: `pod{namespace="namespace01",pod=~"nginx-.*"}`,
			want: model.Log{Name: "pod", Lines: []model.LogLine{
				logmodel.PodLog{Time: "2009-11-10T21:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"},
				logmodel.PodLog{Time: "2009-11-10T21:01:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"},
				logmodel.PodLog{Time: "2009-11-10T21:02:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"},
				logmodel.PodLog{Time: "2009-11-10T22:59:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "lerom ipsum"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "lerom from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:59:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"}}},
		},
		{
			input: `pod{namespace="namespace01",pod!~"nginx-.*"}`,
			want: model.Log{Name: "pod", Lines: []model.LogLine{
				logmodel.PodLog{Time: "2009-11-10T22:56:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:56:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:56:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:57:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:57:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:57:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"}}},
		},
		{
			input: `pod{namespace="namespace01",pod=~"nginx-.*",container="sidecar"}`,
			want: model.Log{Name: "pod", Lines: []model.LogLine{
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "lerom from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "hello from sidecar"}}},
		},
		{
			input: `node{node="node01"}`,
			want: model.Log{Name: "node", Lines: []model.LogLine{
				logmodel.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "lerom ipsum"},
				logmodel.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "lerom from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"}}},
		},
		{
			input: `node{node="node01",process!="kubelet"}`,
			want: model.Log{Name: "node", Lines: []model.LogLine{
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "lerom ipsum"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "lerom from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"}}},
		},
		{
			input: `node{node="node01",process!="kubelet"} |= "hello"`,
			want: model.Log{Name: "node", Lines: []model.LogLine{
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"}}},
		},
		{
			input: `node{node="node01",process!="kubelet"} |= "hello" |= "sidecar"`,
			want: model.Log{Name: "node", Lines: []model.LogLine{
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"}}},
		},
		{
			input: `node{node="node01",process!="kubelet"} |= "hello" != "sidecar"`,
			want: model.Log{Name: "node", Lines: []model.LogLine{
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"}}},
		},
		{
			input: `node{node="node01",process!="kubelet"} |~ "ll.*" !~ "car.*"`,
			want: model.Log{Name: "node", Lines: []model.LogLine{
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"}}},
		},
		{
			input: `node{node="node01",process!="kubelet"} |~ "ll.*" !~ "car.*"`,
			want: model.Log{Name: "node", Lines: []model.LogLine{
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"}}},
		},
		{
			input: `node{node="node01",process!="kubelet"} != "sidecar" |~ "d$" `,
			want: model.Log{Name: "node", Lines: []model.LogLine{
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"}}},
		},
		// warnings
		{
			input:        `node{node!~"node.*"}`,
			wantWarnings: model.Warnings{fmt.Errorf("warnMultiTargets: use operator '=' for selecting target")},
			want:         model.Log{Name: "node", Lines: []model.LogLine{}},
		},
		{
			input:        `pod{namespace!="namespace01"}`,
			wantWarnings: model.Warnings{fmt.Errorf("warnMultiTargets: use operator '=' for selecting target")},
			want: model.Log{Name: "pod", Lines: []model.LogLine{
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "lerom ipsum"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "lerom from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"}}},
		},
		{
			input:        `pod{namespace=~"namespace.*"}`,
			wantWarnings: model.Warnings{fmt.Errorf("warnMultiTargets: use operator '=' for selecting target")},
			want: model.Log{Name: "pod", Lines: []model.LogLine{
				logmodel.PodLog{Time: "2009-11-10T21:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"},
				logmodel.PodLog{Time: "2009-11-10T21:01:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"},
				logmodel.PodLog{Time: "2009-11-10T21:02:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"},
				logmodel.PodLog{Time: "2009-11-10T22:56:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:56:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:56:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:59:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "lerom ipsum"},
				logmodel.PodLog{Time: "2009-11-10T22:57:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:57:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:57:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "lerom from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:59:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "lerom ipsum"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "lerom from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace02", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"}}},
		},
		{
			input:        `node{node=~"node.*"}`,
			wantWarnings: model.Warnings{fmt.Errorf("warnMultiTargets: use operator '=' for selecting target")},
			want: model.Log{Name: "node", Lines: []model.LogLine{
				logmodel.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "lerom ipsum"},
				logmodel.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "lerom from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "containerd", Log: "lerom ipsum"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "dockerd", Log: "lerom from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"}}},
		},
		{
			input:        `node{node=~"node.*"}`,
			wantWarnings: model.Warnings{fmt.Errorf("warnMultiTargets: use operator '=' for selecting target")},
			want: model.Log{Name: "node", Lines: []model.LogLine{
				logmodel.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "lerom ipsum"},
				logmodel.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "lerom from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "containerd", Log: "lerom ipsum"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "dockerd", Log: "lerom from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"}}},
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i, tc.input), func(t *testing.T) {
			clock.SetPlaygroundMode(true)
			defer clock.SetPlaygroundMode(false)

			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			logService := newLogService(t)
			var got parser.Value
			var ws model.Warnings
			var evalErr error
			expr, parseErr := parser.ParseExpr(tc.input)
			if tc.wantParseError {
				require.Error(t, parseErr)
			} else {
				require.NoError(t, parseErr)

				now := clock.Now()
				ev := evaluator{
					logService:     logService,
					ctx:            context.TODO(),
					startTimestamp: 0,
					endTimestamp:   0,
					interval:       0,
					start:          now.Add(-4 * time.Hour),
					end:            now,
				}
				got, ws, evalErr = ev.Eval(expr)
			}
			if tc.wantEvalError {
				require.Error(t, evalErr)
			} else {
				require.NoError(t, evalErr)
			}
			require.Equal(t, tc.want, got)
			require.Equal(t, tc.wantWarnings, ws)
		},
		)
	}
}
