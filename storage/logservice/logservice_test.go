package logservice

import (
	// "fmt"
	"testing"
	// "time"

	// "github.com/kuoss/lethe/clock"
	// "github.com/kuoss/lethe/letheql/model"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.NotEmpty(t, logService1)
}

// func TestGetLog(t *testing.T) {
// 	testCases := []struct {
// 		LogSelector model.LogSelector
// 		want        model.Logs
// 		wantError   string
// 	}{
// 		// basic test
// 		{
// 			model.LogSelector{},
// 			[]model.LogLine{}, "not supported log type: ''",
// 		},
// 		{
// 			model.LogSelector{LogType: model.LogTypeNode},
// 			[]model.LogLine{}, "targetFilter is zero",
// 		},
// 		{
// 			model.LogSelector{LogType: model.LogTypeNode, NodeLogFilter: model.NodeLogFilter{}},
// 			[]model.LogLine{}, "targetFilter is zero",
// 		},
// 		// ======== LabelFilter ========
// 		{
// 			model.LogSelector{LogType: model.LogTypeNode, NodeLogFilter: model.NodeLogFilter{NodeFilter: model.LabelFilter{}}},
// 			[]model.LogLine{}, "targetFilter is zero",
// 		},
// 		// LabelFilterOperatorEqual
// 		{
// 			model.LogSelector{LogType: model.LogTypeNode, NodeLogFilter: model.NodeLogFilter{NodeFilter: model.LabelFilter{Operator: model.LabelFilterOperatorEqual}}},
// 			[]model.LogLine{},
// 			"",
// 		},
// 		{
// 			model.LogSelector{LogType: model.LogTypeNode, NodeLogFilter: model.NodeLogFilter{NodeFilter: model.LabelFilter{Operator: model.LabelFilterOperatorEqual, Operand: ""}}},
// 			[]model.LogLine{},
// 			"",
// 		},
// 		{
// 			model.LogSelector{LogType: model.LogTypeNode, NodeLogFilter: model.NodeLogFilter{NodeFilter: model.LabelFilter{Operator: model.LabelFilterOperatorEqual, Operand: "hello"}}},
// 			[]model.LogLine{},
// 			"",
// 		},
// 		{
// 			model.LogSelector{LogType: model.LogTypeNode, NodeLogFilter: model.NodeLogFilter{NodeFilter: model.LabelFilter{Operator: model.LabelFilterOperatorEqual, Operand: "node"}}},
// 			[]model.LogLine{},
// 			"",
// 		},
// 		{
// 			model.LogSelector{LogType: model.LogTypeNode, NodeLogFilter: model.NodeLogFilter{NodeFilter: model.LabelFilter{Operator: model.LabelFilterOperatorEqual, Operand: "node01"}}},
// 			[]model.LogLine{
// 				model.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "lerom ipsum"},
// 				model.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "lerom from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
// 				model.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"}},
// 			"",
// 		},
// 		// LabelFilterOperatorNotEqual
// 		{
// 			model.LogSelector{LogType: model.LogTypeNode, NodeLogFilter: model.NodeLogFilter{NodeFilter: model.LabelFilter{Operator: model.LabelFilterOperatorNotEqual, Operand: "node02"}}},
// 			[]model.LogLine{
// 				model.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "lerom ipsum"},
// 				model.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "lerom from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
// 				model.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"}},
// 			"",
// 		},
// 		// LabelFilterOperatorEqualRegex
// 		{
// 			model.LogSelector{LogType: model.LogTypeNode, NodeLogFilter: model.NodeLogFilter{NodeFilter: model.LabelFilter{Operator: model.LabelFilterOperatorEqualRegex, Operand: "node"}}},
// 			[]model.LogLine{
// 				model.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "lerom ipsum"},
// 				model.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "lerom from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
// 				model.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
// 				model.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "containerd", Log: "hello world"},
// 				model.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "containerd", Log: "lerom ipsum"},
// 				model.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "containerd", Log: "hello world"},
// 				model.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "dockerd", Log: "hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "dockerd", Log: "lerom from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "dockerd", Log: "hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T21:58:00.000000Z", Node: "node02", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"}},
// 			"",
// 		},
// 		{
// 			model.LogSelector{LogType: model.LogTypeNode, NodeLogFilter: model.NodeLogFilter{NodeFilter: model.LabelFilter{Operator: model.LabelFilterOperatorEqualRegex, Operand: "node01"}}},
// 			[]model.LogLine{
// 				model.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "lerom ipsum"},
// 				model.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "lerom from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
// 				model.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"}},
// 			"",
// 		},
// 		// LabelFilterOperatorNotRegex
// 		{
// 			model.LogSelector{LogType: model.LogTypeNode, NodeLogFilter: model.NodeLogFilter{NodeFilter: model.LabelFilter{Operator: model.LabelFilterOperatorNotRegex, Operand: "node02"}}},
// 			[]model.LogLine{
// 				model.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:56:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "lerom ipsum"},
// 				model.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:57:00.000000Z", Node: "node01", Process: "kubelet", Log: "I0525 20:00:45.752587   17221 scope.go:110] \"RemoveContainer\" hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "lerom from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
// 				model.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"}},
// 			"",
// 		},
// 		{
// 			model.LogSelector{LogType: model.LogTypeNode, NodeLogFilter: model.NodeLogFilter{NodeFilter: model.LabelFilter{Operator: model.LabelFilterOperatorNotRegex, Operand: "node"}}},
// 			[]model.LogLine{},
// 			"",
// 		},
// 		// NodeFilter + ProcessFilter
// 		{
// 			model.LogSelector{LogType: model.LogTypeNode, NodeLogFilter: model.NodeLogFilter{
// 				NodeFilter:    model.LabelFilter{Operator: model.LabelFilterOperatorNotRegex, Operand: "node02"},
// 				ProcessFilter: model.LabelFilter{Operator: model.LabelFilterOperatorNotRegex, Operand: "kubelet"},
// 			}},
// 			[]model.LogLine{
// 				model.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "lerom ipsum"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "lerom from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:59:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"},
// 				model.NodeLog{Time: "2009-11-10T23:00:00.000000Z", Node: "node01", Process: "containerd", Log: "hello world"}},
// 			"",
// 		},
// 		// NodeFilter + ProcessFilter + LineFilters
// 		{
// 			model.LogSelector{LogType: model.LogTypeNode,
// 				NodeLogFilter: model.NodeLogFilter{
// 					NodeFilter:    model.LabelFilter{Operator: model.LabelFilterOperatorNotRegex, Operand: "node02"},
// 					ProcessFilter: model.LabelFilter{Operator: model.LabelFilterOperatorNotRegex, Operand: "kubelet"}},
// 				LineFilters: []model.LineFilter{
// 					{Operator: model.LineFilterOperatorPipeEqual, Operand: "hello"},
// 					{Operator: model.LineFilterOperatorPipeEqual, Operand: "sidecar"}}},
// 			[]model.LogLine{
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "lerom from sidecar"},
// 				model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"}},
// 			"",
// 		},
// 	}
// 	for i, tc := range testCases {
// 		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
// 			got, err := logService1.GetLogs(tc.LogSelector)
// 			if tc.wantError == "" {
// 				assert.NoError(t, err)
// 			} else {
// 				assert.EqualError(t, err, tc.wantError)
// 			}
// 			assert.Equal(t, tc.want, got)
// 		})
// 	}
// }

// func TestAppendLogLine(t *testing.T) {
// 	now := clock.Now()
// 	timeRange := model.TimeRange{Start: now.Add(-240 * time.Hour), End: now}
// 	testCases := []struct {
// 		line        string
// 		LogSelector model.LogSelector
// 		want        []model.LogLine
// 	}{
// 		{
// 			"",
// 			model.LogSelector{},
// 			[]model.LogLine{}, // no time separator
// 		},
// 		{
// 			"hello",
// 			model.LogSelector{},
// 			[]model.LogLine{}, // no time separator
// 		},
// 		{
// 			"2009-11-10T22:58:00.000000Z[node01|dockerd] hello from sidecar",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange},
// 			[]model.LogLine{}, // not in time range
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypeNode, TimeRange: timeRange},
// 			[]model.LogLine{},
// 		},
// 		{
// 			"2009-11-10T22:58:00.000000Z[node01|dockerd] hello from sidecar",
// 			model.LogSelector{LogType: model.LogTypeNode, TimeRange: timeRange},
// 			[]model.LogLine{model.NodeLog{Time: "2009-11-10T22:58:00.000000Z", Node: "node01", Process: "dockerd", Log: "hello from sidecar"}},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange},
// 			[]model.LogLine{model.PodLog{Time: "2009-11-10T23:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world [ddd12wewe]"}},
// 		},
// 		// ========= PodLogFilter ========
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, PodLogFilter: model.PodLogFilter{}},
// 			[]model.LogLine{model.PodLog{Time: "2009-11-10T23:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world [ddd12wewe]"}},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, PodLogFilter: model.PodLogFilter{NamespaceFilter: model.LabelFilter{}}},
// 			[]model.LogLine{model.PodLog{Time: "2009-11-10T23:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world [ddd12wewe]"}},
// 		},
// 		// An target(namespace, node) filter does nothing. Because the appendmodel.LogLine() function doesn't handle that.
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, PodLogFilter: model.PodLogFilter{NamespaceFilter: model.LabelFilter{Operator: model.LabelFilterOperatorEqualRegex, Operand: ".*01"}}},
// 			[]model.LogLine{model.PodLog{Time: "2009-11-10T23:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world [ddd12wewe]"}},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, PodLogFilter: model.PodLogFilter{NamespaceFilter: model.LabelFilter{Operator: model.LabelFilterOperatorEqualRegex, Operand: ".*02"}}},
// 			[]model.LogLine{model.PodLog{Time: "2009-11-10T23:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world [ddd12wewe]"}},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, PodLogFilter: model.PodLogFilter{PodFilter: model.LabelFilter{Operator: model.LabelFilterOperatorEqualRegex, Operand: "nginx"}}},
// 			[]model.LogLine{model.PodLog{Time: "2009-11-10T23:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world [ddd12wewe]"}},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, PodLogFilter: model.PodLogFilter{PodFilter: model.LabelFilter{Operator: model.LabelFilterOperatorEqualRegex, Operand: "nginx.*"}}},
// 			[]model.LogLine{model.PodLog{Time: "2009-11-10T23:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world [ddd12wewe]"}},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, PodLogFilter: model.PodLogFilter{PodFilter: model.LabelFilter{Operator: model.LabelFilterOperatorEqualRegex, Operand: ".*nginx"}}},
// 			[]model.LogLine{model.PodLog{Time: "2009-11-10T23:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world [ddd12wewe]"}},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, PodLogFilter: model.PodLogFilter{PodFilter: model.LabelFilter{Operator: model.LabelFilterOperatorEqualRegex, Operand: ".*-nginx"}}},
// 			[]model.LogLine{},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange},
// 			[]model.LogLine{model.PodLog{Time: "2009-11-10T23:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world [ddd12wewe]"}},
// 		},
// 		// line filter
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, LineFilters: []model.LineFilter{}},
// 			[]model.LogLine{model.PodLog{Time: "2009-11-10T23:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world [ddd12wewe]"}},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, LineFilters: []model.LineFilter{{Operator: model.LineFilterOperatorPipeEqual, Operand: "hello"}}},
// 			[]model.LogLine{model.PodLog{Time: "2009-11-10T23:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world [ddd12wewe]"}},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, LineFilters: []model.LineFilter{{Operator: model.LineFilterOperatorNotEqual, Operand: "hello"}}},
// 			[]model.LogLine{},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, LineFilters: []model.LineFilter{{Operator: model.LineFilterOperatorPipeRegex, Operand: "hello"}}},
// 			[]model.LogLine{model.PodLog{Time: "2009-11-10T23:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world [ddd12wewe]"}},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, LineFilters: []model.LineFilter{{Operator: model.LineFilterOperatorPipeRegex, Operand: "hello.*"}}},
// 			[]model.LogLine{model.PodLog{Time: "2009-11-10T23:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world [ddd12wewe]"}},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, LineFilters: []model.LineFilter{{Operator: model.LineFilterOperatorPipeRegex, Operand: ".*hello"}}},
// 			[]model.LogLine{model.PodLog{Time: "2009-11-10T23:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world [ddd12wewe]"}},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, LineFilters: []model.LineFilter{{Operator: model.LineFilterOperatorNotRegex, Operand: "hello"}}},
// 			[]model.LogLine{},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, LineFilters: []model.LineFilter{{Operator: model.LineFilterOperatorNotEqual, Operand: "foo"}}},
// 			[]model.LogLine{model.PodLog{Time: "2009-11-10T23:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world [ddd12wewe]"}},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, LineFilters: []model.LineFilter{
// 				{Operator: model.LineFilterOperatorPipeEqual, Operand: "hello"},
// 				{Operator: model.LineFilterOperatorPipeEqual, Operand: "world"},
// 			}},
// 			[]model.LogLine{model.PodLog{Time: "2009-11-10T23:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world [ddd12wewe]"}},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, LineFilters: []model.LineFilter{
// 				{Operator: model.LineFilterOperatorPipeRegex, Operand: "hello"},
// 				{Operator: model.LineFilterOperatorPipeRegex, Operand: "world"},
// 			}},
// 			[]model.LogLine{model.PodLog{Time: "2009-11-10T23:00:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world [ddd12wewe]"}},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, LineFilters: []model.LineFilter{
// 				{Operator: model.LineFilterOperatorNotRegex, Operand: "hello"},
// 				{Operator: model.LineFilterOperatorPipeRegex, Operand: "world"},
// 			}},
// 			[]model.LogLine{},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, LineFilters: []model.LineFilter{
// 				{Operator: model.LineFilterOperatorPipeRegex, Operand: "hello"},
// 				{Operator: model.LineFilterOperatorNotRegex, Operand: "world"},
// 			}},
// 			[]model.LogLine{},
// 		},
// 		{
// 			"2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]",
// 			model.LogSelector{LogType: model.LogTypePod, TimeRange: timeRange, LineFilters: []model.LineFilter{
// 				{Operator: model.LineFilterOperatorNotRegex, Operand: "hello"},
// 				{Operator: model.LineFilterOperatorNotRegex, Operand: "world"},
// 			}},
// 			[]model.LogLine{},
// 		},
// 	}
// 	for i, tc := range testCases {
// 		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
// 			ls := tc.LogSelector
// 			funcSet := ls.GetFilterFuncSet()
// 			// fmt.Println(ls.PodLogFilter.PodFilter)

// 			got := []model.LogLine{}
// 			appendLogLine(&got, tc.line, &tc.LogSelector, &funcSet)
// 			assert.Equal(t, tc.want, got)
// 		})
// 	}
// }
