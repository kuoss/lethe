package driver

import (
	"io"
)

type FileWriter interface {
	io.WriteCloser
	Size() int64
	Cancel() error
	Commit() error
}
