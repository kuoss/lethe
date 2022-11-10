package file

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/util"
)

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func ListFiles() []LogFile {
	var logFiles []LogFile
	out, _, err := util.Execute(fmt.Sprintf("cd %s; find -type f -maxdepth 3 -mindepth 3", config.GetLogRoot()))
	if err != nil {
		fmt.Println("Warning: cannot find files")
		return logFiles
	}
	files := strings.Split(strings.TrimSpace(out), "\n")
	separator := regexp.MustCompile("[/.]+")
	for _, file := range files {
		cols := separator.Split(file, -1)
		logFiles = append(logFiles, LogFile{
			Typ:       cols[0],
			Target:    cols[1],
			Name:      cols[2],
			Extention: cols[3],
			SubPath:   file[1:],
			FullPath:  config.GetLogRoot() + file[1:],
		})
	}
	return logFiles
}

func ListFilesWithSize() []LogFile {
	logFiles := ListFiles()
	for i, logFile := range logFiles {
		stat, err := os.Stat(logFile.FullPath)
		if err != nil {
			fmt.Printf("Warning: cannot stat: %s\n", logFile.FullPath)
			continue
		}
		logFiles[i].Size = stat.Size()
	}
	return logFiles
}

func ListDirs() []LogDir {
	var logDirs []LogDir
	out, _, err := util.Execute(fmt.Sprintf("cd %s; find -type d -maxdepth 2 -mindepth 2", config.GetLogRoot()))
	if err != nil {
		fmt.Println("Warning: cannot find dirs")
		return logDirs
	}
	dirs := strings.Split(strings.TrimSpace(out), "\n")
	fmt.Println(dirs)
	separator := regexp.MustCompile("[/.]+")
	for _, dir := range dirs {
		cols := separator.Split(dir, -1)
		logDirs = append(logDirs, LogDir{
			Typ:      cols[0],
			Target:   cols[1],
			SubPath:  dir[1:],
			FullPath: config.GetLogRoot() + dir[1:],
		})
	}
	return logDirs
}

func ListDirsWithSize() []LogDir {
	logDirs := ListDirs()
	for i, logDir := range logDirs {
		var err error
		var size int64
		var files []fs.FileInfo

		size, err = DirSize(logDir.FullPath)
		if err != nil {
			fmt.Printf("Warning: cannot get size of directory: %s\n", logDir.FullPath)
		}
		logDirs[i].Size = size

		files, err = ioutil.ReadDir(logDir.FullPath)
		if err != nil {
			fmt.Printf("Warning: cannot get file count of directory: %s\n", logDir.FullPath)
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

func ListTargets() []LogDir {
	logDirs := ListDirsWithSize()
	for i, logDir := range logDirs {
		if logDir.LastFile == "" {
			continue
		}
		command := fmt.Sprintf(`tail -1 %s/%s`, logDir.FullPath, logDir.LastFile)
		out, _, err := util.Execute(command)
		if err != nil {
			continue
		}
		logDirs[i].LastForward = out[:20]
	}
	return logDirs
}
