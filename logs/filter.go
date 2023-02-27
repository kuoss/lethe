package logs

import (
	"errors"
	"regexp"
	"strings"
)

const (
	include      = "|="
	exclude      = "!="
	includeRegex = "|~"
	excludeRegex = "!~"
)

func IsFilterExist(query string) (ok bool, filter string, err error) {
	// move to heap Filterlist
	filterList := []string{include, exclude, includeRegex, excludeRegex}
	var filtersInQuery []string
	for _, v := range filterList {
		if strings.Contains(query, v) {
			filtersInQuery = append(filtersInQuery, v)
		}
	}
	switch len(filtersInQuery) {
	case 0:
		return false, "", nil
	case 1:
		return true, filtersInQuery[0], nil
	}
	return false, "", errors.New("filter must be only one or no filters")
}

func FilterFromQuery(query string) (Filter, error) {
	_, filterType, _ := IsFilterExist(query)
	parts := strings.Split(query, filterType)
	switch filterType {
	case include:
		return &includeFilter{keyword: strings.TrimSpace(parts[1])}, nil
	case exclude:
		return &excludeFilter{keyword: strings.TrimSpace(parts[1])}, nil
	case includeRegex:
		f := &includeRegexFilter{keyword: strings.TrimSpace(parts[1])}
		if f.isRegexFilter() {
			return f, nil
		}
		return &includeRegexFilter{}, errors.New("wrong regex expression")
	case excludeRegex:
		f := &excludeRegexFilter{keyword: strings.TrimSpace(parts[1])}
		if f.isRegexFilter() {
			return f, nil
		}
		return &excludeRegexFilter{}, errors.New("wrong regex expression")
	}
	return nil, errors.New("there is wrong filter in query")
}

type Filter interface {
	match(string) bool
}

type RegexFilter interface {
	isRegex() bool
}

// for just test build
type TempExportFilter struct{}

func (f TempExportFilter) match(line string) bool {
	return true
}

type includeFilter struct {
	keyword string
}

func (f *includeFilter) match(line string) bool {
	return strings.Contains(line, f.keyword)
}

type excludeFilter struct {
	keyword string
}

func (f *excludeFilter) match(line string) bool {
	return !strings.Contains(line, f.keyword)
}

type includeRegexFilter struct {
	regex   *regexp.Regexp
	keyword string
}

func (f *includeRegexFilter) match(line string) bool {
	return f.regex.MatchString(line)
}

func (f *includeRegexFilter) isRegexFilter() bool {
	compile, err := regexp.Compile(f.keyword)
	f.regex = compile
	if err != nil {
		return false
	}
	if err != nil {
		return false
	}
	return true
}

type excludeRegexFilter struct {
	regex   *regexp.Regexp
	keyword string
}

func (f *excludeRegexFilter) match(line string) bool {
	return !f.regex.MatchString(line)
}

func (f *excludeRegexFilter) isRegexFilter() bool {
	compile, err := regexp.Compile(f.keyword)
	f.regex = compile
	if err != nil {
		return false
	}
	if err != nil {
		return false
	}
	return true
}
