package logs

import (
	rotatorTest "github.com/kuoss/lethe/logs/rotator"
	"github.com/kuoss/lethe/testutil"
)

var rotator *rotatorTest.Rotator

func init() {
	testutil.Init()
	testutil.SetTestLogFiles()
	rotator = rotatorTest.NewRotator()
}
