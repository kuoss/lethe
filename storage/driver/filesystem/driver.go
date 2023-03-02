package filesystem

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/kuoss/lethe/storage"
	storagedriver "github.com/kuoss/lethe/storage/driver"
	"github.com/kuoss/lethe/storage/driver/factory"
)

const (
	driverName           = "filesystem"
	defaultRootDirectory = "/tmp/log"
)

type DriverParameters struct {
	RootDirectory string
}

func init() {
	factory.Register(driverName, &filesystemDriverFactory{})
}

type driver struct {
	rootDirectory string
}

func New(params DriverParameters) storagedriver.StorageDriver {
	//return &driver{rootDirectory: defaultRootDirectory}
	if params.RootDirectory != "" {
		return &driver{rootDirectory: params.RootDirectory}
	}

	return &driver{rootDirectory: defaultRootDirectory}
}
func (d *driver) RootDirectory() string {
	return d.rootDirectory
}

func (d *driver) Name() string {
	return driverName
}

func (d *driver) fullPath(subPath string) string {
	return path.Join(d.rootDirectory, subPath)
}

func (d *driver) GetContent(path string) ([]byte, error) {
	rc, err := d.Reader(path)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	p, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (d *driver) PutContent(subPath string, content []byte) error {
	writer, err := d.Writer(subPath)
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = io.Copy(writer, bytes.NewReader(content))
	if err != nil {
		writer.Cancel()
		return err
	}
	return writer.Commit()
}

func (d *driver) Reader(path string) (io.ReadCloser, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, storage.PathNotFoundError{Path: path}
		}
		return nil, err
	}
	//TODO seek

	return file, nil
}

func (d *driver) Writer(subPath string) (storagedriver.FileWriter, error) {
	fullPath := d.fullPath(subPath)
	parentDir := path.Dir(fullPath)
	if err := os.MkdirAll(parentDir, 0777); err != nil {
		return nil, err
	}

	fp, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	writer := &fileWriter{
		file: fp,
		bw:   bufio.NewWriter(fp),
	}
	return writer, nil
}

func (d *driver) Stat(path string) (storagedriver.FileInfo, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, storage.PathNotFoundError{Path: path}
		}
		return nil, err
	}

	return fileInfo{
		path:     path,
		FileInfo: fi,
	}, nil
}

func (d *driver) List(path string) ([]string, error) {

	dir, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, storage.PathNotFoundError{Path: path}
		}
		return nil, err
	}
	defer dir.Close()

	fileNames, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(fileNames))
	for _, fileName := range fileNames {
		keys = append(keys, filepath.Join(path, fileName))
	}
	return keys, nil
}

func (d *driver) Move(sourcePath, destPath string) error {
	source := d.fullPath(sourcePath)
	dest := d.fullPath(destPath)

	if _, err := os.Stat(source); os.IsNotExist(err) {
		return storage.PathNotFoundError{Path: sourcePath}
	}

	if err := os.MkdirAll(path.Dir(dest), 0755); err != nil {
		return err
	}

	// TODO check windows
	//Rename replaces it
	err := os.Rename(source, dest)
	return err
}

func (d *driver) Delete(path string) error {
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err != nil {
		return storage.PathNotFoundError{Path: path}
	}
	err = os.RemoveAll(path)
	return err
}

// return only files
func (d *driver) Walk(from string) ([]storagedriver.FileInfo, error) {

	var infos []storagedriver.FileInfo

	err := filepath.Walk(from, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() {
			infos = append(infos, fileInfo{
				FileInfo: info,
				path:     path,
			})
		}
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", from, err)
		return []storagedriver.FileInfo{}, err
	}
	return infos, nil
}

// WalkDir method return only directories
func (d *driver) WalkDir(from string) ([]string, error) {
	dirs := []string{}

	err := filepath.WalkDir(from, func(path string, dir fs.DirEntry, err error) error {

		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		//fmt.Printf("visited logs or dir: %q\n", path)
		if dir.IsDir() {
			rel, _ := filepath.Rel(from, path)
			if err != nil {
				return err
			}
			dirs = append(dirs, rel)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", from, err)
		return dirs, err
	}
	return dirs, nil
}

//todo if listing every file for acquire target info makes performance issue, the consider interface only directories, not file
/*
func (d *driver) WalkDirWithDepth(from string, depth int) ([]string, error) {
	dirs := []string{}

	err := filepath.WalkDir(from, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		rel, _ := filepath.Rel(from, path)
		if strings.Count(rel, string(os.PathSeparator)) > depth {
			return filepath.SkipDir
		}
		dirs = append(dirs, path)
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", from, err)
		return dirs, err
	}
	return dirs, nil
}

*/

// For object-storage backend
type fileWriter struct {
	file      *os.File
	size      int64
	bw        *bufio.Writer
	closed    bool
	committed bool
	cancelled bool
}

func (fw *fileWriter) Write(p []byte) (int, error) {
	if fw.closed {
		return 0, fmt.Errorf("already closed")
	} else if fw.committed {
		return 0, fmt.Errorf("already committed")
	} else if fw.cancelled {
		return 0, fmt.Errorf("already cancelled")
	}
	n, err := fw.bw.Write(p)
	fw.size += int64(n)
	return n, err
}

func (fw *fileWriter) Size() int64 {
	return fw.size
}

func (fw *fileWriter) Close() error {
	if fw.closed {
		return fmt.Errorf("already closed")
	}

	if err := fw.bw.Flush(); err != nil {
		return err
	}

	if err := fw.file.Sync(); err != nil {
		return err
	}

	if err := fw.file.Close(); err != nil {
		return err
	}
	fw.closed = true
	return nil
}

func (fw *fileWriter) Cancel() error {
	if fw.closed {
		return fmt.Errorf("already closed")
	}

	fw.cancelled = true
	fw.file.Close()
	return os.Remove(fw.file.Name())
}

func (fw *fileWriter) Commit() error {
	if fw.closed {
		return fmt.Errorf("already closed")
	} else if fw.committed {
		return fmt.Errorf("already committed")
	} else if fw.cancelled {
		return fmt.Errorf("already cancelled")
	}

	if err := fw.bw.Flush(); err != nil {
		return err
	}

	if err := fw.file.Sync(); err != nil {
		return err
	}

	fw.committed = true
	return nil
}

// for compile
var _ storagedriver.FileInfo = fileInfo{}

type fileInfo struct {
	os.FileInfo
	path string
}

func (fi fileInfo) Path() string {
	return fi.path
}

func (fi fileInfo) Size() int64 {
	if fi.IsDir() {
		return 0
	}

	return fi.FileInfo.Size()
}

func (fi fileInfo) ModTime() time.Time {
	return fi.FileInfo.ModTime()
}

func (fi fileInfo) IsDir() bool {
	return fi.FileInfo.IsDir()
}
