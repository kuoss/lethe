package driver

import (
	"io"
)

type Driver interface {
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
	RootDirectory() string
}
