package model

import (
	"fmt"
	"time"

	"github.com/kuoss/lethe/letheql/parser"
	"github.com/prometheus/prometheus/model/labels"
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
	Name          string
	LabelMatchers []*labels.Matcher
	LineMatchers  []*LineMatcher
	TimeRange     TimeRange
}

const ValueTypeLogSelector parser.ValueType = "logselector"

// implements parser.Expr
func (ls LogSelector) PromQLExpr()                         {}
func (ls LogSelector) Pretty(int) string                   { return ls.String() }
func (ls LogSelector) PositionRange() parser.PositionRange { return parser.PositionRange{} }

// implements parser.Value
func (ls LogSelector) Type() parser.ValueType { return ValueTypeLogSelector }
func (ls LogSelector) String() string         { return fmt.Sprintf("%#v", ls) }
