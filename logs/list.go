package logs

import (
	"fmt"
	"github.com/kuoss/lethe/storage"
	"os"
	"path/filepath"
)

// LogFile
type LogFile struct {
	FullPath  string
	SubPath   string
	LogType   string
	Target    string
	Name      string
	Extention string
	Size      int64
}

type LogDir struct {
	FullPath    string
	SubPath     string
	LogType     string
	Target      string
	FileCount   int
	FirstFile   string
	LastFile    string
	Size        int64
	LastForward string
}

func (rotator *Rotator) ListFiles() []LogFile {
	var logFiles []LogFile

	rootDirectory := rotator.driver.RootDirectory()
	fileInfos, err := rotator.driver.Walk(rootDirectory)
	if err != nil {
		return logFiles
	}
	for _, fileInfo := range fileInfos {

		logPath := storage.LogPath{RootDirectory: rootDirectory}
		rel, err := filepath.Rel(logPath.RootDirectory, fileInfo.Path())
		if err != nil {
			return nil
		}
		logPath.SetFullPath(rel)

		if logPath.Depth() == storage.FILE {
			logFiles = append(logFiles, LogFile{
				FullPath:  logPath.FullPath(),
				SubPath:   logPath.Filename(),
				LogType:   logPath.LogType(),
				Target:    logPath.Target(),
				Name:      logPath.Filename(),
				Extention: filepath.Ext(logPath.Filename()),
				Size:      fileInfo.Size(),
			})
		}
	}
	return logFiles
}

// Deprecated?
func (rotator *Rotator) ListFilesWithSize() []LogFile {
	logFiles := rotator.ListFiles()

	return logFiles
}

func (rotator *Rotator) ListDirs() []LogDir {
	var logDirs []LogDir

	rootDirecotry := rotator.driver.RootDirectory()
	directories, err := rotator.driver.WalkDir(rootDirecotry)
	if err != nil {
		fmt.Println(err)
		return logDirs
	}
	for _, dir := range directories {
		logPath := storage.LogPath{RootDirectory: rootDirecotry}
		logPath.SetFullPath(dir)
		if logPath.Depth() == storage.TARGET {
			logDirs = append(logDirs, LogDir{
				LogType:  logPath.LogType(),
				Target:   logPath.Target(),
				SubPath:  dir,
				FullPath: logPath.FullPath(),
			})
		}
	}
	return logDirs
}

func (rotator *Rotator) ListDirsWithSize() []LogDir {
	logDirs := rotator.ListDirs()
	for i, logDir := range logDirs {
		var size int64
		size, err := rotator.DirSize(logDir.FullPath)
		if err != nil {
			fmt.Printf("Warning: cannot get size of directory: %s\n", logDir.FullPath)
		}
		logDirs[i].Size = size

		files, err := os.ReadDir(logDir.FullPath)
		if err != nil {
			fmt.Printf("Warning: cannot get logs count of directory: %s\n", logDir.FullPath)
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

func (rotator *Rotator) DirSize(path string) (int64, error) {

	files, err := rotator.driver.List(path)
	if err != nil {
		return 0, err
	}
	var size int64
	for _, file := range files {
		info, err := rotator.driver.Stat(file)
		if err != nil {
			return 0, err
		}
		size += info.Size()
	}
	return size, err
}

// ListTargets method returns ListDirWithSize() + LastForward(timestamp)
func (rotator *Rotator) ListTargets() []LogDir {

	logDirs := rotator.ListDirsWithSize()
	for i, logDir := range logDirs {
		if logDir.LastFile == "" {
			continue
		}

		b, err := rotator.driver.GetContent(filepath.Join(logDir.FullPath, logDir.LastFile))
		if err != nil {
			fmt.Println(err)
			return nil
		}
		content := string(b)

		// todo
		//  if timestamp ?
		logDirs[i].LastForward = content[:20]
	}
	return logDirs
}
