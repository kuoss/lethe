//go:build linux
// +build linux

package logs

import (
	"golang.org/x/sys/unix"
)

func (rotator *Rotator) GetDiskUsedBytes(path string) (int, error) {
	var stat unix.Statfs_t
	err := unix.Statfs(path, &stat)
	if err != nil {
		return 0, err
	}
	return int(stat.Blocks - stat.Bavail*uint64(stat.Bsize)), nil
}
