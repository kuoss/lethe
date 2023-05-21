package match

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kuoss/lethe/letheql/model"
	"github.com/kuoss/lethe/letheql/parser"
)

func getLineMatchFuncs(sel *model.LogSelector) ([]MatchFunc, error) {
	var funcs []MatchFunc
	for _, m := range sel.LineMatchers {
		f, err := getLineMatchFunc(m)
		if err != nil {
			return nil, fmt.Errorf("getLineMatchFunc err: %w", err)
		}
		funcs = append(funcs, f)
	}
	return funcs, nil
}

func getLineMatchFunc(m *model.LineMatcher) (MatchFunc, error) {
	switch m.Op {
	case parser.PIPE_EQL: // |=
		return func(s string) bool { return strings.Contains(s, m.Value) }, nil
	case parser.NEQ: // !=
		return func(s string) bool { return !strings.Contains(s, m.Value) }, nil
	case parser.PIPE_REGEX: // |~
		re, err := regexp.Compile(m.Value)
		if err != nil {
			return nil, err
		}
		return func(s string) bool { return re.MatchString(s) }, nil
	case parser.NEQ_REGEX: // !~
		re, err := regexp.Compile(m.Value)
		if err != nil {
			return nil, err
		}
		return func(s string) bool { return !re.MatchString(s) }, nil
	}
	return nil, fmt.Errorf("unknown match op: %s", m.Op)
}
