package logs

import (
	"runtime"
	"testing"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_RoutineDelete_50k(t *testing.T) {
	testutil.SetTestLogFiles()

	config.GetConfig().Set("retention.time", "100m")
	config.GetConfig().Set("retention.size", "50k")
	rotator.RunOnce()

	if runtime.GOOS == "linux" {
		assert.NoFileExists(t, "./tmp/log/node/node01/2009-11-10_21.log")
		assert.NoFileExists(t, "./tmp/log/node/node01/2009-11-10_22.log")

		assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-01_00.log")
		assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-10_21.log")

		assert.NoFileExists(t, "./tmp/log/pod/namespace01/2000-01-01_00.log")
		// assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_21.log")
		assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_22.log")

		assert.NoFileExists(t, "./tmp/log/pod/namespace02/0000-00-00_00.log")
		assert.NoFileExists(t, "./tmp/log/pod/namespace02/2009-11-10_22.log")
	} else {
		assert.FileExists(t, "./tmp/log/node/node01/2009-11-10_21.log")
		assert.FileExists(t, "./tmp/log/node/node01/2009-11-10_22.log")

		assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-01_00.log")
		assert.FileExists(t, "./tmp/log/node/node02/2009-11-10_21.log")

		assert.NoFileExists(t, "./tmp/log/pod/namespace01/2000-01-01_00.log")
		assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_21.log")
		assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_22.log")

		assert.NoFileExists(t, "./tmp/log/pod/namespace02/0000-00-00_00.log")
		assert.FileExists(t, "./tmp/log/pod/namespace02/2009-11-10_22.log")
	}
}

func Test_RoutineDelete_15k(t *testing.T) {
	testutil.SetTestLogFiles()

	config.GetConfig().Set("retention.time", "100m")
	config.GetConfig().Set("retention.size", "15k")
	rotator.RunOnce()

	if runtime.GOOS == "linux" {
		assert.NoFileExists(t, "./tmp/log/node/node01/2009-11-10_21.log")
		assert.NoFileExists(t, "./tmp/log/node/node01/2009-11-10_22.log")

		assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-01_00.log")
		assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-10_21.log")

		assert.NoFileExists(t, "./tmp/log/pod/namespace01/2000-01-01_00.log")
		// assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_21.log")
		assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_22.log")

		assert.NoFileExists(t, "./tmp/log/pod/namespace02/0000-00-00_00.log")
		assert.NoFileExists(t, "./tmp/log/pod/namespace02/2009-11-10_22.log")
	} else {
		assert.NoFileExists(t, "./tmp/log/node/node01/2009-11-10_21.log")
		// assert.NoFileExists(t, "./tmp/log/node/node01/2009-11-10_22.log")

		assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-01_00.log")
		assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-10_21.log")

		assert.NoFileExists(t, "./tmp/log/pod/namespace01/2000-01-01_00.log")
		// assert.NoFileExists(t, "./tmp/log/pod/namespace01/2009-11-10_21.log")
		// assert.NoFileExists(t, "./tmp/log/pod/namespace01/2009-11-10_22.log")

		assert.NoFileExists(t, "./tmp/log/pod/namespace02/0000-00-00_00.log")
		// assert.FileExists(t, "./tmp/log/pod/namespace02/2009-11-10_22.log")
	}
}

func Test_RoutineDelete_1d(t *testing.T) {
	testutil.SetTestLogFiles()

	config.GetConfig().Set("retention.time", "1d")
	config.GetConfig().Set("retention.size", "100g")
	rotator.RunOnce()

	assert.FileExists(t, "./tmp/log/node/node01/2009-11-10_21.log")
	assert.FileExists(t, "./tmp/log/node/node01/2009-11-10_22.log")

	assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-01_00.log")
	assert.FileExists(t, "./tmp/log/node/node02/2009-11-10_21.log")

	assert.NoFileExists(t, "./tmp/log/pod/namespace01/2000-01-01_00.log")
	assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_21.log")
	assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_22.log")

	assert.NoFileExists(t, "./tmp/log/pod/namespace02/0000-00-00_00.log")
	assert.NoFileExists(t, "./tmp/log/pod/namespace02/2009-11-10_21.log")
}
