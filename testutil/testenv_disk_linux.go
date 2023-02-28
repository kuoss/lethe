//go:build linux
// +build linux

package testutil

func getDiskAvailableBytes(path string) (string, error) {
	// get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("cannot get absoluth path for [%s]: %s", path, err)
	}

	// find mount
	// ignoring mountpoint "/init" because it is not a directory
	mount, err := fscryptFilesystem.FindMount(absPath)
	if err != nil {
		return "", fmt.Errorf("cannot find mount for [%s]: %s", absPath, err)
	}

	// get disk available bytes
	var stat unix.Statfs_t
	err = unix.Statfs(mount.Path, &stat)
	if err != nil {
		return "", fmt.Errorf("cannot get disk available bytes for [%s]: %s", mount.Path, err)
	}
	return fmt.Sprintf("%d", stat.Bavail*uint64(stat.Bsize)), nil
}
