package letheql

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/letheql/model"
	"github.com/kuoss/lethe/letheql/parser"
	"github.com/kuoss/lethe/storage/logservice"
)

type evaluator struct {
	logService *logservice.LogService

	ctx context.Context

	startTimestamp int64 // Start time in milliseconds.
	endTimestamp   int64 // End time in milliseconds.

	start time.Time
	end   time.Time
}

func (ev *evaluator) error(err error) {
	panic(err)
}

func (ev *evaluator) recover(expr parser.Expr, ws *model.Warnings, errp *error) {
	e := recover()
	if e == nil {
		return
	}

	switch err := e.(type) {
	case runtime.Error:
		buf := make([]byte, 64<<10)
		buf = buf[:runtime.Stack(buf, false)]

		logger.Errorf("msg: runtime panic in parser. expr: %s, err: %s, stacktrace: %s", expr.String(), e, string(buf))
		*errp = fmt.Errorf("unexpected error: %w", err)
	case model.ErrWithWarnings:
		*errp = err.Err
		*ws = append(*ws, err.Warnings...)
	case error:
		*errp = err
	default:
		*errp = fmt.Errorf("%v", err)
	}
}

func (ev *evaluator) Eval(expr parser.Expr) (v parser.Value, ws model.Warnings, err error) {
	defer ev.recover(expr, &ws, &err)
	val, ws := ev.eval(expr)

	ls, ok := val.(*model.LogSelector)
	if ok {
		val, ws = ev.evalWithWarnings(ls, &ws)
	}
	return val, ws, nil
}

func (ev *evaluator) eval(expr parser.Expr) (parser.Value, model.Warnings) {
	if err := contextDone(ev.ctx, "expression evaluation"); err != nil {
		ev.error(err)
	}

	switch e := expr.(type) {
	case *parser.BinaryExpr:
		return ev.evalBinaryExpr(e)
	case *parser.StringLiteral:
		return String{V: e.Val, T: ev.startTimestamp}, nil
	case *parser.VectorSelector:
		return ev.evalVectorSelector(e)
	case *model.LogSelector:
		return ev.evalLogSelector(e)
	}
	panic(fmt.Errorf("eval: unhandled expr: %#v", expr))
}

func (ev *evaluator) evalWithWarnings(expr parser.Expr, warnings *model.Warnings) (parser.Value, model.Warnings) {
	val, ws := ev.eval(expr)
	*warnings = append(*warnings, ws...)
	return val, *warnings
}

func (ev *evaluator) evalVectorSelector(vs *parser.VectorSelector) (*model.LogSelector, model.Warnings) {
	return &model.LogSelector{
		Name:          vs.Name,
		LabelMatchers: vs.LabelMatchers,
		TimeRange:     model.TimeRange{Start: ev.start, End: ev.end},
	}, nil
}

func (ev *evaluator) evalLogSelector(ls *model.LogSelector) (parser.Value, model.Warnings) {
	val, ws, err := ev.logService.SelectLog(ls)
	if err != nil {
		ev.error(err)
	}
	return val, ws
}

func (ev *evaluator) evalBinaryExpr(expr *parser.BinaryExpr) (parser.Value, model.Warnings) {
	if !expr.Op.IsFilterOperator() {
		ev.error(fmt.Errorf("evalBinaryExpr err: not filter operator: %s", expr.Op))
	}

	switch lhs := expr.LHS.(type) {
	case *parser.BinaryExpr:
		newLHS, warnings := ev.eval(lhs)
		switch nl := newLHS.(type) {
		case *model.LogSelector:
			expr.LHS = nl
			return ev.evalWithWarnings(expr, &warnings)
		}
	case *parser.VectorSelector:
		newLHS, warnings := ev.evalVectorSelector(lhs)
		expr.LHS = newLHS
		return ev.evalWithWarnings(expr, &warnings)
	case *model.LogSelector:
		switch rhs := expr.RHS.(type) {
		case *parser.StringLiteral:
			lhs.LineMatchers = append(lhs.LineMatchers, &model.LineMatcher{Op: expr.Op, Value: rhs.Val})
			return lhs, nil
		default:
			ev.error(fmt.Errorf("unknown type rhs: %#v", rhs))
		}
	}
	panic(fmt.Errorf("evalBinaryExpr unhandles op:%s, lhs:%s, rhs:%s", expr.Op, expr.LHS, expr.RHS))
}
