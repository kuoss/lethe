package querier

import (
	"context"

	"github.com/prometheus/prometheus/model/labels"
	promstorage "github.com/prometheus/prometheus/storage"
)

// A LetheQueryable is used for testing purposes so that a Lethe Querier can be used.
type LetheQueryable struct {
	LetheQuerier promstorage.Querier
}

func (q *LetheQueryable) Querier(context.Context, int64, int64) (promstorage.Querier, error) {
	return q.LetheQuerier, nil
}

// LetheQuerier is used for test purposes to Lethe the selected series that is returned.
type LetheQuerier struct {
	SelectLetheFunction func(sortSeries bool, hints *promstorage.SelectHints, matchers ...*labels.Matcher) promstorage.SeriesSet
}

func (q *LetheQuerier) LabelValues(string, ...*labels.Matcher) ([]string, promstorage.Warnings, error) {
	return nil, nil, nil
}

func (q *LetheQuerier) LabelNames(...*labels.Matcher) ([]string, promstorage.Warnings, error) {
	return nil, nil, nil
}

func (q *LetheQuerier) Close() error {
	return nil
}

func (q *LetheQuerier) Select(sortSeries bool, hints *promstorage.SelectHints, matchers ...*labels.Matcher) promstorage.SeriesSet {
	return q.SelectLetheFunction(sortSeries, hints, matchers...)
}
