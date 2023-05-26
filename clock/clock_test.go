package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	playgroundTime_test = time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)
)

func TestNow(t *testing.T) {

	justBefore := time.Now()

	// normal mode
	assert.NotEqual(t, playgroundTime_test, Now())
	assert.Greater(t, Now(), justBefore)

	// playground mode
	SetPlaygroundMode(true)
	assert.Equal(t, playgroundTime_test, Now())

	// normal mode
	SetPlaygroundMode(false)
	assert.NotEqual(t, playgroundTime_test, Now())
	assert.Greater(t, Now(), justBefore)
}
