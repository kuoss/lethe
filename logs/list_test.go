package logs

import (
	"github.com/kuoss/lethe/storage/driver/factory"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"

	"github.com/kuoss/lethe/testutil"
)

func TestListFiles(t *testing.T) {
	testutil.SetTestLogs()

	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp", "log")})
	rotator := Rotator{driver: d}

	want := []LogFile{
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), NODE_TYPE, "node01", "2009-11-10_22.log"),
			SubPath:   "2009-11-10_22.log",
			LogType:   NODE_TYPE,
			Target:    "node01",
			Name:      "2009-11-10_22.log",
			Extention: ".log",
			Size:      1057,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), NODE_TYPE, "node01", "2009-11-10_23.log"),
			SubPath:   "2009-11-10_23.log",
			LogType:   NODE_TYPE,
			Target:    "node01",
			Name:      "2009-11-10_23.log",
			Extention: ".log",
			Size:      177,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), NODE_TYPE, "node02", "2009-11-10_22.log"),
			SubPath:   "2009-11-10_22.log",
			LogType:   NODE_TYPE,
			Target:    "node02",
			Name:      "2009-11-10_22.log",
			Extention: ".log",
			Size:      1116,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), POD_TYPE, "namespace01", "2009-11-10_22.log"),
			SubPath:   "2009-11-10_22.log",
			LogType:   POD_TYPE,
			Target:    "namespace01",
			Name:      "2009-11-10_22.log",
			Extention: ".log",
			Size:      1031,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), POD_TYPE, "namespace01", "2009-11-10_23.log"),
			SubPath:   "2009-11-10_23.log",
			LogType:   POD_TYPE,
			Target:    "namespace01",
			Name:      "2009-11-10_23.log",
			Extention: ".log",
			Size:      279,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), POD_TYPE, "namespace02", "2009-11-10_22.log"),
			SubPath:   "2009-11-10_22.log",
			LogType:   POD_TYPE,
			Target:    "namespace02",
			Name:      "2009-11-10_22.log",
			Extention: ".log",
			Size:      1125,
		},
	}
	got := rotator.ListFiles()
	assert.Equal(t, want, got)
}

func TestListDirs(t *testing.T) {
	testutil.SetTestLogs()

	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp", "log")})
	rotator := Rotator{driver: d}
	want := []LogDir{
		{
			FullPath: filepath.Join(rotator.driver.RootDirectory(), NODE_TYPE, "node01"),
			SubPath:  filepath.Join(NODE_TYPE, "node01"),
			LogType:  NODE_TYPE,
			Target:   "node01",
		},
		{
			FullPath: filepath.Join(rotator.driver.RootDirectory(), NODE_TYPE, "node02"),
			SubPath:  filepath.Join(NODE_TYPE, "node02"),
			LogType:  NODE_TYPE,
			Target:   "node02",
		},
		{
			FullPath: filepath.Join(rotator.driver.RootDirectory(), POD_TYPE, "namespace01"),
			SubPath:  filepath.Join(POD_TYPE, "namespace01"),
			LogType:  POD_TYPE,
			Target:   "namespace01",
		},
		{
			FullPath: filepath.Join(rotator.driver.RootDirectory(), POD_TYPE, "namespace02"),
			SubPath:  filepath.Join(POD_TYPE, "namespace02"),
			LogType:  POD_TYPE,
			Target:   "namespace02",
		},
	}
	got := rotator.ListDirs()
	assert.Equal(t, want, got)
}

func TestListDirWithSize(t *testing.T) {
	testutil.SetTestLogs()

	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp", "log")})
	rotator := Rotator{driver: d}

	want := []LogDir{
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), NODE_TYPE, "node01"),
			SubPath:   filepath.Join(NODE_TYPE, "node01"),
			LogType:   NODE_TYPE,
			Target:    "node01",
			FileCount: 2,
			FirstFile: "2009-11-10_22.log",
			LastFile:  "2009-11-10_23.log",
			Size:      1234,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), NODE_TYPE, "node02"),
			SubPath:   filepath.Join(NODE_TYPE, "node02"),
			LogType:   NODE_TYPE,
			Target:    "node02",
			FileCount: 1,
			FirstFile: "2009-11-10_22.log",
			LastFile:  "2009-11-10_22.log",
			Size:      1116,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), POD_TYPE, "namespace01"),
			SubPath:   filepath.Join(POD_TYPE, "namespace01"),
			LogType:   POD_TYPE,
			Target:    "namespace01",
			FileCount: 2,
			FirstFile: "2009-11-10_22.log",
			LastFile:  "2009-11-10_23.log",
			Size:      1310,
		},
		{
			FullPath:  filepath.Join(rotator.driver.RootDirectory(), POD_TYPE, "namespace02"),
			SubPath:   filepath.Join(POD_TYPE, "namespace02"),
			LogType:   POD_TYPE,
			Target:    "namespace02",
			FileCount: 1,
			FirstFile: "2009-11-10_22.log",
			LastFile:  "2009-11-10_22.log",
			Size:      1125,
		},
	}
	got := rotator.ListDirsWithSize()
	assert.Equal(t, want, got)
}

func TestListTargets(t *testing.T) {
	testutil.SetTestLogs()

	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp", "log")})
	rotator := Rotator{driver: d}
	want := []LogDir{
		{
			FullPath:    filepath.Join(rotator.driver.RootDirectory(), NODE_TYPE, "node01"),
			SubPath:     filepath.Join(NODE_TYPE, "node01"),
			LogType:     NODE_TYPE,
			Target:      "node01",
			FileCount:   2,
			FirstFile:   "2009-11-10_22.log",
			LastFile:    "2009-11-10_23.log",
			Size:        1234,
			LastForward: "2009-11-10T23:00:00.",
		},
		{
			FullPath:    filepath.Join(rotator.driver.RootDirectory(), NODE_TYPE, "node02"),
			SubPath:     filepath.Join(NODE_TYPE, "node02"),
			LogType:     NODE_TYPE,
			Target:      "node02",
			FileCount:   1,
			FirstFile:   "2009-11-10_22.log",
			LastFile:    "2009-11-10_22.log",
			Size:        1116,
			LastForward: "2009-11-10T22:58:00.",
		},
		{
			FullPath:    filepath.Join(rotator.driver.RootDirectory(), POD_TYPE, "namespace01"),
			SubPath:     filepath.Join(POD_TYPE, "namespace01"),
			LogType:     POD_TYPE,
			Target:      "namespace01",
			FileCount:   2,
			FirstFile:   "2009-11-10_22.log",
			LastFile:    "2009-11-10_23.log",
			Size:        1310,
			LastForward: "2009-11-10T23:00:00.",
		},
		{
			FullPath:    filepath.Join(rotator.driver.RootDirectory(), POD_TYPE, "namespace02"),
			SubPath:     filepath.Join(POD_TYPE, "namespace02"),
			LogType:     POD_TYPE,
			Target:      "namespace02",
			FileCount:   1,
			FirstFile:   "2009-11-10_22.log",
			LastFile:    "2009-11-10_22.log",
			Size:        1125,
			LastForward: "2009-11-10T22:58:00.",
		},
	}
	got := rotator.ListTargets()
	assert.Equal(t, want, got)
}
