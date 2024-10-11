package letheql

import (
	"context"

	"github.com/kuoss/lethe/letheql/parser"
)

type Query interface {
	Cancel()
	Exec(ctx context.Context) *Result
	Statement() parser.Statement
	String() string
}

type query struct {
	q      string
	stmt   parser.Statement
	cancel func()
	ng     *Engine
}

func (q *query) Cancel() {
	if q.cancel != nil {
		q.cancel()
	}
}

func (q *query) Exec(ctx context.Context) *Result {
	value, warnings, err := q.ng.exec(ctx, q)
	return &Result{err, value, warnings}
}

func (q *query) Statement() parser.Statement {
	return q.stmt
}

func (q *query) String() string {
	return q.q
}
