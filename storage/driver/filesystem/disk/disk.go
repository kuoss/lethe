//go:build linux
// +build linux

package disk

import (
	"fmt"
	"path/filepath"

	fscryptFilesystem "github.com/google/fscrypt/filesystem"
	"golang.org/x/sys/unix"
)

func getDiskAvailableBytes(path string) (string, error) {
	// get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("cannot get absoluth path for [%s]: %w", path, err)
	}

	// find mount
	// ignoring mountpoint "/init" because it is not a directory
	mount, err := fscryptFilesystem.FindMount(absPath)
	if err != nil {
		return "", fmt.Errorf("cannot find mount for [%s]: %w", absPath, err)
	}

	// get disk available bytes
	var stat unix.Statfs_t
	err = unix.Statfs(mount.Path, &stat)
	if err != nil {
		return "", fmt.Errorf("cannot get disk available bytes for [%s]: %w", mount.Path, err)
	}
	return fmt.Sprintf("%d", stat.Bavail*uint64(stat.Bsize)), nil
}
