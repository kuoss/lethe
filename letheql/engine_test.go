package letheql

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kuoss/lethe/clock"
	"github.com/kuoss/lethe/letheql/model"
	"github.com/kuoss/lethe/storage/logservice/logmodel"
	"github.com/kuoss/lethe/storage/querier"
	"github.com/prometheus/prometheus/storage"
	"github.com/stretchr/testify/assert"
)

func TestNewEngine(t *testing.T) {
	assert.NotNil(t, engine1)
}

func TestNewInstantQuery(t *testing.T) {

	t.Run("QueryableFunc", func(t *testing.T) {
		queryable := storage.QueryableFunc(func(ctx context.Context, mint, maxt int64) (storage.Querier, error) {
			return nil, nil
		})
		ctx, cancelCtx := context.WithCancel(context.Background())
		defer cancelCtx()

		qry, err := engine1.NewInstantQuery(ctx, queryable, "pod", time.Unix(1, 0))
		assert.NoError(t, err)
		assert.NotEmpty(t, qry)
	})

	t.Run("LetheQueryable", func(t *testing.T) {
		// LetheQueryable
		queryable := &querier.LetheQueryable{LetheQuerier: &querier.LetheQuerier{}}
		testCases := []struct {
			qs        string
			wantError string
			want      *Result
		}{
			{
				`pod`,
				"getTargets err: target matcher err: not found label 'namespace' for logType 'pod'",
				&Result{},
			},
			{
				`pod{namespace="namespace01"}`,
				"",
				&Result{Value: model.Log{Name: "pod", Lines: []model.LogLine{}}, Warnings: model.Warnings(nil)},
			},
		}
		for i, tc := range testCases {
			t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
				qry, err := engine1.NewInstantQuery(context.TODO(), queryable, tc.qs, clock.Now())
				assert.NoError(t, err)
				got := qry.Exec(context.TODO())
				err = got.Err
				got.Err = nil
				if tc.wantError == "" {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, tc.wantError)
				}
				assert.Equal(t, tc.want, got)
			})
		}
	})
}

func TestNewRangeQuery(t *testing.T) {

	ago10d := clock.Now().Add(-240 * time.Hour)
	ago2m := clock.Now().Add(-2 * time.Minute)
	now := clock.Now()

	queryable := &querier.LetheQueryable{LetheQuerier: &querier.LetheQuerier{}}
	testCases := []struct {
		qs        string
		start     time.Time
		end       time.Time
		wantError string
		want      *Result
	}{
		{
			`pod`, ago10d, now,
			"getTargets err: target matcher err: not found label 'namespace' for logType 'pod'",
			&Result{},
		},
		{
			`pod{namespace="namespace01"}`, ago2m, now,
			"",
			&Result{Err: error(nil), Value: model.Log{Name: "pod", Lines: []model.LogLine{
				logmodel.PodLog{Time: "2009-11-10T22:59:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "lerom ipsum"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "lerom from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:58:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "sidecar", Log: "hello from sidecar"},
				logmodel.PodLog{Time: "2009-11-10T22:59:00.000000Z", Namespace: "namespace01", Pod: "nginx-deployment-75675f5897-7ci7o", Container: "nginx", Log: "hello world"}}}, Warnings: model.Warnings(nil)},
		},
		{
			`pod{namespace="namespace01"}`, ago10d, now,
			"",
			&Result{Err: error(nil), Value: model.Log{Name: "pod", Lines: []model.LogLine{
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
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			ctx := context.TODO()
			qry, err := engine1.NewRangeQuery(ctx, queryable, tc.qs, tc.start, tc.end, 0)
			assert.NoError(t, err)
			got := qry.Exec(ctx)
			err = got.Err
			got.Err = nil
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			assert.Equal(t, tc.want, got)
		})
	}

}
