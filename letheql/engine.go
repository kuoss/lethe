package letheql

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/kuoss/lethe/letheql/model"
	"github.com/kuoss/lethe/letheql/parser"
	"github.com/kuoss/lethe/storage/logservice"
	"github.com/prometheus/prometheus/model/timestamp"
	"github.com/prometheus/prometheus/storage"
)

type Engine struct {
	logService *logservice.LogService
}

func New(logService *logservice.LogService) *Engine {
	return &Engine{
		logService: logService,
	}
}

func (ng *Engine) NewInstantQuery(_ context.Context, q storage.Queryable, qs string, ts time.Time) (Query, error) {
	expr, err := parser.ParseExpr(qs)
	if err != nil {
		return nil, err
	}
	qry, err := ng.newQuery(q, expr, ts, ts, 0)
	if err != nil {
		return nil, err
	}
	qry.q = qs

	return qry, nil
}

func (ng *Engine) NewRangeQuery(_ context.Context, q storage.Queryable, qs string, start, end time.Time, interval time.Duration) (Query, error) {
	expr, err := parser.ParseExpr(qs)
	if err != nil {
		return nil, err
	}
	qry, err := ng.newQuery(q, expr, start, end, interval)
	if err != nil {
		return nil, err
	}
	qry.q = qs

	return qry, nil
}

func (ng *Engine) newQuery(q storage.Queryable, expr parser.Expr, start, end time.Time, interval time.Duration) (*query, error) {
	es := &parser.EvalStmt{
		Expr:     PreprocessExpr(expr, start, end),
		Start:    start,
		End:      end,
		Interval: interval,
	}
	qry := &query{
		stmt:      es,
		ng:        ng,
		queryable: q,
	}
	return qry, nil
}

func (ng *Engine) exec(ctx context.Context, q *query) (v parser.Value, ws model.Warnings, err error) {
	ctx, cancel := context.WithTimeout(ctx, ng.logService.Config().Timeout())
	q.cancel = cancel
	defer q.cancel()

	switch s := q.Statement().(type) {
	case *parser.EvalStmt:
		fmt.Println("engine: exec 4")
		return ng.execEvalStmt(ctx, q, s)
	case parser.TestStmt:
		return nil, nil, s(ctx)
	}
	panic(fmt.Errorf("queryengine.exec: unhandled statement of type %T", q.Statement()))
}

func (ng *Engine) execEvalStmt(ctx context.Context, query *query, s *parser.EvalStmt) (parser.Value, model.Warnings, error) {
	fmt.Println("engine: execEvalStmt 1")
	mint, maxt := ng.findMinMaxTime(s)
	fmt.Println("engine: execEvalStmt 2")
	querier, err := query.queryable.Querier(ctx, mint, maxt)
	if err != nil {
		return nil, nil, fmt.Errorf("querier err: %w", err)
	}
	defer querier.Close()

	fmt.Println("engine: execEvalStmt 3")
	if s.Start == s.End && s.Interval == 0 {
		start := timeMilliseconds(s.Start)
		evaluator := &evaluator{
			logService:     ng.logService,
			start:          s.Start,
			end:            s.Start,
			startTimestamp: start,
			endTimestamp:   start,
			interval:       1,
			ctx:            ctx,
		}
		val, warnings, err := evaluator.Eval(s.Expr)
		if err != nil {
			return nil, warnings, err
		}

		var mat Matrix

		switch result := val.(type) {
		case Matrix:
			mat = result
		case String:
			return result, warnings, nil
		default:
			panic(fmt.Errorf("queryengine.exec: invalid expression type %q", val.Type()))
		}

		query.matrix = mat
		switch s.Expr.Type() {
		case parser.ValueTypeVector:
			// Convert matrix with one value per series into vector.
			vector := make(Vector, len(mat))
			for i, s := range mat {
				// Point might have a different timestamp, force it to the evaluation
				// timestamp as that is when we ran the evaluation.
				if len(s.Histograms) > 0 {
					vector[i] = Sample{Metric: s.Metric, H: s.Histograms[0].H, T: start}
				} else {
					vector[i] = Sample{Metric: s.Metric, F: s.Floats[0].F, T: start}
				}
			}
			return vector, warnings, nil
		case parser.ValueTypeScalar:
			return Scalar{V: mat[0].Floats[0].F, T: start}, warnings, nil
		case parser.ValueTypeMatrix:
			return mat, warnings, nil
		default:
			panic(fmt.Errorf("queryengine.exec: unexpected expression type %q", s.Expr.Type()))
		}
	}

	// Range evaluation.
	evaluator := &evaluator{
		logService:     ng.logService,
		start:          s.Start,
		end:            s.End,
		startTimestamp: timeMilliseconds(s.Start),
		endTimestamp:   timeMilliseconds(s.End),
		interval:       durationMilliseconds(s.Interval),
		ctx:            ctx,
	}
	val, warnings, err := evaluator.Eval(s.Expr)
	if err != nil {
		return nil, warnings, err
	}

	mat, ok := val.(Matrix)
	if !ok {
		panic(fmt.Errorf("queryengine.exec: invalid expression type %q", val.Type()))
	}
	query.matrix = mat

	if err := contextDone(ctx, "expression evaluation"); err != nil {
		return nil, warnings, err
	}
	return mat, warnings, nil
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

func PreprocessExpr(expr parser.Expr, start, end time.Time) parser.Expr {
	isStepInvariant := preprocessExprHelper(expr, start, end)
	if isStepInvariant {
		return newStepInvariantExpr(expr)
	}
	return expr
}

func preprocessExprHelper(expr parser.Expr, start, end time.Time) bool {
	switch n := expr.(type) {
	case *parser.VectorSelector:
		switch n.StartOrEnd {
		case parser.START:
			n.Timestamp = makeInt64Pointer(timestamp.FromTime(start))
		case parser.END:
			n.Timestamp = makeInt64Pointer(timestamp.FromTime(end))
		}
		return n.Timestamp != nil

	case *parser.AggregateExpr:
		return preprocessExprHelper(n.Expr, start, end)

	case *parser.BinaryExpr:
		isInvariant1, isInvariant2 := preprocessExprHelper(n.LHS, start, end), preprocessExprHelper(n.RHS, start, end)
		if isInvariant1 && isInvariant2 {
			return true
		}

		if isInvariant1 {
			n.LHS = newStepInvariantExpr(n.LHS)
		}
		if isInvariant2 {
			n.RHS = newStepInvariantExpr(n.RHS)
		}

		return false

	case *parser.Call:
		_, ok := AtModifierUnsafeFunctions[n.Func.Name]
		isStepInvariant := !ok
		isStepInvariantSlice := make([]bool, len(n.Args))
		for i := range n.Args {
			isStepInvariantSlice[i] = preprocessExprHelper(n.Args[i], start, end)
			isStepInvariant = isStepInvariant && isStepInvariantSlice[i]
		}

		if isStepInvariant {
			// The function and all arguments are step invariant.
			return true
		}

		for i, isi := range isStepInvariantSlice {
			if isi {
				n.Args[i] = newStepInvariantExpr(n.Args[i])
			}
		}
		return false

	case *parser.MatrixSelector:
		return preprocessExprHelper(n.VectorSelector, start, end)

	case *parser.SubqueryExpr:
		// Since we adjust offset for the @ modifier evaluation,
		// it gets tricky to adjust it for every subquery step.
		// Hence we wrap the inside of subquery irrespective of
		// @ on subquery (given it is also step invariant) so that
		// it is evaluated only once w.r.t. the start time of subquery.
		isInvariant := preprocessExprHelper(n.Expr, start, end)
		if isInvariant {
			n.Expr = newStepInvariantExpr(n.Expr)
		}
		switch n.StartOrEnd {
		case parser.START:
			n.Timestamp = makeInt64Pointer(timestamp.FromTime(start))
		case parser.END:
			n.Timestamp = makeInt64Pointer(timestamp.FromTime(end))
		}
		return n.Timestamp != nil

	case *parser.ParenExpr:
		return preprocessExprHelper(n.Expr, start, end)

	case *parser.UnaryExpr:
		return preprocessExprHelper(n.Expr, start, end)

	case *parser.StringLiteral, *parser.NumberLiteral:
		return true
	}

	panic(fmt.Sprintf("found unexpected node %#v", expr))
}

func makeInt64Pointer(val int64) *int64 {
	valp := new(int64)
	*valp = val
	return valp
}

func newStepInvariantExpr(expr parser.Expr) parser.Expr {
	return &parser.StepInvariantExpr{Expr: expr}
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
