package letheql

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/letheql/model"
	"github.com/kuoss/lethe/letheql/parser"
	"github.com/kuoss/lethe/storage/logservice"
	"github.com/prometheus/prometheus/model/timestamp"
	"github.com/prometheus/prometheus/storage"
)

type Engine struct {
	cfg        *config.Config
	logService *logservice.LogService
}

func NewEngine(cfg *config.Config, logService *logservice.LogService) *Engine {
	return &Engine{cfg, logService}
}

func (ng *Engine) NewInstantQuery(_ context.Context, q storage.Queryable, qs string, ts time.Time) (Query, error) {
	expr, err := parser.ParseExpr(qs)
	if err != nil {
		return nil, err
	}
	qry, err := ng.newQuery(q, expr, ts, ts)
	if err != nil {
		return nil, err
	}
	qry.q = qs

	return qry, nil
}

func (ng *Engine) NewRangeQuery(_ context.Context, q storage.Queryable, qs string, start, end time.Time) (Query, error) {
	logger.Infof("newRangeQuery qs: %s", qs)
	expr, err := parser.ParseExpr(qs)
	if err != nil {
		return nil, err
	}
	qry, err := ng.newQuery(q, expr, start, end)
	if err != nil {
		return nil, err
	}
	qry.q = qs

	return qry, nil
}

func (ng *Engine) newQuery(q storage.Queryable, expr parser.Expr, start, end time.Time) (*query, error) {
	es := &parser.EvalStmt{
		Expr:  expr,
		Start: start,
		End:   end,
	}
	qry := &query{
		stmt:      es,
		ng:        ng,
		queryable: q,
	}
	return qry, nil
}

func (ng *Engine) exec(ctx context.Context, q *query) (v parser.Value, ws model.Warnings, err error) {
	ctx, cancel := context.WithTimeout(ctx, ng.cfg.Query.Timeout)
	q.cancel = cancel
	defer q.cancel()
	switch s := q.Statement().(type) {
	case *parser.EvalStmt:
		return ng.execEvalStmt(ctx, q, s)
	case parser.TestStmt:
		return nil, nil, s(ctx)
	}
	panic(fmt.Errorf("letheql.exec: unhandled statement of type %T", q.Statement()))
}

func (ng *Engine) execEvalStmt(ctx context.Context, query *query, s *parser.EvalStmt) (parser.Value, model.Warnings, error) {
	mint, maxt := ng.findMinMaxTime(s)
	querier, err := query.queryable.Querier(ctx, mint, maxt)
	if err != nil {
		return nil, nil, fmt.Errorf("querier err: %w", err)
	}
	defer querier.Close()

	// Range evaluation.
	evaluator := &evaluator{
		logService:     ng.logService,
		start:          s.Start,
		end:            s.End,
		startTimestamp: timeMilliseconds(s.Start),
		endTimestamp:   timeMilliseconds(s.End),
		ctx:            ctx,
	}
	val, warnings, err := evaluator.Eval(s.Expr)
	if err != nil {
		return nil, warnings, err
	}
	switch result := val.(type) {
	case model.Log:
		return result, warnings, nil
	case String:
		return result, warnings, nil
	default:
		panic(fmt.Errorf("promql.Engine.exec: invalid expression type %q", val.Type()))
	}
}
func (ng *Engine) findMinMaxTime(s *parser.EvalStmt) (int64, int64) {
	var minTimestamp, maxTimestamp int64 = math.MaxInt64, math.MinInt64
	// Whenever a MatrixSelector is evaluated, evalRange is set to the corresponding range.
	// The evaluation of the VectorSelector inside then evaluates the given range and unsets
	// the variable.

	var evalRange time.Duration
	parser.Inspect(s.Expr, func(node parser.Node, path []parser.Node) error {
		switch n := node.(type) {
		case *parser.VectorSelector:
			start, end := ng.getTimeRangesForSelector(s, n, path, evalRange)
			if start < minTimestamp {
				minTimestamp = start
			}
			if end > maxTimestamp {
				maxTimestamp = end
			}
			evalRange = 0
		case *parser.MatrixSelector:
			evalRange = n.Range
		}
		return nil
	})

	if maxTimestamp == math.MinInt64 {
		// This happens when there was no selector. Hence no time range to select.
		minTimestamp = 0
		maxTimestamp = 0
	}

	return minTimestamp, maxTimestamp
}
func (ng *Engine) getTimeRangesForSelector(s *parser.EvalStmt, n *parser.VectorSelector, path []parser.Node, evalRange time.Duration) (int64, int64) {
	start, end := timestamp.FromTime(s.Start), timestamp.FromTime(s.End)
	subqOffset, subqRange, subqTs := subqueryTimes(path)

	if subqTs != nil {
		// The timestamp on the subquery overrides the eval statement time ranges.
		start = *subqTs
		end = *subqTs
	}

	if n.Timestamp != nil {
		// The timestamp on the selector overrides everything.
		start = *n.Timestamp
		end = *n.Timestamp
	} else {
		offsetMilliseconds := durationMilliseconds(subqOffset)
		start = start - offsetMilliseconds - durationMilliseconds(subqRange)
		end -= offsetMilliseconds
	}

	if evalRange == 0 {
		start -= durationMilliseconds(s.LookbackDelta)
	} else {
		// For all matrix queries we want to ensure that we have (end-start) + range selected
		// this way we have `range` data before the start time
		start -= durationMilliseconds(evalRange)
	}

	offsetMilliseconds := durationMilliseconds(n.OriginalOffset)
	start -= offsetMilliseconds
	end -= offsetMilliseconds

	return start, end
}

func contextDone(ctx context.Context, env string) error {
	if err := ctx.Err(); err != nil {
		return contextErr(err, env)
	}
	return nil
}

func contextErr(err error, env string) error {
	switch {
	case errors.Is(err, context.Canceled):
		return model.ErrQueryCanceled(env)
	case errors.Is(err, context.DeadlineExceeded):
		return model.ErrQueryTimeout(env)
	default:
		return err
	}
}

func subqueryTimes(path []parser.Node) (time.Duration, time.Duration, *int64) {
	var (
		subqOffset, subqRange time.Duration
		ts                    int64 = math.MaxInt64
	)
	for _, node := range path {
		if n, ok := node.(*parser.SubqueryExpr); ok {
			subqOffset += n.OriginalOffset
			subqRange += n.Range
			if n.Timestamp != nil {
				// The @ modifier on subquery invalidates all the offset and
				// range till now. Hence resetting it here.
				subqOffset = n.OriginalOffset
				subqRange = n.Range
				ts = *n.Timestamp
			}
		}
	}
	var tsp *int64
	if ts != math.MaxInt64 {
		tsp = &ts
	}
	return subqOffset, subqRange, tsp
}

func timeMilliseconds(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond/time.Nanosecond)
}

func durationMilliseconds(d time.Duration) int64 {
	return int64(d / (time.Millisecond / time.Nanosecond))
}
