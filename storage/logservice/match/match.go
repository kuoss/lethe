package match

import (
	"fmt"

	"github.com/kuoss/lethe/letheql/model"
)

type MatchFunc func(s string) bool
type MatchFuncSet struct {
	LabelMatchFuncs []MatchFunc
	LineMatchFuncs  []MatchFunc
}

func GetMatchFuncSet(sel *model.LogSelector) (*MatchFuncSet, error) {
	labelMatchFuncs, err := getLabelMatchFuncs(sel)
	if err != nil {
		return nil, fmt.Errorf("getLabelMatchFuncs err: %w", err)
	}
	lineMatchFuncs, err := getLineMatchFuncs(sel)
	if err != nil {
		return nil, fmt.Errorf("getLineMatchFuncs err: %w", err)
	}
	return &MatchFuncSet{
		LabelMatchFuncs: labelMatchFuncs,
		LineMatchFuncs:  lineMatchFuncs,
	}, nil
}
