//go:build windows
// +build windows

package disk

// import (
// 	"fmt"
// 	"golang.org/x/sys/windows"
// )

// func getDiskAvailableBytes(path string) (string, error) {
// 	var free, total, available uint64
// 	pathPtr, err := windows.UTF16PtrFromString(path)
// 	if err != nil {
// 		return "", fmt.Errorf("cannot get utf16ptr from string [%s]: %s", path, err)
// 	}
// 	err = windows.GetDiskFreeSpaceEx(pathPtr, &free, &total, &available)
// 	if err != nil {
// 		return "", fmt.Errorf("cannot get disk free space for [%s]: %s", path, err)
// 	}
// 	return fmt.Sprintf("%d", available), nil
// }
