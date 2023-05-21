package querier

import (
	"github.com/prometheus/prometheus/model/labels"
	promstorage "github.com/prometheus/prometheus/storage"
)

// A LetheQueryable is used for testing purposes so that a Lethe Querier can be used.
type LetheQueryable struct {
	LetheQuerier promstorage.Querier
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

// temporary dummy storage for working
type LetheStorage struct {
	// exemplarStorage tsdb.ExemplarStorage
}

func (s LetheStorage) Close() error {
	return nil
}

// func (s LetheStorage) ExemplarAppender() promstorage.ExemplarAppender {
// 	return s
// }

// func (s LetheStorage) ExemplarQueryable() promstorage.ExemplarQueryable {
// 	return s.exemplarStorage
// }

// func (s LetheStorage) AppendExemplar(ref promstorage.SeriesRef, l labels.Labels, e exemplar.Exemplar) (promstorage.SeriesRef, error) {
// 	return ref, s.exemplarStorage.AddExemplar(l, e)
// }
