package fileservice

import (
	"sort"
	"testing"

	"github.com/kuoss/common/tester"
	"github.com/kuoss/lethe/config"
	"github.com/stretchr/testify/require"
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
		{"data/log", "data/log/pod", "pod"},
		{"data/log", "data/log/pod/ns1", "pod/ns1"},
		{"data/log", "data/log/pod/ns1/2023", "pod/ns1/2023"},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i), func(t *testing.T) {
			got := fullpath2subpath(tc.rootDir, tc.fullpath)
			require.Equal(t, tc.want, got)
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
		{"hello", 0, "Path not found: hello, err: open err: open data/log/hello: no such file or directory"},
		{"node", 0, ""},
		{"pod", 0, ""},
		{"node/node01", 1234, ""},
		{"node/node02", 1116, ""},
		{"pod/namespace01", 2620, ""},
		{"pod/namespace02", 1137, ""},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i, tc.path), func(t *testing.T) {
			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			cfg, err := config.New("test")
			require.NoError(t, err)
			fileService, err := New(cfg)
			require.NoError(t, err)

			got, err := fileService.dirSize(tc.path)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			require.Equal(t, tc.want, got)
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
			[]string{"node", "pod"},
			"",
		},
		{
			".",
			[]string{"node", "pod"},
			"",
		},
		{
			"hello",
			nil,
			"list err: Path not found: hello, err: open err: open data/log/hello: no such file or directory",
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
			[]string{"pod/namespace01/2000-01-01_00.log", "pod/namespace01/2009-11-10_21.log", "pod/namespace01/2009-11-10_22.log", "pod/namespace01/2029-11-10_23.log"},
			"",
		},
		{
			"pod/namespace01/2029-11-10_23.log",
			nil,
			"list err: readdirnames err: readdirent data/log/pod/namespace01/2029-11-10_23.log: not a directory",
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			cfg, err := config.New("test")
			require.NoError(t, err)
			fileService, err := New(cfg)
			require.NoError(t, err)

			got, err := fileService.List(tc.subpath)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			sort.Strings(got)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestListLogDirs(t *testing.T) {
	_, cleanup := tester.SetupDir(t, map[string]string{
		"@/testdata/log": "data/log",
	})
	defer cleanup()

	cfg, err := config.New("test")
	require.NoError(t, err)
	fileService, err := New(cfg)
	require.NoError(t, err)

	want := []LogDir{
		{Fullpath: "data/log/node/node01", Subpath: "node/node01", LogType: "node", Target: "node01", FileCount: 0, FirstFile: "", LastFile: "", Size: 0, LastForward: ""},
		{Fullpath: "data/log/node/node02", Subpath: "node/node02", LogType: "node", Target: "node02", FileCount: 0, FirstFile: "", LastFile: "", Size: 0, LastForward: ""},
		{Fullpath: "data/log/pod/namespace01", Subpath: "pod/namespace01", LogType: "pod", Target: "namespace01", FileCount: 0, FirstFile: "", LastFile: "", Size: 0, LastForward: ""},
		{Fullpath: "data/log/pod/namespace02", Subpath: "pod/namespace02", LogType: "pod", Target: "namespace02", FileCount: 0, FirstFile: "", LastFile: "", Size: 0, LastForward: ""}}
	got := fileService.ListLogDirs()
	require.Equal(t, want, got)
}

func TestListLogDirsWithSize(t *testing.T) {
	_, cleanup := tester.SetupDir(t, map[string]string{
		"@/testdata/log": "data/log",
	})
	defer cleanup()

	cfg, err := config.New("test")
	require.NoError(t, err)
	fileService, err := New(cfg)
	require.NoError(t, err)

	want := []LogDir{
		{Fullpath: "data/log/node/node01", Subpath: "node/node01", LogType: "node", Target: "node01", FileCount: 2, FirstFile: "2009-11-10_21.log", LastFile: "2009-11-10_22.log", Size: 1234, LastForward: ""},
		{Fullpath: "data/log/node/node02", Subpath: "node/node02", LogType: "node", Target: "node02", FileCount: 2, FirstFile: "2009-11-01_00.log", LastFile: "2009-11-10_21.log", Size: 1116, LastForward: ""},
		{Fullpath: "data/log/pod/namespace01", Subpath: "pod/namespace01", LogType: "pod", Target: "namespace01", FileCount: 4, FirstFile: "2000-01-01_00.log", LastFile: "2029-11-10_23.log", Size: 2620, LastForward: ""},
		{Fullpath: "data/log/pod/namespace02", Subpath: "pod/namespace02", LogType: "pod", Target: "namespace02", FileCount: 2, FirstFile: "0000-00-00_00.log", LastFile: "2009-11-10_22.log", Size: 1137, LastForward: ""}}
	got := fileService.ListLogDirsWithSize()
	require.Equal(t, want, got)
}

func TestListTargets(t *testing.T) {
	_, cleanup := tester.SetupDir(t, map[string]string{
		"@/testdata/log": "data/log",
	})
	defer cleanup()

	cfg, err := config.New("test")
	require.NoError(t, err)
	fileService, err := New(cfg)
	require.NoError(t, err)

	want := []LogDir{
		{Fullpath: "data/log/node/node01", Subpath: "node/node01", LogType: "node", Target: "node01", FileCount: 2, FirstFile: "2009-11-10_21.log", LastFile: "2009-11-10_22.log", Size: 1234, LastForward: "2009-11-10T23:00:00Z"},
		{Fullpath: "data/log/node/node02", Subpath: "node/node02", LogType: "node", Target: "node02", FileCount: 2, FirstFile: "2009-11-01_00.log", LastFile: "2009-11-10_21.log", Size: 1116, LastForward: "2009-11-10T21:58:00Z"},
		{Fullpath: "data/log/pod/namespace01", Subpath: "pod/namespace01", LogType: "pod", Target: "namespace01", FileCount: 4, FirstFile: "2000-01-01_00.log", LastFile: "2029-11-10_23.log", Size: 2620, LastForward: "2009-11-10T23:00:00Z"},
		{Fullpath: "data/log/pod/namespace02", Subpath: "pod/namespace02", LogType: "pod", Target: "namespace02", FileCount: 2, FirstFile: "0000-00-00_00.log", LastFile: "2009-11-10_22.log", Size: 1137, LastForward: "2009-11-10T22:58:00Z"}}
	got := fileService.ListTargets()
	require.Equal(t, want, got)
}
