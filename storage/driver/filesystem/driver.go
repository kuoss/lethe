package filesystem

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	storagedriver "github.com/kuoss/lethe/storage/driver"
)

const (
	driverName = "filesystem"
)

type Params struct {
	RootDirectory string
}

type driver struct {
	rootDirectory string
}

func New(params Params) storagedriver.Driver {
	return &driver{rootDirectory: params.RootDirectory}
}

func (d *driver) RootDirectory() string {
	return d.rootDirectory
}

func (d *driver) Name() string {
	return driverName
}

func (d *driver) fullPath(subpath string) string {
	return path.Join(d.rootDirectory, subpath)
}

func (d *driver) GetContent(subpath string) ([]byte, error) {
	rc, err := d.Reader(subpath)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	p, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (d *driver) PutContent(subpath string, content []byte) error {
	writer, err := d.Writer(subpath)
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = io.Copy(writer, bytes.NewReader(content))
	if err != nil {
		err2 := writer.Cancel()
		if err2 != nil {
			return fmt.Errorf("copy err: %w, cancel err: %w", err, err2)
		}
		return fmt.Errorf("copy err: %w", err)
	}
	err = writer.Commit()
	if err != nil {
		return fmt.Errorf("commit err: %w", err)
	}
	return nil
}

func (d *driver) Reader(subpath string) (io.ReadCloser, error) {
	file, err := os.OpenFile(d.fullPath(subpath), os.O_RDONLY, 0644)
	if err != nil {
		return nil, storagedriver.PathNotFoundError{Path: subpath, Err: fmt.Errorf("openFile err: %w", err)}
	}
	// TODO: seek
	return file, nil
}

func (d *driver) Writer(subpath string) (storagedriver.FileWriter, error) {
	fullpath := d.fullPath(subpath)
	dir := path.Dir(fullpath)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return nil, err
	}

	fp, err := os.OpenFile(fullpath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	writer := &fileWriter{
		file: fp,
		bw:   bufio.NewWriter(fp),
	}
	return writer, nil
}

func (d *driver) Stat(subpath string) (storagedriver.FileInfo, error) {
	fullpath := d.fullPath(subpath)
	fi, err := os.Stat(fullpath)
	if err != nil {
		return nil, storagedriver.PathNotFoundError{Path: subpath, Err: fmt.Errorf("stat err: %w", err)}
	}
	return FileInfo{
		fullpath:   fullpath,
		osFileInfo: fi,
	}, nil
}

func (d *driver) List(subpath string) ([]string, error) {
	fullpath := d.fullPath(subpath)
	dir, err := os.Open(fullpath)
	if err != nil {
		return nil, storagedriver.PathNotFoundError{Path: subpath, Err: fmt.Errorf("open err: %w", err)}
	}
	defer dir.Close()

	fileNames, err := dir.Readdirnames(0)
	if err != nil {
		return nil, fmt.Errorf("readdirnames err: %w", err)
	}

	keys := make([]string, 0, len(fileNames))
	for _, fileName := range fileNames {
		keys = append(keys, filepath.Join(subpath, fileName))
	}
	return keys, nil
}

func (d *driver) Move(sourcePath, targetPath string) error {
	if len(strings.Split(sourcePath, string(os.PathSeparator))) < 1 {
		return fmt.Errorf("moving 0-1 depth directory is not allowed")
	}
	source := d.fullPath(sourcePath)
	target := d.fullPath(targetPath)

	if _, err := os.Stat(source); err != nil {
		return storagedriver.PathNotFoundError{Path: sourcePath, Err: fmt.Errorf("stat err: %w", err)}
	}

	if err := os.MkdirAll(path.Dir(target), 0755); err != nil {
		return fmt.Errorf("mkdirAll err: %w", err)
	}

	// TODO check windows
	//Rename replaces it
	err := os.Rename(source, target)
	return err
}

func (d *driver) Delete(subpath string) error {
	if len(strings.Split(subpath, string(os.PathSeparator))) < 2 {
		return fmt.Errorf("deleting 0-1 depth directory is not allowed")
	}
	fullpath := d.fullPath(subpath)
	_, err := os.Stat(fullpath)
	if err != nil {
		return storagedriver.PathNotFoundError{Path: subpath, Err: fmt.Errorf("stat err: %w", err)}
	}
	err = os.RemoveAll(fullpath)
	if err != nil {
		return fmt.Errorf("removeAll err: %w", err)
	}
	return nil
}

// return only files
func (d *driver) Walk(subpath string) ([]storagedriver.FileInfo, error) {
	fullpath := d.fullPath(subpath)
	infos := []storagedriver.FileInfo{}
	err := filepath.Walk(fullpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walkFunc err: %w", err)
		}
		if !info.IsDir() {
			infos = append(infos, FileInfo{info, path})
		}
		return nil
	})
	if err != nil {
		return infos, fmt.Errorf("walk err: %w", err)
	}
	return infos, nil
}

// WalkDir method return only directories
func (d *driver) WalkDir(subpath string) ([]string, error) {
	fullpath := d.fullPath(subpath)
	dirs := []string{}
	err := filepath.WalkDir(fullpath, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walkDirFunc err: %w", err)
		}
		if !dir.IsDir() {
			return nil
		}
		subpath, err := filepath.Rel(fullpath, path)
		if err != nil {
			return fmt.Errorf("rel err: %w", err)
		}
		dirs = append(dirs, subpath)
		return nil
	})
	if err != nil {
		return dirs, fmt.Errorf("walkdir err: %w", err)
	}
	return dirs, nil
}

func (d *driver) Mkdir(subpath string) error {

	if err := os.MkdirAll(subpath, 0777); err != nil {
		return err
	}
	return nil
}
