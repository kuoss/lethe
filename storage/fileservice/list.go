package fileservice

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kuoss/common/logger"
	storagedriver "github.com/kuoss/lethe/storage/driver"
)

func fullpath2subpath(rootDir string, fullpath string) string {
	subpath, err := filepath.Rel(rootDir, fullpath)
	if err != nil {
		logger.Warnf("rel err: %s", err.Error())
		return ""
	}
	return subpath
}

func (s *FileService) dirSize(subpath string) (int64, error) {
	files, err := s.driver.List(subpath)
	if err != nil {
		return 0, err
	}
	var size int64
	for _, file := range files {
		info, err := s.driver.Stat(file)
		if err != nil {
			return 0, err
		}
		size += info.Size()
	}
	return size, err
}

func (s *FileService) List(subpath string) ([]string, error) {
	dirs, err := s.driver.List(subpath)
	if err != nil {
		return nil, fmt.Errorf("list err: %w", err)
	}
	return dirs, nil
}

func (s *FileService) ListLogDirs() []LogDir {
	dirs, err := s.driver.WalkDir(".")
	if err != nil {
		logger.Warnf("walkDir err: %s", err.Error())
		return []LogDir{}
	}

	rootDir := s.driver.RootDirectory()
	logDirs := []LogDir{}
	for _, dir := range dirs {
		logPath := storagedriver.LogPath{RootDirectory: rootDir, Subpath: dir}

		if logPath.Depth() == storagedriver.DepthTarget {
			logDirs = append(logDirs, LogDir{
				LogType:  logPath.LogType(),
				Target:   logPath.Target(),
				Subpath:  dir,
				Fullpath: logPath.Fullpath(),
			})
		}
	}
	return logDirs
}

func (s *FileService) listLogDirsWithSize() []LogDir {
	logDirs := s.ListLogDirs()
	for i, logDir := range logDirs {
		var size int64
		size, err := s.dirSize(logDir.Subpath)
		if err != nil {
			logger.Warnf("dirSize err: %s, path:%s", err.Error(), logDir.Fullpath)
			continue
		}
		logDirs[i].Size = size

		files, err := os.ReadDir(logDir.Fullpath)
		if err != nil {
			logger.Warnf("readDir err: %s, path:%s", err.Error(), logDir.Fullpath)
			continue
		}
		fileCount := len(files)
		logDirs[i].FileCount = fileCount
		if fileCount > 0 {
			logDirs[i].FirstFile = files[0].Name()
			logDirs[i].LastFile = files[fileCount-1].Name()
		}
	}
	return logDirs
}

func (s *FileService) ListTargets() []LogDir {
	logDirs := s.listLogDirsWithSize()
	for i, logDir := range logDirs {
		if logDir.LastFile == "" {
			continue
		}
		b, err := s.driver.GetContent(filepath.Join(logDir.Subpath, logDir.LastFile))
		if err != nil {
			logger.Warnf("getContent err: %s", err.Error())
			continue
		}
		content := string(b)
		// TODO: if timestamp ?
		logDirs[i].LastForward = content[:19] + "Z"
	}
	return logDirs
}

func (s *FileService) ListFiles() ([]LogFile, error) {
	fileInfos, err := s.driver.Walk(".")
	if err != nil {
		return []LogFile{}, fmt.Errorf("walk err: %e", err)
	}
	logFiles := []LogFile{}
	rootDir := s.driver.RootDirectory()
	for _, fileInfo := range fileInfos {
		logPath := storagedriver.LogPath{RootDirectory: rootDir, Subpath: fullpath2subpath(s.config.LogDataPath(), fileInfo.Fullpath())}
		if logPath.Depth() == storagedriver.DepthFile {
			logFiles = append(logFiles, LogFile{
				Fullpath:  logPath.Fullpath(),
				Subpath:   logPath.Subpath,
				LogType:   logPath.LogType(),
				Target:    logPath.Target(),
				Name:      logPath.Filename(),
				Extension: filepath.Ext(logPath.Filename()),
				Size:      fileInfo.Size(),
			})
		}
	}
	return logFiles, nil
}
