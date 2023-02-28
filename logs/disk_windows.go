//go:build windows
// +build windows

package logs

import (
	"os"
	"strconv"

	"golang.org/x/sys/windows"
)

func (rotator *Rotator) GetDiskUsedBytes(path string) (int, error) {
	if os.Getenv("TEST_MODE") == "1" {
		return rotator.GetDiskUsedBytesInTest(path)
	}
	var free, total, available uint64
	pathPtr, err := windows.UTF16PtrFromString(`C:\Users`)
	if err != nil {
		return -1, err
	}
	err = windows.GetDiskFreeSpaceEx(pathPtr, &free, &total, &available)
	if err != nil {
		return -1, err
	}
	return int(total) - int(available), nil
}

func (rotator *Rotator) GetDiskUsedBytesInTest(path string) (int, error) {
	// initialAvail
	initialAvail, err := strconv.Atoi(os.Getenv("TEST_INITIAL_DISK_AVAILABLE_BYTES"))
	if err != nil {
		return -1, err
	}

	// avail
	var free, total, avail uint64
	pathPtr, err := windows.UTF16PtrFromString(`C:\Users`)
	if err != nil {
		return -1, err
	}
	err = windows.GetDiskFreeSpaceEx(pathPtr, &free, &total, &avail)
	if err != nil {
		return -1, err
	}
	return initialAvail - int(int64(avail)), nil
}
