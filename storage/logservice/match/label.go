package match

import (
	"fmt"
	"regexp"

	"github.com/kuoss/lethe/letheql/model"
	"github.com/prometheus/prometheus/model/labels"
)

func getLabelMatchFuncs(sel *model.LogSelector) ([]MatchFunc, error) {
	switch sel.Name {
	case "node":
		return getLabelMatchFuncsDetail(sel, "process")
	case "pod":
		return getLabelMatchFuncsDetail(sel, "pod", "container")
	}
	return nil, fmt.Errorf("unknwon logType: %s", sel.Name)
}

func getLabelMatchFuncsDetail(sel *model.LogSelector, names ...string) ([]MatchFunc, error) {
	var funcs []MatchFunc
	for _, name := range names {
		f, err := getLabelMatchFunc(sel, name)
		if err != nil {
			return nil, fmt.Errorf("getLabelMatchFunc err: %w", err)
		}
		funcs = append(funcs, f)
	}
	return funcs, nil
}

func getLabelMatchFunc(sel *model.LogSelector, name string) (MatchFunc, error) {
	m := getLabelMatcher(sel, name)
	if m == nil {
		return nil, nil // ok (empty)
	}
	switch m.Type {
	case labels.MatchEqual:
		return func(s string) bool { return s == m.Value }, nil
	case labels.MatchNotEqual:
		return func(s string) bool { return s != m.Value }, nil
	case labels.MatchRegexp:
		re, err := regexp.Compile("^(?:" + m.Value + ")$")
		if err != nil {
			return nil, err
		}
		return func(s string) bool { return re.MatchString(s) }, nil
	case labels.MatchNotRegexp:
		re, err := regexp.Compile("^(?:" + m.Value + ")$")
		if err != nil {
			return nil, err
		}
		return func(s string) bool { return !re.MatchString(s) }, nil
	}
	return nil, fmt.Errorf("unknown match type: %s", m.Type)
}

func getLabelMatcher(sel *model.LogSelector, name string) *labels.Matcher {
	for _, matcher := range sel.LabelMatchers {
		if matcher.Name == name {
			return matcher
		}
	}
	return nil // ok (empty)
}
