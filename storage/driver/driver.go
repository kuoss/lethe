package driver

import (
	"io"
)

// https://github.com/distribution/distribution/blob/v2.8.2/registry/storage/driver/storagedriver.go#L41
type Driver interface {
	Name() string
	GetContent(path string) ([]byte, error)
	PutContent(path string, content []byte) error
	Reader(path string) (io.ReadCloser, error)
	Writer(path string) (FileWriter, error)
	Stat(path string) (FileInfo, error)
	List(path string) ([]string, error)
	Move(sourcePath string, destPath string) error
	Delete(path string) error
	Walk(path string) ([]FileInfo, error)
	WalkDir(path string) ([]string, error)
	Mkdir(path string) error
	RootDirectory() string
}
