package queryservice

import (
	"context"
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
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i), func(t *testing.T) {
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
