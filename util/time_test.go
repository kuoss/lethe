package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_GetDurationFromAge(t *testing.T) {
	var d, d2 time.Duration

	d, _ = GetDurationFromAge("2m")
	assert.Equal(t, time.Duration(120000000000), d)

	d, _ = GetDurationFromAge("2h")
	assert.Equal(t, time.Duration(7200000000000), d)

	d, _ = GetDurationFromAge("2d")
	d2, _ = time.ParseDuration("48h")
	assert.Equal(t, d2, d)

	d, _ = GetDurationFromAge("100d")
	d2, _ = time.ParseDuration("2400h")
	assert.Equal(t, d2, d)
}

func Test_FloatStringToTime(t *testing.T) {
	want := time.Date(2015, time.July, 1, 20, 10, 51, 780999898, time.UTC)
	got := FloatStringToTime("1435781451.781")
	assert.True(t, want.Equal(got))
}
