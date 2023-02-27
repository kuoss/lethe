package clock

import (
	"os"
	"testing"
	"time"

	"github.com/kuoss/lethe/testutil"
)

var testTime time.Time = time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)

func init() {
	testutil.Init()
}

func Test_GetNow_TestMode(t *testing.T) {
	if GetNow() != testTime {
		t.Fatalf("Now != testTime in test mode.")
	}
}

func Test_GetNow_ProdMode(t *testing.T) {
	os.Setenv("TEST_MODE", "0")
	if GetNow() == testTime {
		t.Fatalf("Now == testTime in prod mode.")
	}
	os.Setenv("TEST_MODE", "1")
}
