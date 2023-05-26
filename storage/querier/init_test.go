package querier

import (
	"github.com/kuoss/lethe/clock"
	"github.com/kuoss/lethe/util/testutil"
)

func init() {
	testutil.ChdirRoot()
	testutil.ResetLogData()
	clock.SetPlaygroundMode(true)
}
