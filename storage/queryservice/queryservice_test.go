package queryservice

import (
	"context"
	"fmt"
	"testing"

	"github.com/kuoss/lethe/letheql"
	"github.com/kuoss/lethe/letheql/model"
	"github.com/kuoss/lethe/storage/logservice/logmodel"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.NotEmpty(t, queryService)
}

func TestQuery(t *testing.T) {
	testCases := []struct {
		qs        string
		tr        model.TimeRange
		wantError string
		want      *letheql.Result
	}{
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
				logmodel.PodLog{Time: "2009-11-10T22:59:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "lerom ipsum"},
				logmodel.PodLog{Time: "2009-11-10T22:59:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"}}}, Warnings: model.Warnings(nil)},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			res := queryService.Query(context.TODO(), tc.qs, tc.tr)
			if tc.wantError == "" {
				assert.NoError(t, res.Err)
			} else {
				assert.EqualError(t, res.Err, tc.wantError)
			}
			res.Err = nil
			assert.Equal(t, tc.want, res)
		})
	}
}
