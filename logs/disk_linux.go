//go:build linux
// +build linux

package logs

import "golang.org/x/sys/unix"

func (rotator *Rotator) GetDiskUsedBytes(path string) (int, error) {
	if os.Getenv("TEST_MODE") == "1" {
		return rotator.GetDiskUsedBytesInTest(path)
	}
	var stat unix.Statfs_t
	err := unix.Statfs(path, &stat)
	if err != nil {
		return -1, err
	}
	return int(int64(stat.Blocks-stat.Bavail) * stat.Bsize), nil
}

func (rotator *Rotator) GetDiskUsedBytesInTest(path string) (int, error) {
	var stat unix.Statfs_t
	err := unix.Statfs(path, &stat)
	if err != nil {
		return -1, err
	}
	avail := int(int64(stat.Bavail) * stat.Bsize)
	initialAvail, err := strconv.Atoi(os.Getenv("TEST_INITIAL_DISK_AVAILABLE_BYTES"))
	if err != nil {
		return -1, err
	}
	return initialAvail - avail, nil
}
