package main

import (
	"testing"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_task_deleteByAge_10d(t *testing.T) {
	testutil.Init()
	testutil.SetTestLogFiles()

	config.Viper().Set("retention.time", "20d")
	execute("task", "delete-by-age")

	assert.FileExists(t, "./tmp/log/node/node01/2009-11-10_21.log")
	assert.FileExists(t, "./tmp/log/node/node01/2009-11-10_22.log")

	assert.FileExists(t, "./tmp/log/node/node02/2009-11-01_00.log")
	assert.FileExists(t, "./tmp/log/node/node02/2009-11-10_21.log")

	// assert.NoFileExists(t, "./tmp/log/pod/namespace01/2000-01-01_00.log")
	assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_21.log")
	assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_22.log")

	// assert.NoFileExists(t, "./tmp/log/pod/namespace02/0000-00-00_00.log")
	assert.NoFileExists(t, "./tmp/log/pod/namespace02/2009-11-10_21.log")
}

func Test_task_deleteByAge_1d(t *testing.T) {
	testutil.Init()
	testutil.SetTestLogFiles()

	config.Viper().Set("retention.time", "2d")
	execute("task", "delete-by-age")

	assert.FileExists(t, "./tmp/log/node/node01/2009-11-10_21.log")
	assert.FileExists(t, "./tmp/log/node/node01/2009-11-10_22.log")

	// assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-01_00.log")
	assert.FileExists(t, "./tmp/log/node/node02/2009-11-10_21.log")

	// assert.NoFileExists(t, "./tmp/log/pod/namespace01/2000-01-01_00.log")
	assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_21.log")
	assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_22.log")

	// assert.NoFileExists(t, "./tmp/log/pod/namespace02/0000-00-00_00.log")
	assert.NoFileExists(t, "./tmp/log/pod/namespace02/2009-11-10_21.log")
}

func Test_task_deleteByAge_1h(t *testing.T) {
	testutil.Init()
	testutil.SetTestLogFiles()

	config.Viper().Set("retention.time", "1h")
	execute("task", "delete-by-age")

	// assert.NoFileExists(t, "./tmp/log/node/node01/2009-11-10_21.log")
	assert.FileExists(t, "./tmp/log/node/node01/2009-11-10_22.log")

	// assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-01_00.log")
	// assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-10_21.log")

	// assert.NoFileExists(t, "./tmp/log/pod/namespace01/2000-01-01_00.log")
	// assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_21.log")
	assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_22.log")

	// assert.NoFileExists(t, "./tmp/log/pod/namespace02/0000-00-00_00.log")
	assert.NoFileExists(t, "./tmp/log/pod/namespace02/2009-11-10_21.log")
}

func Test_task_deleteBySize_1m(t *testing.T) {
	testutil.Init()
	testutil.SetTestLogFiles()

	config.Viper().Set("retention.size", "1m")
	execute("task", "delete-by-size")

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

func Test_task_deleteBySize_1k(t *testing.T) {
	testutil.Init()
	testutil.SetTestLogFiles()

	config.Viper().Set("retention.size", "1k")
	execute("task", "delete-by-size")

	// assert.NoFileExists(t, "./tmp/log/node/node01/2009-11-10_21.log")
	// assert.NoFileExists(t, "./tmp/log/node/node01/2009-11-10_22.log")

	// assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-01_00.log")
	// assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-10_21.log")

	// assert.NoFileExists(t, "./tmp/log/pod/namespace01/2000-01-01_00.log")
	// assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_21.log")
	// assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_22.log")

	// assert.NoFileExists(t, "./tmp/log/pod/namespace02/0000-00-00_00.log")
	assert.NoFileExists(t, "./tmp/log/pod/namespace02/2009-11-10_21.log")
}
