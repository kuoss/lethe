package driver

import (
	"io"
	"time"
)

type StorageDriver interface {
	Name() string
	GetContent(string) ([]byte, error)
	PutContent(string, []byte) error
	Reader(string) (io.ReadCloser, error)
	Writer(string) (FileWriter, error)
	Stat(string) (FileInfo, error)
	List(string) ([]string, error)
	Move(string, string) error
	Delete(string) error
	Walk(string) ([]FileInfo, error)
	WalkDir(string) ([]string, error)
	//WalkDirWithDepth(from string, depth int) ([]string, error)
	RootDirectory() string
}

type FileWriter interface {
	io.WriteCloser
	Size() int64
	Cancel() error
	Commit() error
}

type FileInfo interface {
	Path() string
	Size() int64
	ModTime() time.Time
	IsDir() bool
}
