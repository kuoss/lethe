package testutil

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"runtime"

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
	logRoot := "./tmp/log"
	os.Setenv("TEST_MODE", "1")
	changeWorkingDirectoryToProjectRoot()

	config.LoadConfig()
	config.GetConfig().Set("retention.time", "3h")
	config.GetConfig().Set("retention.size", "10m")
	config.GetConfig().Set("retention.sizingStrategy", "files")
	config.SetLimit(1000)
	config.SetLogRoot(logRoot)

	ClearTestLogFiles()
	setenvIntialDiskAvailableBytes()
	fmt.Println("Test environment initialized...")
}

func setenvIntialDiskAvailableBytes() {
	if os.Getenv("TEST_INITIAL_DISK_AVAILABLE_BYTES") != "" {
		return
	}
	avail, err := getDiskAvailableBytes(config.GetLogRoot())
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
		log.Fatalf("Cannot change directory to [%s]", dir)
	}
}

func ClearTestLogFiles() {
	logDirectory := config.GetLogRoot()
	fmt.Printf("clear logDirectory: %s\n", logDirectory)
	err := os.RemoveAll(logDirectory)
	if err != nil {
		log.Fatalf("Cannot remove logDirectory [%s]: %s", logDirectory, err)
	}
	os.MkdirAll(logDirectory, 0755)
	time.Sleep(200 * time.Millisecond)
}

func SetTestLogFiles() {
	ClearTestLogFiles()
	logDirectory := config.GetLogRoot()
	fmt.Printf("copy to logDirectory: %s\n", logDirectory)
	copyRecursively("./testutil/log", logDirectory)
}

func copyRecursively(src string, dest string) {
	f, err := os.Open(src)
	if err != nil {
		log.Fatalf("Cannot open [%s]: %s", src, err)
	}
	file, err := f.Stat()
	if err != nil {
		log.Fatalf("Cannot stat [%s]: %s", file, err)
	}
	if !file.IsDir() {
		log.Fatalf("Source [%s] is not a directory: %s", file.Name(), err)
	}
	// log.Println("make directory:", dest)
	err = os.MkdirAll(dest, 0755)
	if err != nil {
		log.Fatalf("Cannot make directory [%s]: %s", dest, err)
	}
	files, err := ioutil.ReadDir(src)
	if err != nil {
		log.Fatalf("Cannot read directory [%s]: %s", dest, err)
	}
	for _, f := range files {
		srcFile := src + "/" + f.Name()
		destFile := dest + "/" + f.Name()
		if f.IsDir() {
			copyRecursively(srcFile, destFile)
			continue
		}
		// log.Println("copy file:", srcFile, destFile)
		content, err := ioutil.ReadFile(srcFile)
		if err != nil {
			log.Fatalf("Cannot read file [%s]: %s", srcFile, err)
		}
		err = ioutil.WriteFile(destFile, content, 0755)
		if err != nil {
			log.Fatalf("Cannot write file [%s]: %s", destFile, err)
		}
	}
}
