package model

import (
	"fmt"
	"time"

	"github.com/kuoss/lethe/letheql/parser"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser/posrange"
)

type TimeRange struct {
	Start time.Time
	End   time.Time
}

type LineMatcher struct {
	Op    parser.ItemType
	Value string
}

// LogSelector represents a Log selection.
type LogSelector struct {
	Name          string            // pod
	LabelMatchers []*labels.Matcher // {namespace="default",pod="nginx"}
	LineMatchers  []*LineMatcher    // |= "err" |~ "4.."
	TimeRange     TimeRange
}

const ValueTypeLogSelector parser.ValueType = "logselector"

// implements parser.Expr
func (e *LogSelector) PromQLExpr()                           {}
func (e *LogSelector) Pretty(int) string                     { return e.String() }
func (e *LogSelector) PositionRange() posrange.PositionRange { return posrange.PositionRange{} }

// implements parser.Value
func (e *LogSelector) Type() parser.ValueType { return ValueTypeLogSelector }
func (e *LogSelector) String() string         { return fmt.Sprintf("%#v", e) }
