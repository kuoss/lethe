package queryservice

import (
	"context"
	"fmt"
	"testing"

	"github.com/kuoss/common/tester"
	"github.com/kuoss/lethe/clock"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/letheql"
	"github.com/kuoss/lethe/letheql/model"
	"github.com/kuoss/lethe/storage/fileservice"
	"github.com/kuoss/lethe/storage/logservice"
	"github.com/kuoss/lethe/storage/logservice/logmodel"
	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	testCases := []struct {
		qs        string
		tr        model.TimeRange
		wantError string
		want      *letheql.Result
	}{
		// letheQL
		{
			`foo`,
			model.TimeRange{},
			"unknown logType: foo",
			&letheql.Result{},
		},
		{
			`pod`,
			model.TimeRange{},
			"getTargets err: target matcher err: not found label 'namespace' for logType 'pod'",
			&letheql.Result{},
		},
		{
			`pod{namespace="namespace01"}`,
			model.TimeRange{},
			"",
			&letheql.Result{Err: error(nil), Value: model.Log{Name: "pod", Lines: []model.LogLine{
				logmodel.PodLog{Time: "2009-11-10T22:59:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "lorem ipsum"},
				logmodel.PodLog{Time: "2009-11-10T22:59:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"}}}, Warnings: model.Warnings(nil)},
		},
		// logQL
		{
			`{job="pod"}`,
			model.TimeRange{},
			"getTargets err: target matcher err: not found label 'namespace' for logType 'pod'",
			&letheql.Result{},
		},
		{
			`{job="pod",namespace="namespace01"}`,
			model.TimeRange{},
			"",
			&letheql.Result{Err: error(nil), Value: model.Log{Name: "pod", Lines: []model.LogLine{
				logmodel.PodLog{Time: "2009-11-10T22:59:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "lorem ipsum"},
				logmodel.PodLog{Time: "2009-11-10T22:59:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"}}}, Warnings: model.Warnings(nil)},
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i, tc.qs), func(t *testing.T) {
			clock.SetPlaygroundMode(true)
			defer clock.SetPlaygroundMode(false)

			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			cfg, err := config.New("test")
			require.NoError(t, err)
			fileService, err := fileservice.New(cfg)
			require.NoError(t, err)
			logService := logservice.New(cfg, fileService)
			queryService := New(cfg, logService)
			require.NotEmpty(t, queryService)

			res := queryService.Query(context.TODO(), tc.qs, tc.tr)
			if tc.wantError == "" {
				require.NoError(t, res.Err)
			} else {
				require.EqualError(t, res.Err, tc.wantError)
			}
			res.Err = nil
			require.Equal(t, tc.want, res)
		})
	}
}

// TestToLetheQL tests the toLetheQL function against various scenarios.
func TestToLetheQL(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Standard cases
		{
			input: `{job="pod"}`,
			want:  `pod`,
		},
		{
			input: `{job="pod",namespace="namespace01"}`,
			want:  `pod{namespace="namespace01"}`,
		},
		{
			input: `{namespace="namespace01",job="pod"}`,
			want:  `pod{namespace="namespace01"}`,
		},
		{
			input: `{job="node"}`,
			want:  `node`,
		},
		{
			input: `{job="node",node="node01"}`,
			want:  `node{node="node01"}`,
		},
		{
			input: `{node="node01",job="node"}`,
			want:  `node{node="node01"}`,
		},
		// Cases with additional string " |= "hello world""
		{
			input: `{job="pod"} |= "hello world"`,
			want:  `pod |= "hello world"`,
		},
		{
			input: `{job="pod",namespace="namespace01"} |= "hello world"`,
			want:  `pod{namespace="namespace01"} |= "hello world"`,
		},
		{
			input: `{namespace="namespace01",job="pod"} |= "hello world"`,
			want:  `pod{namespace="namespace01"} |= "hello world"`,
		},
		{
			input: `{job="node"} |= "hello world"`,
			want:  `node |= "hello world"`,
		},
		{
			input: `{job="node",node="node01"} |= "hello world"`,
			want:  `node{node="node01"} |= "hello world"`,
		},
		{
			input: `{node="node01",job="node"} |= "hello world"`,
			want:  `node{node="node01"} |= "hello world"`,
		},
		// Edge cases and error handling
		{
			input: `just a string`, // Does not start with '{'
			want:  `just a string`,
		},
		{
			input: `some{job="test"}string`, // Does not start with '{'
			want:  `some{job="test"}string`,
		},
		{
			input: ``, // Empty input
			want:  ``,
		},
		{
			input: `{}`, // Empty braces, no job
			want:  `{}`,
		},
		{
			input: `{key="value"}`, // No job attribute
			want:  `{key="value"}`,
		},
		{
			input: `{job="test"}`, // Basic job
			want:  `test`,
		},
		{
			input: `{job="test",attr1="val1",attr2="val2"}`, // job at beginning
			want:  `test{attr1="val1",attr2="val2"}`,
		},
		{
			input: `{attr1="val1",job="test",attr2="val2"}`, // job in middle
			want:  `test{attr1="val1",attr2="val2"}`,
		},
		{
			input: `{attr1="val1",attr2="val2",job="test"}`, // job at end
			want:  `test{attr1="val1",attr2="val2"}`,
		},
		{
			input: `{job="test", }`, // job with trailing comma
			want:  `test`,
		},
		{
			input: `{ ,job="test"}`, // job with leading comma
			want:  `test`,
		},
		{
			input: `{ ,,job="test",,attr="value",, }`, // Many commas and spaces
			want:  `test{attr="value"}`,
		},
		{
			input: `{job="pod", namespace="namespace01" }`, // Spaces around attributes
			want:  `pod{namespace="namespace01"}`,
		},
		{
			input: `{job="pod" , namespace="namespace01"}`, // Spaces around commas
			want:  `pod{namespace="namespace01"}`,
		},
		{
			input: `{job="test"`, // Malformed: missing closing brace
			want:  `{job="test"`,
		},
		{
			input: `{attr="value",job="test`, // Malformed: missing closing quote for job
			want:  `{attr="value",job="test`,
		},
	}

	for i, tc := range tests {
		t.Run(tester.CaseName(i, tc.input), func(t *testing.T) {
			actual := toLetheQL(tc.input)
			require.Equal(t, tc.want, actual, fmt.Sprintf("Input: %q", tc.input))
		})
	}
}
