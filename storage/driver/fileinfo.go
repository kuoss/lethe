package driver

import (
	"time"
)

type FileInfo interface {
	Fullpath() string
	Size() int64
	ModTime() time.Time
	IsDir() bool
}
