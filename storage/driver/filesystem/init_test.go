package filesystem

import (
	"log"
	"os"
	"path"
	"runtime"
)

func init() {
	changeWorkingDirectoryToProjectRoot()
}

func changeWorkingDirectoryToProjectRoot() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../..")
	err := os.Chdir(dir)
	if err != nil {
		log.Fatalf("Cannot change directory to [%s]", dir)
	}
}
