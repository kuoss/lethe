package match

import (
	"github.com/kuoss/lethe/letheql/model"
	"github.com/kuoss/lethe/letheql/parser"
	"github.com/prometheus/prometheus/model/labels"
	"strings"
	"testing"
)

func FuzzNodeLabelMatchEqual(f *testing.F) {
	nodeSel := &model.LogSelector{
		Name: "node",
		LabelMatchers: []*labels.Matcher{
			&labels.Matcher{
				Type:  labels.MatchEqual,
				Name:  "process",
				Value: "label_value",
			},
		},
		LineMatchers: []*model.LineMatcher{
			&model.LineMatcher{
				Op:    parser.PIPE_REGEX,
				Value: "line_value",
			},
		},
		TimeRange: model.TimeRange{},
	}

	funcs, err := getLabelMatchFuncs(nodeSel)
	if err != nil {
		f.Fail()
	}

	var verfityFunc MatchFunc
	verfityFunc = funcs[0]

	testcases := []string{"label_value", "non_label_value"}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, label string) {
		isMatch := verfityFunc(label)
		if label == "label_value" && isMatch == true {
			// as we expected
		}
		if label != "label_value" && isMatch == true {
			t.Errorf("%q label is matched with ", label)
		}
	})
}

func FuzzNodeLabelMatchNotEqual(f *testing.F) {
	nodeSel := &model.LogSelector{
		Name: "node",
		LabelMatchers: []*labels.Matcher{
			&labels.Matcher{
				Type:  labels.MatchNotEqual,
				Name:  "process",
				Value: "label_value",
			},
		},
		LineMatchers: []*model.LineMatcher{
			&model.LineMatcher{
				Op:    parser.PIPE_REGEX,
				Value: "line_value",
			},
		},
		TimeRange: model.TimeRange{},
	}

	funcs, err := getLabelMatchFuncs(nodeSel)
	if err != nil {
		f.Fail()
	}

	var verfityFunc MatchFunc
	verfityFunc = funcs[0]

	testcases := []string{"label_value", "non_label_value"}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, label string) {
		isMatch := verfityFunc(label)
		if label != "label_value" && isMatch == true {
			// as we expected
		}
		if label != "label_value" && isMatch == false {
			t.Errorf("label is not matched with %q", label)
		}
	})
}

func FuzzNodeLabelMatchRegexOnlyContains(f *testing.F) {
	nodeSel := &model.LogSelector{
		Name: "node",
		LabelMatchers: []*labels.Matcher{
			&labels.Matcher{
				Type:  labels.MatchRegexp,
				Name:  "process",
				Value: ".+bar",
			},
		},
		LineMatchers: []*model.LineMatcher{
			&model.LineMatcher{
				Op:    parser.PIPE_REGEX,
				Value: "line_value",
			},
		},
		TimeRange: model.TimeRange{},
	}

	funcs, err := getLabelMatchFuncs(nodeSel)
	if err != nil {
		f.Fail()
	}

	var verfityFunc MatchFunc
	verfityFunc = funcs[0]

	testcases := []string{"foobar", "foobaar", "regardless"}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, label string) {
		isMatch := verfityFunc(label)
		if !strings.Contains(label, "bar") && isMatch == true {
			t.Errorf("label is matched with %q unexpectedly", label)
		}
	})
}

func FuzzNodeLabelMatchRegexNotEqualOnlyContains(f *testing.F) {
	nodeSel := &model.LogSelector{
		Name: "node",
		LabelMatchers: []*labels.Matcher{
			&labels.Matcher{
				Type:  labels.MatchNotRegexp,
				Name:  "process",
				Value: ".+bar",
			},
		},
		LineMatchers: []*model.LineMatcher{
			&model.LineMatcher{
				Op:    parser.PIPE_REGEX,
				Value: "line_value",
			},
		},
		TimeRange: model.TimeRange{},
	}

	funcs, err := getLabelMatchFuncs(nodeSel)
	if err != nil {
		f.Fail()
	}

	var verfityFunc MatchFunc
	verfityFunc = funcs[0]

	testcases := []string{"foobar", "foobaar", "regardless"}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, label string) {
		isMatch := verfityFunc(label)
		if !strings.Contains(label, "bar") && isMatch == false {
			t.Errorf("label is matched with %q unexpectedly", label)
		}
	})
}
