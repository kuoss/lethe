package rotator

import (
	"testing"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_DeleteByAge_10d(t *testing.T) {
	testutil.Init()
	testutil.SetTestLogFiles()

	config.GetConfig().Set("retention.time", "20d")
	NewRotator().DeleteByAge()

	assert.FileExists(t, "./tmp/log/node/node01/2009-11-10_21.log")
	assert.FileExists(t, "./tmp/log/node/node01/2009-11-10_22.log")

	assert.FileExists(t, "./tmp/log/node/node02/2009-11-01_00.log")
	assert.FileExists(t, "./tmp/log/node/node02/2009-11-10_21.log")

	assert.NoFileExists(t, "./tmp/log/pod/namespace01/2000-01-01_00.log")
	assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_21.log")
	assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_22.log")

	assert.NoFileExists(t, "./tmp/log/pod/namespace02/0000-00-00_00.log")
	assert.NoFileExists(t, "./tmp/log/pod/namespace02/2009-11-10_21.log")
}

func Test_DeleteByAge_1d(t *testing.T) {
	testutil.Init()
	testutil.SetTestLogFiles()

	config.GetConfig().Set("retention.time", "2d")
	config.GetConfig().Set("retention.size", "100m")
	NewRotator().DeleteByAge()

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

func Test_DeleteByAge_1h(t *testing.T) {
	testutil.SetTestLogFiles()

	config.GetConfig().Set("retention.time", "1h")
	config.GetConfig().Set("retention.size", "100m")
	NewRotator().DeleteByAge()

	assert.NoFileExists(t, "./tmp/log/node/node01/2009-11-10_21.log")
	assert.FileExists(t, "./tmp/log/node/node01/2009-11-10_22.log")

	assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-01_00.log")
	assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-10_21.log")

	assert.NoFileExists(t, "./tmp/log/pod/namespace01/2000-01-01_00.log")
	// assert.NoFileExists(t, "./tmp/log/pod/namespace01/2009-11-10_21.log")
	assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_22.log")

	assert.NoFileExists(t, "./tmp/log/pod/namespace02/0000-00-00_00.log")
	assert.NoFileExists(t, "./tmp/log/pod/namespace02/2009-11-10_21.log")
}

func Test_DeleteBySize_1m(t *testing.T) {
	testutil.SetTestLogFiles()

	config.GetConfig().Set("retention.size", "1m")
	NewRotator().DeleteBySize()

	assert.FileExists(t, "./tmp/log/node/node01/2009-11-10_21.log")
	assert.FileExists(t, "./tmp/log/node/node01/2009-11-10_22.log")

	assert.FileExists(t, "./tmp/log/node/node02/2009-11-01_00.log")
	assert.FileExists(t, "./tmp/log/node/node02/2009-11-10_21.log")

	assert.FileExists(t, "./tmp/log/pod/namespace01/2000-01-01_00.log")
	assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_21.log")
	assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_22.log")

	assert.FileExists(t, "./tmp/log/pod/namespace02/0000-00-00_00.log")
	assert.NoFileExists(t, "./tmp/log/pod/namespace02/2009-11-10_21.log")
}

func Test_DeleteBySize_3k(t *testing.T) {
	testutil.Init()
	testutil.SetTestLogFiles()

	config.GetConfig().Set("retention.size", "3k")
	NewRotator().DeleteBySize()

	assert.NoFileExists(t, "./tmp/log/node/node01/2009-11-10_21.log")
	assert.FileExists(t, "./tmp/log/node/node01/2009-11-10_22.log")

	assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-01_00.log")
	assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-10_21.log")

	assert.NoFileExists(t, "./tmp/log/pod/namespace01/2000-01-01_00.log")
	assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_21.log")
	assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_22.log")

	assert.NoFileExists(t, "./tmp/log/pod/namespace02/0000-00-00_00.log")
	assert.NoFileExists(t, "./tmp/log/pod/namespace02/2009-11-10_21.log")
}

func Test_DeleteBySize_2k(t *testing.T) {
	testutil.Init()
	testutil.SetTestLogFiles()

	config.GetConfig().Set("retention.size", "2k")
	NewRotator().DeleteBySize()

	assert.NoFileExists(t, "./tmp/log/node/node01/2009-11-10_21.log")
	assert.NoFileExists(t, "./tmp/log/node/node01/2009-11-10_22.log")

	assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-01_00.log")
	assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-10_21.log")

	assert.NoFileExists(t, "./tmp/log/pod/namespace01/2000-01-01_00.log")
	assert.NoFileExists(t, "./tmp/log/pod/namespace01/2009-11-10_21.log")
	assert.NoFileExists(t, "./tmp/log/pod/namespace01/2009-11-10_22.log")

	assert.NoFileExists(t, "./tmp/log/pod/namespace02/0000-00-00_00.log")
	assert.NoFileExists(t, "./tmp/log/pod/namespace02/2009-11-10_21.log")
}