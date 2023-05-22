package driver

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kuoss/common/logger"
)

// <RootDirectory>/<LogType>/<target>/<logfile>
// example) RootDirectory = /tmp/log
// path
//  /tmp
//    └── log
//        ├── node
//        │     ├── node01
//        │     │     ├── 2009-11-10_22.log
//        │     │     └── 2009-11-10_23.log
//        │     └── node02
//        │         └── 2009-11-10_22.log
//        └── pod
//            ├── namespace01
//            │     ├── 2009-11-10_22.log
//            │     └── 2009-11-10_23.log
//            └── namespace02
//                └── 2009-11-10_22.log

type Depth int

const (
	DepthUnknown Depth = iota
	DepthType
	DepthTarget
	DepthFile
)

type LogPath struct {
	RootDirectory string
	Subpath       string
	target        string
	logType       string
	file          string
}

func (l *LogPath) Depth() Depth {
	parts := strings.Split(l.Subpath, string(os.PathSeparator))
	switch len(parts) {
	case 1:
		l.logType = parts[0]
		return DepthType
	case 2:
		l.logType = parts[0]
		l.target = parts[1]
		return DepthTarget
	case 3:
		l.logType = parts[0]
		l.target = parts[1]
		l.file = parts[2]
		return DepthFile
	}
	logger.Warnf("path is too deep. depth: %d, rel: %s", len(parts), l.Subpath)
	return DepthUnknown
}

func (l *LogPath) Fullpath() string {
	return filepath.Join(l.RootDirectory, l.Subpath)
}
func (l *LogPath) LogType() string {
	return l.logType
}
func (l *LogPath) Target() string {
	return l.target
}
func (l *LogPath) Filename() string {
	return l.file
}
