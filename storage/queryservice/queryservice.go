package queryservice

import (
	"context"
	"reflect"
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/clock"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/letheql"
	"github.com/kuoss/lethe/letheql/model"
	"github.com/kuoss/lethe/storage/logservice"
	"github.com/kuoss/lethe/storage/querier"
)

type QueryService struct {
	engine    *letheql.Engine
	queryable *querier.LetheQueryable
}

func New(cfg *config.Config, logService *logservice.LogService) *QueryService {
	return &QueryService{
		engine: letheql.NewEngine(cfg, logService),
		queryable: &querier.LetheQueryable{
			LetheQuerier: &querier.LetheQuerier{},
		},
	}
}

func (s *QueryService) Query(ctx context.Context, qs string, tr model.TimeRange) *letheql.Result {
	if reflect.ValueOf(tr).IsZero() {
		now := clock.Now()
		tr = model.TimeRange{
			Start: now.Add(-1 * time.Minute),
			End:   now,
		}
	}
	qry, err := s.engine.NewRangeQuery(ctx, s.queryable, qs, tr.Start, tr.End, 0)
	if err != nil {
		return &letheql.Result{Err: err}
	}
	res := qry.Exec(ctx)
	if res.Err != nil {
		logger.Errorf("exec err: %s", res.Err.Error())
	}
	return res
}
