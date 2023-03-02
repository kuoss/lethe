package logs

import (
	"testing"

	"github.com/kuoss/lethe/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_GetDiskUsedBytes(t *testing.T) {
	testutil.SetTestLogFiles()

	actual, err := rotator.GetDiskUsedBytes(rotator.driver.RootDirectory())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotZero(t, actual)
}
