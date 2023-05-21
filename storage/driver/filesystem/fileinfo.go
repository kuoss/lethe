package filesystem

import (
	"os"
	"time"
	// "github.com/kuoss/lethe/storagedriver/types"
)

// var _ types.FileInfo = FileInfo{}

type FileInfo struct {
	osFileInfo os.FileInfo
	fullpath   string
}

func (i FileInfo) Fullpath() string {
	return i.fullpath
}

func (i FileInfo) Size() int64 {
	if i.IsDir() {
		return 0
	}
	return i.osFileInfo.Size()
}

func (i FileInfo) ModTime() time.Time {
	return i.osFileInfo.ModTime()
}

func (fi FileInfo) IsDir() bool {
	return fi.osFileInfo.IsDir()
}
