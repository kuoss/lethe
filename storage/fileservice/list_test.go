package fileservice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullpath2subpath(t *testing.T) {
	testCases := []struct {
		rootDir  string
		fullpath string
		want     string
	}{
		{"", "", "."},
		{"hello", "", "../."},
		{"", "hello", "hello"},
		{"tmp/init", "tmp/init/pod", "pod"},
		{"tmp/init", "tmp/init/pod/ns1", "pod/ns1"},
		{"tmp/init", "tmp/init/pod/ns1/2023", "pod/ns1/2023"},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got := fullpath2subpath(tc.rootDir, tc.fullpath)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDirSize(t *testing.T) {
	testCases := []struct {
		path      string
		want      int64
		wantError string
	}{
		{"", 0, ""},
		{"hello", 0, "Path not found: hello"},
		{"node", 0, ""},
		{"pod", 0, ""},
		{"node/node01", 1234, ""},
		{"node/node02", 1116, ""},
		{"pod/namespace01", 2620, ""},
		{"pod/namespace02", 1137, ""},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got, err := fileService.dirSize(tc.path)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestList(t *testing.T) {
	testCases := []struct {
		subpath   string
		want      []string
		wantError string
	}{
		{
			"",
			[]string{"pod", "node"},
			"",
		},
		{
			"hello",
			nil,
			"list err: Path not found: hello",
		},
		{
			"node",
			[]string{"node/node01", "node/node02"},
			"",
		},
		{
			"pod",
			[]string{"pod/namespace01", "pod/namespace02"},
			"",
		},
		{
			"pod/namespace01",
			[]string{"pod/namespace01/2029-11-10_23.log", "pod/namespace01/2009-11-10_22.log", "pod/namespace01/2000-01-01_00.log", "pod/namespace01/2009-11-10_21.log"},
			"",
		},
		{
			"pod/namespace01/2029-11-10_23.log",
			nil,
			"list err: readdirnames err: readdirent tmp/init/pod/namespace01/2029-11-10_23.log: not a directory",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got, err := fileService.List(tc.subpath)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestListLogDirs(t *testing.T) {
	// TODO: if runtime.GOOS == "windows"
	want := []LogDir{
		{Fullpath: "tmp/init/node/node01", Subpath: "node/node01", LogType: "node", Target: "node01", FileCount: 0, FirstFile: "", LastFile: "", Size: 0, LastForward: ""},
		{Fullpath: "tmp/init/node/node02", Subpath: "node/node02", LogType: "node", Target: "node02", FileCount: 0, FirstFile: "", LastFile: "", Size: 0, LastForward: ""},
		{Fullpath: "tmp/init/pod/namespace01", Subpath: "pod/namespace01", LogType: "pod", Target: "namespace01", FileCount: 0, FirstFile: "", LastFile: "", Size: 0, LastForward: ""},
		{Fullpath: "tmp/init/pod/namespace02", Subpath: "pod/namespace02", LogType: "pod", Target: "namespace02", FileCount: 0, FirstFile: "", LastFile: "", Size: 0, LastForward: ""}}
	got := fileService.ListLogDirs()
	assert.Equal(t, want, got)
}

func TestListLogDirsWithSize(t *testing.T) {
	want := []LogDir{
		{Fullpath: "tmp/init/node/node01", Subpath: "node/node01", LogType: "node", Target: "node01", FileCount: 2, FirstFile: "2009-11-10_21.log", LastFile: "2009-11-10_22.log", Size: 1234, LastForward: ""},
		{Fullpath: "tmp/init/node/node02", Subpath: "node/node02", LogType: "node", Target: "node02", FileCount: 2, FirstFile: "2009-11-01_00.log", LastFile: "2009-11-10_21.log", Size: 1116, LastForward: ""},
		{Fullpath: "tmp/init/pod/namespace01", Subpath: "pod/namespace01", LogType: "pod", Target: "namespace01", FileCount: 4, FirstFile: "2000-01-01_00.log", LastFile: "2029-11-10_23.log", Size: 2620, LastForward: ""},
		{Fullpath: "tmp/init/pod/namespace02", Subpath: "pod/namespace02", LogType: "pod", Target: "namespace02", FileCount: 2, FirstFile: "0000-00-00_00.log", LastFile: "2009-11-10_22.log", Size: 1137, LastForward: ""}}
	got := fileService.listLogDirsWithSize()
	assert.Equal(t, want, got)
}

func TestListTargets(t *testing.T) {
	// TODO: if runtime.GOOS == "windows"
	want := []LogDir{
		{Fullpath: "tmp/init/node/node01", Subpath: "node/node01", LogType: "node", Target: "node01", FileCount: 2, FirstFile: "2009-11-10_21.log", LastFile: "2009-11-10_22.log", Size: 1234, LastForward: "2009-11-10T23:00:00Z"},
		{Fullpath: "tmp/init/node/node02", Subpath: "node/node02", LogType: "node", Target: "node02", FileCount: 2, FirstFile: "2009-11-01_00.log", LastFile: "2009-11-10_21.log", Size: 1116, LastForward: "2009-11-10T21:58:00Z"},
		{Fullpath: "tmp/init/pod/namespace01", Subpath: "pod/namespace01", LogType: "pod", Target: "namespace01", FileCount: 4, FirstFile: "2000-01-01_00.log", LastFile: "2029-11-10_23.log", Size: 2620, LastForward: "2009-11-10T23:00:00Z"},
		{Fullpath: "tmp/init/pod/namespace02", Subpath: "pod/namespace02", LogType: "pod", Target: "namespace02", FileCount: 2, FirstFile: "0000-00-00_00.log", LastFile: "2009-11-10_22.log", Size: 1137, LastForward: "2009-11-10T22:58:00Z"}}
	got := fileService.ListTargets()
	assert.Equal(t, want, got)
}
