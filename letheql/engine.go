package letheql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/letheql/model"
	"github.com/kuoss/lethe/letheql/parser"
	"github.com/kuoss/lethe/storage/logservice"
)

type Engine struct {
	cfg        *config.Config
	logService *logservice.LogService
}

func NewEngine(cfg *config.Config, logService *logservice.LogService) *Engine {
	return &Engine{cfg, logService}
}

func (ng *Engine) NewRangeQuery(qs string, start, end time.Time) (Query, error) {
	logger.Debugf("newRangeQuery qs: %s", qs)
	expr, err := parser.ParseExpr(qs)
	if err != nil {
		return nil, err
	}
	qry := &query{
		q:    qs,
		stmt: &parser.EvalStmt{Expr: expr, Start: start, End: end},
		ng:   ng,
	}
	return qry, nil
}

func (ng *Engine) exec(ctx context.Context, q *query) (v parser.Value, ws model.Warnings, err error) {
	ctx, cancel := context.WithTimeout(ctx, ng.cfg.Query.Timeout)
	q.cancel = cancel
	defer q.cancel()

	switch s := q.Statement().(type) {
	case *parser.EvalStmt:
		return ng.execEvalStmt(ctx, s)
	case parser.TestStmt:
		return nil, nil, s(ctx)
	}
	panic(fmt.Errorf("letheql.exec: unhandled statement of type %T", q.Statement()))
}

func (ng *Engine) execEvalStmt(ctx context.Context, stmt *parser.EvalStmt) (parser.Value, model.Warnings, error) {
	// Range evaluation.
	evaluator := &evaluator{
		logService:     ng.logService,
		start:          stmt.Start,
		end:            stmt.End,
		startTimestamp: timeMilliseconds(stmt.Start),
		endTimestamp:   timeMilliseconds(stmt.End),
		ctx:            ctx,
	}
	val, warnings, err := evaluator.Eval(stmt.Expr)
	if err != nil {
		return nil, warnings, err
	}
	switch result := val.(type) {
	case model.Log:
		return result, warnings, nil
	default:
		panic(fmt.Errorf("promql.Engine.exec: invalid expression type %q", val.Type()))
	}
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

func timeMilliseconds(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond/time.Nanosecond)
}
