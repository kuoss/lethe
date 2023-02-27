package logs

import (
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

	assert.NoFileExists(t, "./tmp/log/node/node01/2009-11-10_21.log")
	assert.NoFileExists(t, "./tmp/log/node/node01/2009-11-10_22.log")

	assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-01_00.log")
	assert.NoFileExists(t, "./tmp/log/node/node02/2009-11-10_21.log")

	assert.NoFileExists(t, "./tmp/log/pod/namespace01/2000-01-01_00.log")
	assert.NoFileExists(t, "./tmp/log/pod/namespace01/2009-11-10_21.log")
	// assert.FileExists(t, "./tmp/log/pod/namespace01/2009-11-10_22.log")

	assert.NoFileExists(t, "./tmp/log/pod/namespace02/0000-00-00_00.log")
	assert.NoFileExists(t, "./tmp/log/pod/namespace02/2009-11-10_21.log")
}

func Test_RoutineDelete_1d(t *testing.T) {
	testutil.SetTestLogFiles()

	config.GetConfig().Set("retention.time", "1d")
	config.GetConfig().Set("retention.size", "100g")
	rotator.RunOnce()

	got := rotator.ListFiles()
	want := `[{"FullPath":"tmp/log/node/node01/2009-11-10_21.log","SubPath":"2009-11-10_21.log","LogType":"node","Target":"node01","Name":"2009-11-10_21.log","Extention":".log","Size":1057},{"FullPath":"tmp/log/node/node01/2009-11-10_22.log","SubPath":"2009-11-10_22.log","LogType":"node","Target":"node01","Name":"2009-11-10_22.log","Extention":".log","Size":177},{"FullPath":"tmp/log/node/node02/2009-11-10_21.log","SubPath":"2009-11-10_21.log","LogType":"node","Target":"node02","Name":"2009-11-10_21.log","Extention":".log","Size":1116},{"FullPath":"tmp/log/pod/namespace01/2009-11-10_21.log","SubPath":"2009-11-10_21.log","LogType":"pod","Target":"namespace01","Name":"2009-11-10_21.log","Extention":".log","Size":279},{"FullPath":"tmp/log/pod/namespace01/2009-11-10_22.log","SubPath":"2009-11-10_22.log","LogType":"pod","Target":"namespace01","Name":"2009-11-10_22.log","Extention":".log","Size":1031},{"FullPath":"tmp/log/pod/namespace01/2029-11-10_23.log","SubPath":"2029-11-10_23.log","LogType":"pod","Target":"namespace01","Name":"2029-11-10_23.log","Extention":".log","Size":279},{"FullPath":"tmp/log/pod/namespace02/2009-11-10_22.log","SubPath":"2009-11-10_22.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_22.log","Extention":".log","Size":1125}]`

	testutil.CheckEqualJSON(t, got, want)
	testutil.CheckEqualJSON(t, len(got), `7`)
}
