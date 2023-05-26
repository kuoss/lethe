package letheql

import (
	"context"

	"github.com/kuoss/lethe/letheql/parser"
	"github.com/prometheus/prometheus/storage"
)

type Query interface {
	Cancel()
	Close()
	Exec(ctx context.Context) *Result
	Statement() parser.Statement
	String() string
}

type query struct {
	q         string
	queryable storage.Queryable
	stmt      parser.Statement
	cancel    func()
	ng        *Engine
}

func (q *query) Cancel() {
	if q.cancel != nil {
		q.cancel()
	}
}

func (q *query) Close() {}

func (q *query) Exec(ctx context.Context) *Result {
	res, warnings, err := q.ng.exec(ctx, q)
	return &Result{Err: err, Value: res, Warnings: warnings}
}

func (q *query) Statement() parser.Statement {
	return q.stmt
}

func (q *query) String() string {
	return q.q
}
