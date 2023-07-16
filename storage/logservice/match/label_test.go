package match

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kuoss/lethe/letheql/model"
	"github.com/prometheus/prometheus/model/labels"
)

func FuzzNodeLabelMatchEqual(f *testing.F) {

	f.Add("needle", "haystack")

	f.Fuzz(func(t *testing.T, needle, haystack string) {
		nodeSel := &model.LogSelector{
			Name: "node",
			LabelMatchers: []*labels.Matcher{
				&labels.Matcher{
					Type:  labels.MatchEqual,
					Name:  "process",
					Value: needle,
				},
			},
		}

		funcs, err := getLabelMatchFuncs(nodeSel)
		if err != nil {
			t.Fail()
		}

		verfityFunc := funcs[0]

		isMatch := verfityFunc(haystack)

		if needle == haystack {
			assert.True(t, isMatch)
		} else {
			assert.False(t, isMatch)
		}
	})
}

func FuzzNodeLabelMatchNotEqual(f *testing.F) {
	f.Add("needle", "haystack")

	f.Fuzz(func(t *testing.T, needle, haystack string) {
		nodeSel := &model.LogSelector{
			Name: "node",
			LabelMatchers: []*labels.Matcher{
				&labels.Matcher{
					Type:  labels.MatchNotEqual,
					Name:  "process",
					Value: needle,
				},
			},
		}

		funcs, err := getLabelMatchFuncs(nodeSel)
		if err != nil {
			t.Fail()
		}

		verfityFunc := funcs[0]

		isMatch := verfityFunc(haystack)

		if needle != haystack {
			assert.True(t, isMatch)
		} else {
			assert.False(t, isMatch)
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
