package logs

import (
	"github.com/kuoss/lethe/testutil"
)

var rotator *Rotator

func init() {
	testutil.Init()
	rotator = NewRotator()
}
