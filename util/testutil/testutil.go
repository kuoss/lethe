package testutil

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var root string

func ChdirRoot() {
	if root != "" {
		fmt.Printf("ChdirRoot: %s (skipped)\n", root)
		return
	}

	_, filename, _, _ := runtime.Caller(0)
	root = path.Join(path.Dir(filename), "../..")
	fmt.Printf("ChdirRoot: %s\n", root)
	err := os.Chdir(root)
	if err != nil {
		panic(err)
	}
}

func getTestID() string {
	_, filename, _, _ := runtime.Caller(2)
	if strings.HasSuffix(filename, "/init_test.go") {
		return "init"
	}
	rel, err := filepath.Rel(root, filename)
	if err != nil {
		panic(err)
	}
	rel = strings.ReplaceAll(rel, ".go", "")
	rel = strings.ReplaceAll(rel, string(os.PathSeparator), "_")
	return rel
}

func ResetLogData() {
	testID := getTestID()
	logDataPath := "tmp/" + testID
	if testID != "init" {
		fmt.Printf("remove logDataPath: %s\n", logDataPath)
		err := os.RemoveAll(logDataPath)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("fill   logDataPath: %s\n", logDataPath)
	err := copyRecursively("./testdata/log", logDataPath)
	if err != nil {
		panic(err)
	}
}

func copyRecursively(src string, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open err: %w", err)
	}
	file, err := f.Stat()
	if err != nil {
		return fmt.Errorf("stat err: %w", err)
	}
	if !file.IsDir() {
		return fmt.Errorf("not dir: %s", file.Name())
	}
	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return fmt.Errorf("mkdirAll err: %w", err)
	}
	files, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("readDir err: %w", err)
	}
	for _, f := range files {
		srcFile := src + "/" + f.Name()
		destFile := dest + "/" + f.Name()
		// dir
		if f.IsDir() {
			err := copyRecursively(srcFile, destFile)
			if err != nil {
				return fmt.Errorf("copyRecursively err: %w", err)
			}
			continue
		}
		// file
		content, err := os.ReadFile(srcFile)
		if err != nil {
			return fmt.Errorf("readFile err: %w", err)
		}
		err = os.WriteFile(destFile, content, 0755)
		if err != nil {
			return fmt.Errorf("writeFile err: %w", err)
		}
	}
	return nil
}
