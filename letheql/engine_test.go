package letheql

import (
	"context"
	"testing"
	"time"

	"github.com/kuoss/common/tester"
	"github.com/kuoss/lethe/clock"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/letheql/model"
	"github.com/kuoss/lethe/storage/fileservice"
	"github.com/kuoss/lethe/storage/logservice"
	"github.com/kuoss/lethe/storage/logservice/logmodel"
	"github.com/kuoss/lethe/storage/querier"
	"github.com/prometheus/prometheus/storage"
	"github.com/stretchr/testify/require"
)

func newEngine(t *testing.T) *Engine {
	cfg, err := config.New("test")
	require.NoError(t, err)

	fileService, err := fileservice.New(cfg)
	require.NoError(t, err)

	logService := logservice.New(cfg, fileService)
	return NewEngine(cfg, logService)
}

func TestNewInstantQuery_QueryableFunc(t *testing.T) {
	queryable := storage.QueryableFunc(func(ctx context.Context, mint, maxt int64) (storage.Querier, error) {
		return nil, nil
	})

	en := newEngine(t)
	qry, err := en.NewInstantQuery(context.TODO(), queryable, "pod", time.Unix(1, 0))
	require.NoError(t, err)
	require.NotEmpty(t, qry)
}

func TestNewInstantQuery_LetheQueryable(t *testing.T) {
	queryable := &querier.LetheQueryable{LetheQuerier: &querier.LetheQuerier{}}
	testCases := []struct {
		qs        string
		wantError bool
		want      *Result
	}{
		{
			qs:        `pod`,
			wantError: true,
			want:      &Result{},
		},
		{
			qs:   `pod{namespace="namespace01"}`,
			want: &Result{Value: model.Log{Name: "pod", Lines: []model.LogLine{}}, Warnings: model.Warnings(nil)},
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

			en := newEngine(t)
			qry, err := en.NewInstantQuery(context.TODO(), queryable, tc.qs, clock.Now())
			require.NoError(t, err)
			got := qry.Exec(context.TODO())
			err = got.Err
			got.Err = nil
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.want, got)
		})
	}
}

func TestNewRangeQuery(t *testing.T) {
	clock.SetPlaygroundMode(true)
	defer clock.SetPlaygroundMode(false)

	ago10d := clock.Now().Add(-240 * time.Hour)
	ago2m := clock.Now().Add(-2 * time.Minute)
	now := clock.Now()

	queryable := &querier.LetheQueryable{LetheQuerier: &querier.LetheQuerier{}}
	testCases := []struct {
		qs        string
		start     time.Time
		end       time.Time
		wantError bool
		want      *Result
	}{
		{
			qs:        `pod`,
			start:     ago10d,
			end:       now,
			wantError: true,
			want:      &Result{},
		},
		{
			qs:    `pod{namespace="namespace01"}`,
			start: ago2m,
			end:   now,
			want: &Result{Err: error(nil), Value: model.Log{Name: "pod", Lines: []model.LogLine{
				logmodel.PodLog{Time: "2009-11-10T22:59:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "lerom ipsum"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "lerom from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:59:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"}}}, Warnings: model.Warnings(nil)},
		},
		{
			qs:    `pod{namespace="namespace01"}`,
			start: ago10d,
			end:   now,
			want: &Result{Err: error(nil), Value: model.Log{Name: "pod", Lines: []model.LogLine{
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
				logmodel.PodLog{Time: "2009-11-10T22:59:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"}}}, Warnings: model.Warnings(nil)},
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i), func(t *testing.T) {
			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			ctx := context.TODO()
			engine := newEngine(t)
			qry, err := engine.NewRangeQuery(ctx, queryable, tc.qs, tc.start, tc.end)
			require.NoError(t, err)
			got := qry.Exec(ctx)
			err = got.Err
			got.Err = nil
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.want, got)
		})
	}

}
