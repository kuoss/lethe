package clock

import (
	"time"
)

var (
	playgroundMode = false
	playgroundTime = time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)
)

func Now() time.Time {
	if playgroundMode {
		return playgroundTime
	}
	return time.Now()
}

func SetPlaygroundMode(newPlaygroundMode bool) {
	playgroundMode = newPlaygroundMode
}
