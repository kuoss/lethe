package parser

import (
	"fmt"
	"testing"

	commonModel "github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {

	testCases := []struct {
		input     string
		wantError string
		want      Expr
	}{
		// BinaryExpr - single FilterOperator
		{
			`pod|="hello"`,
			"",
			&BinaryExpr{
				Op: PIPE_EQL,
				LHS: &VectorSelector{
					Name: "pod", LabelMatchers: []*labels.Matcher{MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: PositionRange{Start: 0, End: 3}},
				RHS: &StringLiteral{
					Val:      "hello",
					PosRange: PositionRange{Start: 5, End: 12}},
				ReturnBool: false},
		},
		{
			`pod|~"hel.*"`,
			"",
			&BinaryExpr{
				Op: PIPE_REGEX,
				LHS: &VectorSelector{
					Name: "pod", LabelMatchers: []*labels.Matcher{MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: PositionRange{Start: 0, End: 3}},
				RHS: &StringLiteral{
					Val:      "hel.*",
					PosRange: PositionRange{Start: 5, End: 12}},
				ReturnBool: false},
		},
		{
			`pod!="hello"`,
			"",
			&BinaryExpr{
				Op: NEQ,
				LHS: &VectorSelector{
					Name: "pod", LabelMatchers: []*labels.Matcher{MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: PositionRange{Start: 0, End: 3}},
				RHS: &StringLiteral{
					Val:      "hello",
					PosRange: PositionRange{Start: 5, End: 12}},
				ReturnBool: false},
		},
		{
			`pod!~"hel.*"`,
			"",
			&BinaryExpr{
				Op: NEQ_REGEX,
				LHS: &VectorSelector{
					Name: "pod", LabelMatchers: []*labels.Matcher{MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: PositionRange{Start: 0, End: 3}},
				RHS: &StringLiteral{
					Val:      "hel.*",
					PosRange: PositionRange{Start: 5, End: 12}},
				ReturnBool: false},
		},
		{
			`pod |= "hello"`,
			"",
			&BinaryExpr{
				Op: PIPE_EQL,
				LHS: &VectorSelector{
					Name: "pod", LabelMatchers: []*labels.Matcher{MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: PositionRange{Start: 0, End: 3}},
				RHS: &StringLiteral{
					Val:      "hello",
					PosRange: PositionRange{Start: 7, End: 14}},
				ReturnBool: false},
		},
		{
			`pod{} |= "hello"`,
			"",
			&BinaryExpr{
				Op: PIPE_EQL,
				LHS: &VectorSelector{
					Name: "pod",
					LabelMatchers: []*labels.Matcher{
						MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: PositionRange{Start: 0, End: 5}},
				RHS: &StringLiteral{
					Val:      "hello",
					PosRange: PositionRange{Start: 9, End: 16}},
				ReturnBool: false},
		},
		{
			`pod{} |~ "hello.*"`,
			"",
			&BinaryExpr{
				Op: PIPE_REGEX,
				LHS: &VectorSelector{
					Name: "pod", LabelMatchers: []*labels.Matcher{MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: PositionRange{Start: 0, End: 5}},
				RHS: &StringLiteral{
					Val:      "hello.*",
					PosRange: PositionRange{Start: 9, End: 18}},
				ReturnBool: false},
		},
		{
			`pod{namespace="namespace01"} |= "hello"`,
			"",
			&BinaryExpr{
				Op: PIPE_EQL,
				LHS: &VectorSelector{
					Name: "pod", LabelMatchers: []*labels.Matcher{
						MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
						MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: PositionRange{Start: 0, End: 28}},
				RHS: &StringLiteral{
					Val:      "hello",
					PosRange: PositionRange{Start: 32, End: 39}},
				ReturnBool: false},
		},
		{
			`pod{namespace="namespace01"} |~ "hel.*"`,
			"",
			&BinaryExpr{
				Op: PIPE_REGEX,
				LHS: &VectorSelector{
					Name: "pod", LabelMatchers: []*labels.Matcher{
						MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
						MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: PositionRange{Start: 0, End: 28}},
				RHS: &StringLiteral{
					Val:      "hel.*",
					PosRange: PositionRange{Start: 32, End: 39}},
				ReturnBool: false},
		},
		// BinaryExpr - multi FilterOperator (nested)
		{
			`pod|="hello"!="world"`,
			"",
			&BinaryExpr{
				Op: NEQ,
				LHS: &BinaryExpr{
					Op: PIPE_EQL,
					LHS: &VectorSelector{
						Name: "pod", LabelMatchers: []*labels.Matcher{MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
						PosRange: PositionRange{Start: 0, End: 3}},
					RHS: &StringLiteral{
						Val:      "hello",
						PosRange: PositionRange{Start: 5, End: 12}},
					ReturnBool: false},
				RHS: &StringLiteral{
					Val:      "world",
					PosRange: PositionRange{Start: 14, End: 21}},
				ReturnBool: false},
		},
		{
			`pod|~"hel.*"|="world"`,
			"",
			&BinaryExpr{
				Op: PIPE_EQL,
				LHS: &BinaryExpr{
					Op: PIPE_REGEX,
					LHS: &VectorSelector{
						Name: "pod", LabelMatchers: []*labels.Matcher{MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
						PosRange: PositionRange{Start: 0, End: 3}},
					RHS: &StringLiteral{
						Val:      "hel.*",
						PosRange: PositionRange{Start: 5, End: 12}},
					ReturnBool: false},
				RHS: &StringLiteral{
					Val:      "world",
					PosRange: PositionRange{Start: 14, End: 21}},
				ReturnBool: false},
		},
		{
			`pod|~"hel.*"!~"wor.*"`,
			"",
			&BinaryExpr{
				Op: NEQ_REGEX,
				LHS: &BinaryExpr{
					Op: PIPE_REGEX,
					LHS: &VectorSelector{
						Name: "pod", LabelMatchers: []*labels.Matcher{MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
						PosRange: PositionRange{Start: 0, End: 3}},
					RHS: &StringLiteral{
						Val:      "hel.*",
						PosRange: PositionRange{Start: 5, End: 12}},
					ReturnBool: false},
				RHS: &StringLiteral{
					Val:      "wor.*",
					PosRange: PositionRange{Start: 14, End: 21}},
				ReturnBool: false},
		},

		// NumberLiteral
		{
			`42`,
			"",
			&NumberLiteral{
				Val:      42,
				PosRange: PositionRange{Start: 0, End: 2}},
		},

		{
			`"hello"`,
			"",
			&StringLiteral{
				Val:      "hello",
				PosRange: PositionRange{Start: 0, End: 7}},
		},

		// VectorSelector
		{
			`pod`,
			"",
			&VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: PositionRange{Start: 0, End: 3}},
		},
		{
			`pod{}`,
			"",
			&VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: PositionRange{Start: 0, End: 5}},
		},
		{
			`pod{namespace="namespace01"}`,
			"",
			&VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
					MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: PositionRange{Start: 0, End: 28}},
		},
		{
			`pod{namespace="not-exists"}`,
			"",
			&VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					MustLabelMatcher(labels.MatchEqual, "namespace", "not-exists"),
					MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: PositionRange{Start: 0, End: 27}},
		},
		{
			`pod{namespace="namespace01",pod="nginx"}`,
			"",
			&VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
					MustLabelMatcher(labels.MatchEqual, "pod", "nginx"),
					MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: PositionRange{Start: 0, End: 40}},
		},
		{
			`pod{namespace="namespace01",pod="nginx-*"}`,
			"",
			&VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
					MustLabelMatcher(labels.MatchEqual, "pod", "nginx-*"),
					MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: PositionRange{Start: 0, End: 42}},
		},
		{
			`pod{namespace="namespace01",container="nginx"}`,
			"",
			&VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
					MustLabelMatcher(labels.MatchEqual, "container", "nginx"),
					MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: PositionRange{Start: 0, End: 46}},
		},
		{
			`pod{namespace="namespace*",container="nginx"}`,
			"",
			&VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					MustLabelMatcher(labels.MatchEqual, "namespace", "namespace*"),
					MustLabelMatcher(labels.MatchEqual, "container", "nginx"),
					MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: PositionRange{Start: 0, End: 45}},
		},

		// MatrixSelector
		{
			`pod{namespace="namespace01",pod="nginx-*"}[3m]`,
			"",
			&MatrixSelector{
				VectorSelector: &VectorSelector{
					Name: "pod",
					LabelMatchers: []*labels.Matcher{
						MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
						MustLabelMatcher(labels.MatchEqual, "pod", "nginx-*"),
						MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: PositionRange{Start: 0, End: 42}},
				Range: 180000000000, EndPos: 46},
		},

		// Call
		{
			`count_over_time(pod{namespace="namespace01",pod="nginx-*"}[3m])`,
			"",
			&Call{
				Func: &Function{
					Name:       "count_over_time",
					ArgTypes:   []ValueType{ValueTypeMatrix},
					ReturnType: ValueTypeVector},
				Args: Expressions{
					&MatrixSelector{
						VectorSelector: &VectorSelector{
							Name: "pod",
							LabelMatchers: []*labels.Matcher{
								MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
								MustLabelMatcher(labels.MatchEqual, "pod", "nginx-*"),
								MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
							PosRange: PositionRange{Start: 16, End: 58}},
						Range: 180000000000, EndPos: 62}},
				PosRange: PositionRange{Start: 0, End: 63}},
		},
		{
			`count_over_time(pod{}[3m])`,
			"",
			&Call{
				Func: &Function{
					Name:       "count_over_time",
					ArgTypes:   []ValueType{ValueTypeMatrix},
					ReturnType: ValueTypeVector},
				Args: Expressions{
					&MatrixSelector{
						VectorSelector: &VectorSelector{
							Name:          "pod",
							LabelMatchers: []*labels.Matcher{MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
							PosRange:      PositionRange{Start: 16, End: 21}},
						Range: 180000000000, EndPos: 25}},
				PosRange: PositionRange{Start: 0, End: 26}},
		},

		// BinaryExpr
		{
			`count_over_time(pod{}[3m]) > 10`,
			"",
			&BinaryExpr{
				Op: GTR,
				LHS: &Call{
					Func: &Function{
						Name:       "count_over_time",
						ArgTypes:   []ValueType{ValueTypeMatrix},
						ReturnType: ValueTypeVector},
					Args: Expressions{
						&MatrixSelector{
							VectorSelector: &VectorSelector{
								Name:          "pod",
								LabelMatchers: []*labels.Matcher{MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
								PosRange:      PositionRange{Start: 16, End: 21}},
							Range: 180000000000, EndPos: 25}},
					PosRange: PositionRange{Start: 0, End: 26}},
				RHS: &NumberLiteral{
					Val:      10,
					PosRange: PositionRange{Start: 29, End: 31}}},
		},
		{
			`count_over_time(pod{}[3m]) < 10`,
			"",
			&BinaryExpr{
				Op: LSS,
				LHS: &Call{
					Func: &Function{
						Name:       "count_over_time",
						ArgTypes:   []ValueType{ValueTypeMatrix},
						ReturnType: ValueTypeVector},
					Args: Expressions{
						&MatrixSelector{
							VectorSelector: &VectorSelector{
								Name:          "pod",
								LabelMatchers: []*labels.Matcher{MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
								PosRange:      PositionRange{Start: 16, End: 21}},
							Range: 180000000000, EndPos: 25}},
					PosRange: PositionRange{Start: 0, End: 26}},
				RHS: &NumberLiteral{
					Val:      10,
					PosRange: PositionRange{Start: 29, End: 31}}},
		},
		{
			`count_over_time(pod{}[3m]) == 21`,
			"",
			&BinaryExpr{
				Op: EQLC,
				LHS: &Call{
					Func: &Function{
						Name:       "count_over_time",
						ArgTypes:   []ValueType{ValueTypeMatrix},
						ReturnType: ValueTypeVector},
					Args: Expressions{
						&MatrixSelector{
							VectorSelector: &VectorSelector{
								Name:          "pod",
								LabelMatchers: []*labels.Matcher{MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
								PosRange:      PositionRange{Start: 16, End: 21}},
							Range: 180000000000, EndPos: 25}},
					PosRange: PositionRange{Start: 0, End: 26}},
				RHS: &NumberLiteral{
					Val:      21,
					PosRange: PositionRange{Start: 30, End: 32}},
				VectorMatching: (*VectorMatching)(nil),
				ReturnBool:     false},
		},

		// ######## ERROR
		{
			`pod{namespace="namespace01"} "`,
			"1:30: parse error: unterminated quoted string",
			&VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
					MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: PositionRange{Start: 0, End: 28},
			},
		},
		{
			`pod{namespace="namespace01"} hello`,
			"1:30: parse error: unexpected identifier \"hello\"",
			&VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
					MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: PositionRange{Start: 0, End: 28},
			},
		},
		{
			`pod{namespace="namespace01"} "hello`,
			"1:30: parse error: unterminated quoted string",
			&VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
					MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: PositionRange{Start: 0, End: 28},
			},
		},
		{
			`pod{namespace="namespace01"} "hello"`,
			"1:30: parse error: unexpected string \"\\\"hello\\\"\"",
			&VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
					MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: PositionRange{Start: 0, End: 28},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			expr, err := ParseExpr(tc.input)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			assert.Equal(t, tc.want, expr)
		})
	}
}
