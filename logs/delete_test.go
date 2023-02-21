package logs

import (
	"fmt"
	"github.com/kuoss/lethe/storage/driver/factory"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"

	"github.com/kuoss/lethe/testutil"
)

func Test_GetDiskUsedKB(t *testing.T) {

	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp", "log")})
	rotator := Rotator{driver: d}

	kb, err := rotator.GetDiskUsedKB(rotator.driver.RootDirectory())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("kb=", kb)
}

func Test_DeleteByAge(t *testing.T) {
	testutil.SetTestLogFiles()

	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp2", "log")})
	rotator := Rotator{driver: d}

	rotator.DeleteByAge(false)

	got := rotator.ListFiles()
	//want := `[{"Filepath":"/tmp/log/node/node01/2009-11-10_20.log","LogType":"node","Target":"node01","Name":"2009-11-10_20","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node01/2009-11-10_21.log","LogType":"node","Target":"node01","Name":"2009-11-10_21","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node01/2009-11-10_22.log","LogType":"node","Target":"node01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node01/2009-11-10_23.log","LogType":"node","Target":"node01","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_20.log","LogType":"node","Target":"node02","Name":"2009-11-10_20","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_21.log","LogType":"node","Target":"node02","Name":"2009-11-10_21","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_22.log","LogType":"node","Target":"node02","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_23.log","LogType":"node","Target":"node02","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_20.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_20","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_21.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_21","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_22.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_23.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_20.log","LogType":"pod","Target":"namspace01","Name":"2009-11-10_20","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_21.log","LogType":"pod","Target":"namspace01","Name":"2009-11-10_21","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_22.log","LogType":"pod","Target":"namspace01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_23.log","LogType":"pod","Target":"namspace01","Name":"2009-11-10_23","Extention":"log","KB":0}]`
	want := []LogFile{
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), NODE_TYPE, "node01", "2009-11-10_20.log"),
			SubPath:   "2009-11-10_20.log",
			LogType:   NODE_TYPE,
			Target:    "node01",
			Name:      "2009-11-10_20.log",
			Extention: ".log",
			Size:      0,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), NODE_TYPE, "node01", "2009-11-10_21.log"),
			SubPath:   "2009-11-10_21.log",
			LogType:   NODE_TYPE,
			Target:    "node01",
			Name:      "2009-11-10_21.log",
			Extention: ".log",
			Size:      0,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), NODE_TYPE, "node01", "2009-11-10_22.log"),
			SubPath:   "2009-11-10_22.log",
			LogType:   NODE_TYPE,
			Target:    "node01",
			Name:      "2009-11-10_22.log",
			Extention: ".log",
			Size:      0,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), NODE_TYPE, "node01", "2009-11-10_23.log"),
			SubPath:   "2009-11-10_23.log",
			LogType:   NODE_TYPE,
			Target:    "node01",
			Name:      "2009-11-10_23.log",
			Extention: ".log",
			Size:      0,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), NODE_TYPE, "node02", "2009-11-10_20.log"),
			SubPath:   "2009-11-10_20.log",
			LogType:   NODE_TYPE,
			Target:    "node02",
			Name:      "2009-11-10_20.log",
			Extention: ".log",
			Size:      0,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), NODE_TYPE, "node02", "2009-11-10_21.log"),
			SubPath:   "2009-11-10_21.log",
			LogType:   NODE_TYPE,
			Target:    "node02",
			Name:      "2009-11-10_21.log",
			Extention: ".log",
			Size:      0,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), NODE_TYPE, "node02", "2009-11-10_22.log"),
			SubPath:   "2009-11-10_22.log",
			LogType:   NODE_TYPE,
			Target:    "node02",
			Name:      "2009-11-10_22.log",
			Extention: ".log",
			Size:      0,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), NODE_TYPE, "node02", "2009-11-10_23.log"),
			SubPath:   "2009-11-10_23.log",
			LogType:   NODE_TYPE,
			Target:    "node02",
			Name:      "2009-11-10_23.log",
			Extention: ".log",
			Size:      0,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), POD_TYPE, "namespace01", "2009-11-10_20.log"),
			SubPath:   "2009-11-10_20.log",
			LogType:   POD_TYPE,
			Target:    "namespace01",
			Name:      "2009-11-10_20.log",
			Extention: ".log",
			Size:      0,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), POD_TYPE, "namespace01", "2009-11-10_21.log"),
			SubPath:   "2009-11-10_21.log",
			LogType:   POD_TYPE,
			Target:    "namespace01",
			Name:      "2009-11-10_21.log",
			Extention: ".log",
			Size:      0,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), POD_TYPE, "namespace01", "2009-11-10_22.log"),
			SubPath:   "2009-11-10_22.log",
			LogType:   POD_TYPE,
			Target:    "namespace01",
			Name:      "2009-11-10_22.log",
			Extention: ".log",
			Size:      0,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), POD_TYPE, "namespace01", "2009-11-10_23.log"),
			SubPath:   "2009-11-10_23.log",
			LogType:   POD_TYPE,
			Target:    "namespace01",
			Name:      "2009-11-10_23.log",
			Extention: ".log",
			Size:      0,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), POD_TYPE, "namespace02", "2009-11-10_20.log"),
			SubPath:   "2009-11-10_20.log",
			LogType:   POD_TYPE,
			Target:    "namespace02",
			Name:      "2009-11-10_20.log",
			Extention: ".log",
			Size:      0,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), POD_TYPE, "namespace02", "2009-11-10_21.log"),
			SubPath:   "2009-11-10_21.log",
			LogType:   POD_TYPE,
			Target:    "namespace02",
			Name:      "2009-11-10_21.log",
			Extention: ".log",
			Size:      0,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), POD_TYPE, "namespace02", "2009-11-10_22.log"),
			SubPath:   "2009-11-10_22.log",
			LogType:   POD_TYPE,
			Target:    "namespace02",
			Name:      "2009-11-10_22.log",
			Extention: ".log",
			Size:      0,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), POD_TYPE, "namespace02", "2009-11-10_23.log"),
			SubPath:   "2009-11-10_23.log",
			LogType:   POD_TYPE,
			Target:    "namespace02",
			Name:      "2009-11-10_23.log",
			Extention: ".log",
			Size:      0,
		},
	}
	assert.Equal(t, want, got)
}

func Test_DeleteBySize(t *testing.T) {
	testutil.SetTestLogFiles()

	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp2", "log")})
	rotator := Rotator{driver: d}

	rotator.DeleteBySize(false)
	want := `[{"Filepath":"/tmp/log/node/node01/2009-11-10_22.log","LogType":"node","Target":"node01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node01/2009-11-10_23.log","LogType":"node","Target":"node01","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_22.log","LogType":"node","Target":"node02","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_23.log","LogType":"node","Target":"node02","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_21.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_21","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_22.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_23.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_22.log","LogType":"pod","Target":"namspace01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_23.log","LogType":"pod","Target":"namspace01","Name":"2009-11-10_23","Extention":"log","KB":0}]`
	got := rotator.ListFiles()

	testutil.CheckEqualJSON(t, got, want)
	testutil.CheckEqualJSON(t, len(got), "9")
}

func Test_DeleteBySize_DryRun(t *testing.T) {
	testutil.SetTestLogFiles()

	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp", "log")})
	rotator := Rotator{driver: d}

	rotator.DeleteBySize(true)
}
