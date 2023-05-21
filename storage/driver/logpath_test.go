package driver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	logPath1 = LogPath{
		RootDirectory: "/tmp/init",
		Subpath:       "Subpath",
		target:        "target",
		logType:       "logType",
		file:          "file",
	}
)

func TestGetDepth(t *testing.T) {
	testCases := []struct {
		rootDirectory string
		subpath       string
		want          Depth
		wantLogPath   LogPath
	}{
		{
			"", "",
			DepthType, LogPath{RootDirectory: "", Subpath: "", target: "", logType: "", file: ""},
		},
		{
			"log", "",
			DepthType, LogPath{RootDirectory: "log", Subpath: "", target: "", logType: "", file: ""},
		},
		{
			"", "log",
			DepthType, LogPath{RootDirectory: "", Subpath: "log", target: "", logType: "log", file: ""},
		},
		{
			"tmp/init", "",
			DepthType, LogPath{RootDirectory: "tmp/init", Subpath: "", target: "", logType: "", file: ""},
		},
		{
			"tmp/init", "node",
			DepthType, LogPath{RootDirectory: "tmp/init", Subpath: "node", target: "", logType: "node", file: ""},
		},
		{
			"tmp/init", "hello",
			DepthType, LogPath{RootDirectory: "tmp/init", Subpath: "hello", target: "", logType: "hello", file: ""},
		},
		{
			"tmp/init", "pod",
			DepthType, LogPath{RootDirectory: "tmp/init", Subpath: "pod", target: "", logType: "pod", file: ""},
		},
		{
			"tmp/init", "pod/ns1",
			DepthTarget, LogPath{RootDirectory: "tmp/init", Subpath: "pod/ns1", target: "ns1", logType: "pod", file: ""},
		},
		{
			"tmp/init", "pod/ns1/2022",
			DepthFile, LogPath{RootDirectory: "tmp/init", Subpath: "pod/ns1/2022", target: "ns1", logType: "pod", file: "2022"},
		},
		// error
		{
			"/a", "./pod/ns1/pod1",
			DepthUnknown, LogPath{RootDirectory: "/a", Subpath: "./pod/ns1/pod1", target: "", logType: "", file: ""},
		},
		{
			"log", "pod/ns1/pod1/asdf",
			DepthUnknown, LogPath{RootDirectory: "log", Subpath: "pod/ns1/pod1/asdf", target: "", logType: "", file: ""},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d  %s  %s", i, tc.rootDirectory, tc.subpath), func(t *testing.T) {
			logPath := LogPath{RootDirectory: tc.rootDirectory, Subpath: tc.subpath}
			got := logPath.Depth()
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantLogPath, logPath)
		})
	}
}

func TestLogType(t *testing.T) {
	got := logPath1.LogType()
	assert.Equal(t, "logType", got)
}

func TestTarget(t *testing.T) {
	got := logPath1.Target()
	assert.Equal(t, "target", got)
}

func TestFilename(t *testing.T) {
	got := logPath1.Filename()
	assert.Equal(t, "file", got)
}
