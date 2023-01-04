package storage

import (
	"os"
	"path/filepath"
	"strings"
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

const (
	TYPE    = "TYPE"
	TARGET  = "TARGET"
	FILE    = "FILE"
	UNKNOWN = "UNKNOWN"
)

/*
type pathSpec() interface{
	Stub()
}
*/

type LogPath struct {
	fullPath      string
	target        string
	logType       string
	file          string
	RootDirectory string
}

func (l *LogPath) Depth() string {
	rel, err := filepath.Rel(l.RootDirectory, l.fullPath)
	if err != nil {
		return UNKNOWN
	}
	parts := strings.Split(rel, string(os.PathSeparator))
	switch len(parts) {
	case 1:
		l.logType = parts[0]
		return TYPE
	case 2:
		l.logType = parts[0]
		l.target = parts[1]
		return TARGET

	case 3:
		l.logType = parts[0]
		l.target = parts[1]
		l.file = parts[2]
		return FILE
	}
	return UNKNOWN
}

func (l *LogPath) FullPath() string {
	return l.fullPath
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

//just for testing?
func (l *LogPath) SetFullPath(subPath string) {
	l.fullPath = filepath.Join(l.RootDirectory, subPath)
}
