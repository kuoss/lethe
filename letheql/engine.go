package letheql

import (
	"github.com/VictoriaMetrics/metricsql"
	"github.com/kuoss/lethe/logs/filter"
	"github.com/pkg/errors"
	"strings"
)

type Engine struct {
}

// NewQuery return Query that including filter, keyword, engine
func (e *Engine) newQuery(queryString string) (*query, error) {

	ok, filterType, err := filter.IsFilterExist(queryString)
	if err != nil {
		return nil, err
	}

	var parsableQuery, keyword string
	var f filter.Filter

	if ok {
		parts := strings.Split(queryString, filterType)
		parsableQuery = strings.TrimSpace(parts[0])
		keyword = strings.TrimSpace(parts[1])
		filterFromQuery, err := filter.FilterFromQuery(queryString)
		if err != nil {
			return nil, err
		}
		f = filterFromQuery
	} else {
		parsableQuery = queryString
		keyword = ""
	}

	if len(queryString) < 1 {
		return nil, errors.New("empty queryString")
	}
	return &query{
		q:       parsableQuery,
		filter:  f,
		keyword: keyword,
		engine:  e,
	}, nil
}

func (e *Engine) parseQuery(q *query) {
	_, err := metricsql.Parse(q.q)
	if err != nil {
		return
	}
}
func (e *Engine) exec(q *query) {
	q.Exec()
}
