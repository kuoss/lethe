package parser_test

import (
	"testing"

	"github.com/kuoss/lethe/letheql/parser"
	"github.com/kuoss/lethe/testutil"
	commonModel "github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {

	testCases := map[string]struct {
		input     string
		wantError string
		want      parser.Expr
	}{
		// BinaryExpr - single FilterOperator
		testutil.TC(): {
			`pod|="hello"`,
			"",
			&parser.BinaryExpr{
				Op: parser.PIPE_EQL,
				LHS: &parser.VectorSelector{
					Name: "pod", LabelMatchers: []*labels.Matcher{parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: parser.PositionRange{Start: 0, End: 3}},
				RHS: &parser.StringLiteral{
					Val:      "hello",
					PosRange: parser.PositionRange{Start: 5, End: 12}},
				ReturnBool: false},
		},
		testutil.TC(): {
			`pod|~"hel.*"`,
			"",
			&parser.BinaryExpr{
				Op: parser.PIPE_REGEX,
				LHS: &parser.VectorSelector{
					Name: "pod", LabelMatchers: []*labels.Matcher{parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: parser.PositionRange{Start: 0, End: 3}},
				RHS: &parser.StringLiteral{
					Val:      "hel.*",
					PosRange: parser.PositionRange{Start: 5, End: 12}},
				ReturnBool: false},
		},
		testutil.TC(): {
			`pod!="hello"`,
			"",
			&parser.BinaryExpr{
				Op: parser.NEQ,
				LHS: &parser.VectorSelector{
					Name: "pod", LabelMatchers: []*labels.Matcher{parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: parser.PositionRange{Start: 0, End: 3}},
				RHS: &parser.StringLiteral{
					Val:      "hello",
					PosRange: parser.PositionRange{Start: 5, End: 12}},
				ReturnBool: false},
		},
		testutil.TC(): {
			`pod!~"hel.*"`,
			"",
			&parser.BinaryExpr{
				Op: parser.NEQ_REGEX,
				LHS: &parser.VectorSelector{
					Name: "pod", LabelMatchers: []*labels.Matcher{parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: parser.PositionRange{Start: 0, End: 3}},
				RHS: &parser.StringLiteral{
					Val:      "hel.*",
					PosRange: parser.PositionRange{Start: 5, End: 12}},
				ReturnBool: false},
		},
		testutil.TC(): {
			`pod |= "hello"`,
			"",
			&parser.BinaryExpr{
				Op: parser.PIPE_EQL,
				LHS: &parser.VectorSelector{
					Name: "pod", LabelMatchers: []*labels.Matcher{parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: parser.PositionRange{Start: 0, End: 3}},
				RHS: &parser.StringLiteral{
					Val:      "hello",
					PosRange: parser.PositionRange{Start: 7, End: 14}},
				ReturnBool: false},
		},
		testutil.TC(): {
			`pod{} |= "hello"`,
			"",
			&parser.BinaryExpr{
				Op: parser.PIPE_EQL,
				LHS: &parser.VectorSelector{
					Name: "pod",
					LabelMatchers: []*labels.Matcher{
						parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: parser.PositionRange{Start: 0, End: 5}},
				RHS: &parser.StringLiteral{
					Val:      "hello",
					PosRange: parser.PositionRange{Start: 9, End: 16}},
				ReturnBool: false},
		},
		testutil.TC(): {
			`pod{} |~ "hello.*"`,
			"",
			&parser.BinaryExpr{
				Op: parser.PIPE_REGEX,
				LHS: &parser.VectorSelector{
					Name: "pod", LabelMatchers: []*labels.Matcher{parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: parser.PositionRange{Start: 0, End: 5}},
				RHS: &parser.StringLiteral{
					Val:      "hello.*",
					PosRange: parser.PositionRange{Start: 9, End: 18}},
				ReturnBool: false},
		},
		testutil.TC(): {
			`pod{namespace="namespace01"} |= "hello"`,
			"",
			&parser.BinaryExpr{
				Op: parser.PIPE_EQL,
				LHS: &parser.VectorSelector{
					Name: "pod", LabelMatchers: []*labels.Matcher{
						parser.MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
						parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: parser.PositionRange{Start: 0, End: 28}},
				RHS: &parser.StringLiteral{
					Val:      "hello",
					PosRange: parser.PositionRange{Start: 32, End: 39}},
				ReturnBool: false},
		},
		testutil.TC(): {
			`pod{namespace="namespace01"} |~ "hel.*"`,
			"",
			&parser.BinaryExpr{
				Op: parser.PIPE_REGEX,
				LHS: &parser.VectorSelector{
					Name: "pod", LabelMatchers: []*labels.Matcher{
						parser.MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
						parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: parser.PositionRange{Start: 0, End: 28}},
				RHS: &parser.StringLiteral{
					Val:      "hel.*",
					PosRange: parser.PositionRange{Start: 32, End: 39}},
				ReturnBool: false},
		},
		// BinaryExpr - multi FilterOperator (nested)
		testutil.TC(): {
			`pod|="hello"!="world"`,
			"",
			&parser.BinaryExpr{
				Op: parser.NEQ,
				LHS: &parser.BinaryExpr{
					Op: parser.PIPE_EQL,
					LHS: &parser.VectorSelector{
						Name: "pod", LabelMatchers: []*labels.Matcher{parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
						PosRange: parser.PositionRange{Start: 0, End: 3}},
					RHS: &parser.StringLiteral{
						Val:      "hello",
						PosRange: parser.PositionRange{Start: 5, End: 12}},
					ReturnBool: false},
				RHS: &parser.StringLiteral{
					Val:      "world",
					PosRange: parser.PositionRange{Start: 14, End: 21}},
				ReturnBool: false},
		},
		testutil.TC(): {
			`pod|~"hel.*"|="world"`,
			"",
			&parser.BinaryExpr{
				Op: parser.PIPE_EQL,
				LHS: &parser.BinaryExpr{
					Op: parser.PIPE_REGEX,
					LHS: &parser.VectorSelector{
						Name: "pod", LabelMatchers: []*labels.Matcher{parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
						PosRange: parser.PositionRange{Start: 0, End: 3}},
					RHS: &parser.StringLiteral{
						Val:      "hel.*",
						PosRange: parser.PositionRange{Start: 5, End: 12}},
					ReturnBool: false},
				RHS: &parser.StringLiteral{
					Val:      "world",
					PosRange: parser.PositionRange{Start: 14, End: 21}},
				ReturnBool: false},
		},
		testutil.TC(): {
			`pod|~"hel.*"!~"wor.*"`,
			"",
			&parser.BinaryExpr{
				Op: parser.NEQ_REGEX,
				LHS: &parser.BinaryExpr{
					Op: parser.PIPE_REGEX,
					LHS: &parser.VectorSelector{
						Name: "pod", LabelMatchers: []*labels.Matcher{parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
						PosRange: parser.PositionRange{Start: 0, End: 3}},
					RHS: &parser.StringLiteral{
						Val:      "hel.*",
						PosRange: parser.PositionRange{Start: 5, End: 12}},
					ReturnBool: false},
				RHS: &parser.StringLiteral{
					Val:      "wor.*",
					PosRange: parser.PositionRange{Start: 14, End: 21}},
				ReturnBool: false},
		},

		// NumberLiteral
		testutil.TC(): {
			`42`,
			"",
			&parser.NumberLiteral{
				Val:      42,
				PosRange: parser.PositionRange{Start: 0, End: 2}},
		},

		testutil.TC(): {
			`"hello"`,
			"",
			&parser.StringLiteral{
				Val:      "hello",
				PosRange: parser.PositionRange{Start: 0, End: 7}},
		},

		// VectorSelector
		testutil.TC(): {
			`pod`,
			"",
			&parser.VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: parser.PositionRange{Start: 0, End: 3}},
		},
		testutil.TC(): {
			`pod{}`,
			"",
			&parser.VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: parser.PositionRange{Start: 0, End: 5}},
		},
		testutil.TC(): {
			`pod{namespace="namespace01"}`,
			"",
			&parser.VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					parser.MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
					parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: parser.PositionRange{Start: 0, End: 28}},
		},
		testutil.TC(): {
			`pod{namespace="not-exists"}`,
			"",
			&parser.VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					parser.MustLabelMatcher(labels.MatchEqual, "namespace", "not-exists"),
					parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: parser.PositionRange{Start: 0, End: 27}},
		},
		testutil.TC(): {
			`pod{namespace="namespace01",pod="nginx"}`,
			"",
			&parser.VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					parser.MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
					parser.MustLabelMatcher(labels.MatchEqual, "pod", "nginx"),
					parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: parser.PositionRange{Start: 0, End: 40}},
		},
		testutil.TC(): {
			`pod{namespace="namespace01",pod="nginx-*"}`,
			"",
			&parser.VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					parser.MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
					parser.MustLabelMatcher(labels.MatchEqual, "pod", "nginx-*"),
					parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: parser.PositionRange{Start: 0, End: 42}},
		},
		testutil.TC(): {
			`pod{namespace="namespace01",container="nginx"}`,
			"",
			&parser.VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					parser.MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
					parser.MustLabelMatcher(labels.MatchEqual, "container", "nginx"),
					parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: parser.PositionRange{Start: 0, End: 46}},
		},
		testutil.TC(): {
			`pod{namespace="namespace*",container="nginx"}`,
			"",
			&parser.VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					parser.MustLabelMatcher(labels.MatchEqual, "namespace", "namespace*"),
					parser.MustLabelMatcher(labels.MatchEqual, "container", "nginx"),
					parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: parser.PositionRange{Start: 0, End: 45}},
		},

		// MatrixSelector
		testutil.TC(): {
			`pod{namespace="namespace01",pod="nginx-*"}[3m]`,
			"",
			&parser.MatrixSelector{
				VectorSelector: &parser.VectorSelector{
					Name: "pod",
					LabelMatchers: []*labels.Matcher{
						parser.MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
						parser.MustLabelMatcher(labels.MatchEqual, "pod", "nginx-*"),
						parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
					PosRange: parser.PositionRange{Start: 0, End: 42}},
				Range: 180000000000, EndPos: 46},
		},

		// Call
		testutil.TC(): {
			`count_over_time(pod{namespace="namespace01",pod="nginx-*"}[3m])`,
			"",
			&parser.Call{
				Func: &parser.Function{
					Name:       "count_over_time",
					ArgTypes:   []parser.ValueType{parser.ValueTypeMatrix},
					ReturnType: parser.ValueTypeVector},
				Args: parser.Expressions{
					&parser.MatrixSelector{
						VectorSelector: &parser.VectorSelector{
							Name: "pod",
							LabelMatchers: []*labels.Matcher{
								parser.MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
								parser.MustLabelMatcher(labels.MatchEqual, "pod", "nginx-*"),
								parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
							PosRange: parser.PositionRange{Start: 16, End: 58}},
						Range: 180000000000, EndPos: 62}},
				PosRange: parser.PositionRange{Start: 0, End: 63}},
		},
		testutil.TC(): {
			`count_over_time(pod{}[3m])`,
			"",
			&parser.Call{
				Func: &parser.Function{
					Name:       "count_over_time",
					ArgTypes:   []parser.ValueType{parser.ValueTypeMatrix},
					ReturnType: parser.ValueTypeVector},
				Args: parser.Expressions{
					&parser.MatrixSelector{
						VectorSelector: &parser.VectorSelector{
							Name:          "pod",
							LabelMatchers: []*labels.Matcher{parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
							PosRange:      parser.PositionRange{Start: 16, End: 21}},
						Range: 180000000000, EndPos: 25}},
				PosRange: parser.PositionRange{Start: 0, End: 26}},
		},

		// BinaryExpr
		testutil.TC(): {
			`count_over_time(pod{}[3m]) > 10`,
			"",
			&parser.BinaryExpr{
				Op: parser.GTR,
				LHS: &parser.Call{
					Func: &parser.Function{
						Name:       "count_over_time",
						ArgTypes:   []parser.ValueType{parser.ValueTypeMatrix},
						ReturnType: parser.ValueTypeVector},
					Args: parser.Expressions{
						&parser.MatrixSelector{
							VectorSelector: &parser.VectorSelector{
								Name:          "pod",
								LabelMatchers: []*labels.Matcher{parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
								PosRange:      parser.PositionRange{Start: 16, End: 21}},
							Range: 180000000000, EndPos: 25}},
					PosRange: parser.PositionRange{Start: 0, End: 26}},
				RHS: &parser.NumberLiteral{
					Val:      10,
					PosRange: parser.PositionRange{Start: 29, End: 31}}},
		},
		testutil.TC(): {
			`count_over_time(pod{}[3m]) < 10`,
			"",
			&parser.BinaryExpr{
				Op: parser.LSS,
				LHS: &parser.Call{
					Func: &parser.Function{
						Name:       "count_over_time",
						ArgTypes:   []parser.ValueType{parser.ValueTypeMatrix},
						ReturnType: parser.ValueTypeVector},
					Args: parser.Expressions{
						&parser.MatrixSelector{
							VectorSelector: &parser.VectorSelector{
								Name:          "pod",
								LabelMatchers: []*labels.Matcher{parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
								PosRange:      parser.PositionRange{Start: 16, End: 21}},
							Range: 180000000000, EndPos: 25}},
					PosRange: parser.PositionRange{Start: 0, End: 26}},
				RHS: &parser.NumberLiteral{
					Val:      10,
					PosRange: parser.PositionRange{Start: 29, End: 31}}},
		},
		testutil.TC(): {
			`count_over_time(pod{}[3m]) == 21`,
			"",
			&parser.BinaryExpr{
				Op: parser.EQLC,
				LHS: &parser.Call{
					Func: &parser.Function{
						Name:       "count_over_time",
						ArgTypes:   []parser.ValueType{parser.ValueTypeMatrix},
						ReturnType: parser.ValueTypeVector},
					Args: parser.Expressions{
						&parser.MatrixSelector{
							VectorSelector: &parser.VectorSelector{
								Name:          "pod",
								LabelMatchers: []*labels.Matcher{parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
								PosRange:      parser.PositionRange{Start: 16, End: 21}},
							Range: 180000000000, EndPos: 25}},
					PosRange: parser.PositionRange{Start: 0, End: 26}},
				RHS: &parser.NumberLiteral{
					Val:      21,
					PosRange: parser.PositionRange{Start: 30, End: 32}},
				VectorMatching: (*parser.VectorMatching)(nil),
				ReturnBool:     false},
		},

		// ######## ERROR
		testutil.TC(): {
			`pod{namespace="namespace01"} "`,
			"1:30: parse error: unterminated quoted string",
			&parser.VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					parser.MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
					parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: parser.PositionRange{Start: 0, End: 28},
			},
		},
		testutil.TC(): {
			`pod{namespace="namespace01"} hello`,
			"1:30: parse error: unexpected identifier \"hello\"",
			&parser.VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					parser.MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
					parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: parser.PositionRange{Start: 0, End: 28},
			},
		},
		testutil.TC(): {
			`pod{namespace="namespace01"} "hello`,
			"1:30: parse error: unterminated quoted string",
			&parser.VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					parser.MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
					parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: parser.PositionRange{Start: 0, End: 28},
			},
		},
		testutil.TC(): {
			`pod{namespace="namespace01"} "hello"`,
			"1:30: parse error: unexpected string \"\\\"hello\\\"\"",
			&parser.VectorSelector{
				Name: "pod",
				LabelMatchers: []*labels.Matcher{
					parser.MustLabelMatcher(labels.MatchEqual, "namespace", "namespace01"),
					parser.MustLabelMatcher(labels.MatchEqual, commonModel.MetricNameLabel, "pod")},
				PosRange: parser.PositionRange{Start: 0, End: 28},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name+" "+tc.input, func(t *testing.T) {
			expr, err := parser.ParseExpr(tc.input)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			assert.Equal(t, tc.want, expr)
		})
	}
}
