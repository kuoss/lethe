package logs

import (
	"github.com/kuoss/lethe/testutil"
)

var rotator *Rotator

func init() {
	testutil.Init()
	testutil.SetTestLogFiles()
	rotator = NewRotator()
}
