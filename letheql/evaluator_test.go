package letheql

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kuoss/lethe/letheql/model"
	"github.com/kuoss/lethe/letheql/parser"
	"github.com/kuoss/lethe/storage/logservice"
	"github.com/kuoss/lethe/storage/logservice/logmodel"
	"github.com/stretchr/testify/assert"
)

func TestNewEvaluator(t *testing.T) {
	logService1 := &logservice.LogService{}
	contextBackground := context.Background()
	contextTODO := context.TODO()
	testCases := []struct {
		logService     *logservice.LogService
		ctx            context.Context
		startTimestamp int64
		endTimestamp   int64
		interval       int64
		want           evaluator
	}{
		{
			logService1, contextBackground, 0, 0, 0,
			evaluator{logService1, contextBackground, 0, 0, 0, time.Time{}, time.Time{}},
		},
		{
			logService1, contextTODO, 0, 0, 0,
			evaluator{logService1, contextTODO, 0, 0, 0, time.Time{}, time.Time{}},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got := evaluator{tc.logService, tc.ctx, tc.startTimestamp, tc.endTimestamp, tc.interval, time.Time{}, time.Time{}}
			assert.Equal(t, tc.want, got)
		})
	}
	assert.NotNil(t, evaluator1)
}

func TestEval(t *testing.T) {
	testCases := []struct {
		input          string
		wantParseError string
		wantError      string
		wantWarnings   model.Warnings
		want           parser.Value
	}{
		{
			``,
			"1:1: parse error: no expression found in input", "", nil, nil,
		},
		{
			`pod`,
			"", "getTargets err: target matcher err: not found label 'namespace' for logType 'pod'", nil, nil,
		},
		{
			`"hello"`,
			"", "", nil, String{T: 0, V: "hello"},
		},
		{
			`pod{namespace="hello"}`,
			"", "", nil, model.Log{Name: "pod", Lines: []model.LogLine{}},
		},
		{
			`pod{namespace|="namespace01"}`,
			"1:14: parse error: unexpected character inside braces: '|'", "", nil, nil,
		},
		{
			`pod{namespace="namespace01"}`,
			"", "", nil,
			model.Log{Name: "pod", Lines: []model.LogLine{
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
			`pod{namespace!="namespace01"}`,
			"", "",
			model.Warnings{fmt.Errorf("warnMultiTargets: use operator '=' for selecting target")},
			model.Log{Name: "pod", Lines: []model.LogLine{
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
			`pod{namespace=~"namespace.*"}`,
			"", "",
			model.Warnings{fmt.Errorf("warnMultiTargets: use operator '=' for selecting target")},
			model.Log{Name: "pod", Lines: []model.LogLine{
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
			`pod{namespace="namespace01",pod="nginx-.*"}`,
			"", "", nil,
			model.Log{Name: "pod", Lines: []model.LogLine{}},
		},
		{
			`pod{namespace="namespace01",pod~="nginx-.*"}`,
			"1:32: parse error: unexpected character inside braces: '~'", "", nil, nil,
		},
		{
			`pod{namespace="namespace01",pod=~"nginx-.*"}`,
			"", "", nil,
			model.Log{Name: "pod", Lines: []model.LogLine{
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
			`pod{namespace="namespace01",pod!~"nginx-.*"}`,
			"", "", nil,
			model.Log{Name: "pod", Lines: []model.LogLine{
				logmodel.PodLog{Time: "2009-11-10T22:56:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:56:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:56:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:57:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:57:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:57:00.000000Z", Namespace: "namespace01", Pod: "apache-75675f5897-7ci7o", Container: "httpd", Log: "hello from sidecar"}}},
		},
		{
			`pod{namespace="namespace01",pod=~"nginx-.*",container="sidecar"}`,
			"", "", nil,
			model.Log{Name: "pod", Lines: []model.LogLine{
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "lerom from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "hello from sidecar"}}},
		},
		{
			`node`,
			"", "getTargets err: target matcher err: not found label 'node' for logType 'node'", nil, nil,
		},
		{
			`node{namespace="namespace01"}`,
			"", "getTargets err: target matcher err: not found label 'node' for logType 'node'", nil, nil,
		},
		{
			`node{node="node01"}`,
			"", "", nil,
			model.Log{Name: "node", Lines: []model.LogLine{
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
			`node{node=~"node.*"}`,
			"", "",
			model.Warnings{fmt.Errorf("warnMultiTargets: use operator '=' for selecting target")},
			model.Log{Name: "node", Lines: []model.LogLine{
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
			`node{node=~"node.*"}`,
			"", "",
			model.Warnings{fmt.Errorf("warnMultiTargets: use operator '=' for selecting target")},
			model.Log{Name: "node", Lines: []model.LogLine{
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
			`node{node!~"node.*"}`,
			"", "",
			model.Warnings{fmt.Errorf("warnMultiTargets: use operator '=' for selecting target")},
			model.Log{Name: "node", Lines: []model.LogLine{}},
		},
		{
			`node{node="node01",process!="kubelet"}`,
			"", "", nil,
			model.Log{Name: "node", Lines: []model.LogLine{
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "lerom ipsum"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "lerom from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"}}},
		},
		{
			`node{node="node01",process!="kubelet"} |= "hello"`,
			"", "", nil,
			model.Log{Name: "node", Lines: []model.LogLine{
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"}}},
		},
		{
			`node{node="node01",process!="kubelet"} |= "hello" |= "sidecar"`,
			"", "", nil,
			model.Log{Name: "node", Lines: []model.LogLine{
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
				logmodel.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"}}},
		},
		{
			`node{node="node01",process!="kubelet"} |= "hello" != "sidecar"`,
			"", "", nil,
			model.Log{Name: "node", Lines: []model.LogLine{
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"}}},
		},
		{
			`node{node="node01",process!="kubelet"} |~ "ll.*" !~ "car.*"`,
			"", "", nil,
			model.Log{Name: "node", Lines: []model.LogLine{
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"}}},
		},
		{
			`node{node="node01",process!="kubelet"} |~ "ll.*" !~ "car.*"`,
			"", "", nil,
			model.Log{Name: "node", Lines: []model.LogLine{
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"}}},
		},
		{
			`node{node="node01",process!="kubelet"} != "sidecar" |~ "d$" `,
			"", "", nil,
			model.Log{Name: "node", Lines: []model.LogLine{
				logmodel.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
				logmodel.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"}}},
		},
		{
			`node{node="node01",process!="kubelet"} != "sidecar" |~ "*" `,
			"", "getMatchFuncSet err: getLineMatchFuncs err: getLineMatchFunc err: error parsing regexp: missing argument to repetition operator: `*`", nil, nil,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			var got parser.Value
			var ws model.Warnings
			var err error
			expr, parseErr := parser.ParseExpr(tc.input)
			if tc.wantParseError != "" {
				assert.EqualError(t, parseErr, tc.wantParseError)
			} else {
				assert.NoError(t, parseErr)
				got, ws, err = evaluator1.Eval(expr)
			}
			if tc.wantError != "" {
				assert.EqualError(t, err, tc.wantError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantWarnings, ws)
		})
	}
}
