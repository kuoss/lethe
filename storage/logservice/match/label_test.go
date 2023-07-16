package match

import (
	"github.com/stretchr/testify/assert"
	"regexp"

	"testing"

	"github.com/kuoss/lethe/letheql/model"
	"github.com/prometheus/prometheus/model/labels"
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
	}

	funcs, err := getLabelMatchFuncs(nodeSel)
	if err != nil {
		f.Fail()
	}

	verfityFunc := funcs[0]

	testcases := []string{"label_value", "non_label_value"}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, label string) {
		isMatch := verfityFunc(label)
		if label == "lable_value" {
			assert.True(t, isMatch)
		} else {
			assert.False(t, isMatch)
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
	}

	funcs, err := getLabelMatchFuncs(nodeSel)
	if err != nil {
		f.Fail()
	}

	verfityFunc := funcs[0]

	testcases := []string{"label_value", "non_label_value"}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, label string) {
		isMatch := verfityFunc(label)
		if label == "lable_value" {
			assert.False(t, isMatch)
		} else {
			assert.True(t, isMatch)
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
	}

	funcs, err := getLabelMatchFuncs(nodeSel)
	if err != nil {
		f.Fail()
	}

	verfityFunc := funcs[0]

	testcases := []string{"foobar", "foobaar", "regardless"}
	for _, tc := range testcases {
		f.Add(tc)
	}
	re := regexp.MustCompile("^.+bar$")
	f.Fuzz(func(t *testing.T, label string) {
		got := verfityFunc(label)
		want := re.MatchString(label)
		assert.Equal(t, got, want)
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
	}

	funcs, err := getLabelMatchFuncs(nodeSel)
	if err != nil {
		f.Fail()
	}

	verfityFunc := funcs[0]

	testcases := []string{"foobar", "foobaar", "regardless"}
	for _, tc := range testcases {
		f.Add(tc)
	}
	re := regexp.MustCompile("^.+bar$")
	f.Fuzz(func(t *testing.T, label string) {
		got := verfityFunc(label)
		want := !re.MatchString(label)
		assert.Equal(t, got, want)
	})
}
