package querier

import (
	"context"
	"testing"

	"github.com/prometheus/prometheus/model/labels"
	promstorage "github.com/prometheus/prometheus/storage"
	"github.com/stretchr/testify/require"
)

// MockSeriesSet is a mock implementation of the promstorage.SeriesSet interface for testing purposes.
type MockSeriesSet struct{}

func (m *MockSeriesSet) Next() bool                     { return false }
func (m *MockSeriesSet) At() promstorage.Series         { return nil }
func (m *MockSeriesSet) Err() error                     { return nil }
func (m *MockSeriesSet) Warnings() promstorage.Warnings { return nil }

func TestLetheQueryable_Querier(t *testing.T) {
	mockSeriesSet := &MockSeriesSet{}
	mockSelectFunction := func(sortSeries bool, hints *promstorage.SelectHints, matchers ...*labels.Matcher) promstorage.SeriesSet {
		return mockSeriesSet
	}

	letheQuerier := &LetheQuerier{
		SelectLetheFunction: mockSelectFunction,
	}

	letheQueryable := &LetheQueryable{
		LetheQuerier: letheQuerier,
	}

	// Create a new querier using the LetheQueryable
	querier, err := letheQueryable.Querier(context.Background(), 0, 0)
	require.NoError(t, err, "Expected no error from LetheQueryable.Querier")

	// Verify the returned querier is the same as the one we set up
	require.Equal(t, letheQuerier, querier, "Expected the querier to match the LetheQuerier")

	// Test the Select method of LetheQuerier
	result := querier.Select(false, nil)
	require.Equal(t, mockSeriesSet, result, "Expected the SeriesSet returned by Select to match the mockSeriesSet")
}
