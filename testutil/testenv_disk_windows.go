//go:build windows
// +build windows

package testutil

import (
	"fmt"
	// "golang.org/x/sys/windows"
)

func getDiskAvailableBytes(path string) (string, error) {
	fmt.Println("Currently getDiskAvailableBytes is not supported on Windows.")
	return "99999999", nil
	// var free, total, available uint64
	// pathPtr, err := windows.UTF16PtrFromString(path)
	// if err != nil {
	// 	return "", fmt.Errorf("cannot get utf16ptr from string [%s]: %s", path, err)
	// }
	// err = windows.GetDiskFreeSpaceEx(pathPtr, &free, &total, &available)
	// if err != nil {
	// 	return "", fmt.Errorf("cannot get disk free space for [%s]: %s", path, err)
	// }
	// return fmt.Sprintf("%d", available), nil
}
