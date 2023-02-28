//go:build windows
// +build windows

package logs

import (
	"os"
	"strconv"

	"golang.org/x/sys/windows"
)

func (rotator *Rotator) GetDiskUsedBytes(path string) (int, error) {
	var free, total, avail uint64
	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}
	err = windows.GetDiskFreeSpaceEx(pathPtr, &free, &total, &avail)
	if err != nil {
		return 0, err
	}
	return int(total - avail), nil
}
