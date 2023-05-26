package match

import (
	"fmt"

	"github.com/kuoss/lethe/letheql/model"
	"github.com/prometheus/prometheus/model/labels"
)

func GetTargetMatcher(sel *model.LogSelector) (*labels.Matcher, model.Warnings, error) {
	matcher, err := getTargetLabelMatcher(sel)
	if err != nil {
		return nil, nil, err
	}
	// check matcher type
	switch matcher.Type {
	case labels.MatchNotEqual, labels.MatchRegexp, labels.MatchNotRegexp:
		// for now, selecting multiple targets is discouraged, because cross-target logs are not sorted by time
		return matcher, model.Warnings{fmt.Errorf("warnMultiTargets: use operator '=' for selecting target")}, nil
	case labels.MatchEqual:
		return matcher, nil, nil
	}
	return nil, nil, fmt.Errorf("not supported matcher type: %s", matcher.Type.String())
}

func getTargetLabelMatcher(sel *model.LogSelector) (*labels.Matcher, error) {
	var labelName string
	switch sel.Name {
	case "node":
		labelName = "node"
	case "pod":
		labelName = "namespace"
	default:
		return nil, fmt.Errorf("getTargetMatcher: unknown logType: %s", sel.Name)
	}
	var matcher *labels.Matcher
	for _, m := range sel.LabelMatchers {
		if m.Name == labelName {
			matcher = m
		}
	}
	if matcher == nil {
		return nil, fmt.Errorf("not found label '%s' for logType '%s'", labelName, sel.Name)
	}
	return matcher, nil
}
