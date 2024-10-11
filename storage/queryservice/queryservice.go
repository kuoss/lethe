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
)

type QueryService struct {
	engine *letheql.Engine
}

func New(cfg *config.Config, logService *logservice.LogService) *QueryService {
	return &QueryService{
		engine: letheql.NewEngine(cfg, logService),
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
	qry, err := s.engine.NewRangeQuery(qs, tr.Start, tr.End)
	if err != nil {
		return &letheql.Result{Err: err}
	}
	res := qry.Exec(ctx)
	if res.Err != nil {
		logger.Errorf("exec err: %s", res.Err.Error())
	}
	return res
}
