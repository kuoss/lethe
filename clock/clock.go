package clock

import (
	"os"
	"time"
)

func GetNow() time.Time {
	if os.Getenv("TEST_MODE") == "1" {
		return time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)
	}
	return time.Now()
}
