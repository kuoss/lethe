package testutil

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"runtime"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/config"
)

const (
	POD         = "pod"
	NODE        = "node"
	namespace01 = "namespace01"
	namespace02 = "namespace02"
	node01      = "node01"
	node02      = "node02"
)

func Init() {
	logRoot := filepath.Join(".", "tmp", "log")
	os.Setenv("TEST_MODE", "1")
	changeWorkingDirectoryToProjectRoot()

	err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("error on LoadConfig: %s", err)
	}
	config.Viper().Set("retention.time", "3h")
	config.Viper().Set("retention.size", "10m")
	config.Viper().Set("retention.sizingStrategy", "files")
	config.SetLimit(1000)
	config.SetLogDataPath(logRoot)

	setenvIntialDiskAvailableBytes()
	fmt.Println("Test environment initialized...")
}

func setenvIntialDiskAvailableBytes() {
	if os.Getenv("TEST_INITIAL_DISK_AVAILABLE_BYTES") != "" {
		return
	}
	logDirectory := config.GetLogDataPath()
	_ = os.MkdirAll(logDirectory, 0755)
	avail, err := getDiskAvailableBytes(logDirectory)
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("TEST_INITIAL_DISK_AVAILABLE_BYTES", avail)
}

func changeWorkingDirectoryToProjectRoot() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		log.Fatalf("cannot change directory to [%s]", dir)
	}
}

// func ClearTestLogFiles() {
// 	logDirectory := config.GetLogRoot()
// 	fmt.Printf("clear logDirectory: %s\n", logDirectory)
// 	err := os.RemoveAll(logDirectory)
// 	if err != nil {
// 		log.Fatalf("cannot remove logDirectory [%s]: %s", logDirectory, err)
// 	}
// 	os.MkdirAll(logDirectory, 0755)
// }

func SetTestLogFiles() {
	// ClearTestLogFiles()
	logDirectory := config.GetLogDataPath()
	logger.Infof("SetTestLogFiles: logDirectory=%s", logDirectory)
	err := CopyRecursively("./testutil/log", logDirectory)
	if err != nil {
		logger.Errorf("error on CopyRecursively: %s", err)
	}
}

func CopyRecursively(src string, dest string) error {
	logger.Infof("CopyRecursively... src=%s, dest=%s", src, dest)

	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error on Open: %w", err)
	}
	file, err := f.Stat()
	if err != nil {
		return fmt.Errorf("error on Stat: %w", err)
	}
	if !file.IsDir() {
		return fmt.Errorf("src[%s] is not a dir", file.Name())
	}
	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return fmt.Errorf("error on MkdirAll: %w", err)
	}
	files, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("error on ReadDir: %w", err)
	}
	for _, f := range files {
		srcFile := src + "/" + f.Name()
		destFile := dest + "/" + f.Name()
		// dir
		if f.IsDir() {
			err := CopyRecursively(srcFile, destFile)
			if err != nil {
				logger.Errorf("error on CopyRecursively: %s", err)
			}
			continue
		}
		// file
		content, err := os.ReadFile(srcFile)
		if err != nil {
			logger.Errorf("error on ReadFile: %s", err)
			continue
		}
		err = os.WriteFile(destFile, content, 0755)
		if err != nil {
			logger.Errorf("error on WriteFile: %s", err)
		}
	}
	return nil
}
