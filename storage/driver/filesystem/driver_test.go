package filesystem

import (
	"fmt"
	"sort"
	"testing"

	// storagedriver "github.com/kuoss/lethe/storage/driver"

	"github.com/kuoss/common/tester"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		params Params
		want   *driver
	}{
		{
			Params{},
			&driver{rootDirectory: ""},
		},
		{
			Params{RootDirectory: "asdf"},
			&driver{rootDirectory: "asdf"},
		},
		{
			Params{RootDirectory: "/data/log"},
			&driver{rootDirectory: "/data/log"},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got := New(tc.params)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestGetContent(t *testing.T) {
	testCases := []struct {
		path      string
		want      string
		wantError bool
	}{
		{
			path:      "",
			want:      "",
			wantError: true,
		},
		{
			path:      "node",
			want:      "",
			wantError: true,
		},
		{
			path: "pod/namespace01/2009-11-10_21.log",
			want: "2009-11-10T21:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\n2009-11-10T21:01:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\n2009-11-10T21:02:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\n",
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i), func(t *testing.T) {
			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			d := New(Params{RootDirectory: "data/log"})
			got, err := d.GetContent(tc.path)
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.want, string(got))
		})
	}
}

func TestPutContent(t *testing.T) {
	testCases := []struct {
		path      string
		content   string
		wantError bool
		want      string
	}{
		{
			path:      "",
			content:   "",
			wantError: true,
			want:      "",
		},
		{
			path:      "node",
			content:   "",
			wantError: true,
			want:      "",
		},
		{
			path:    "pod/namespace01/test1.log",
			content: "hello",
			want:    "hello",
		},
		{
			path:    "pod/namespace01/2009-11-10_21.log",
			content: "hello",
			want:    "hello11-10T21:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\n2009-11-10T21:01:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\n2009-11-10T21:02:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\n",
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i, tc.path), func(t *testing.T) {
			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			d := New(Params{RootDirectory: "data/log"})
			err := d.PutContent(tc.path, ([]byte)(tc.content))
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			got, err := d.GetContent(tc.path)
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.want, string(got))
		})
	}
}

func TestReader(t *testing.T) {
	testCases := []struct {
		path      string
		wantError string
	}{
		{
			"",
			"",
		},
		{
			"node",
			"",
		},
		{
			"pod/namespace01/2009-11-10_21.log",
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i, tc.path), func(t *testing.T) {
			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			d := New(Params{RootDirectory: "data/log"})
			got, err := d.Reader(tc.path)
			if tc.wantError == "" {
				require.NoError(t, err)
				require.NotEmpty(t, got)
			} else {
				require.EqualError(t, err, tc.wantError)
				require.Nil(t, got)
			}
		})
	}
}

func TestWriter(t *testing.T) {
	testCases := []struct {
		path      string
		wantError bool
		want      string
	}{
		{
			path:      "",
			wantError: true,
		},
		{
			path:      "node",
			wantError: true,
		},
		{
			path: "hello",
		},
		{
			path: "pod/namespace01/2009-11-10_21.log",
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i, tc.path), func(t *testing.T) {
			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			d := New(Params{RootDirectory: "data/log"})
			got, err := d.Writer(tc.path)
			if tc.wantError {
				require.Error(t, err)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, got)
			}
		})
	}
}

func TestStat(t *testing.T) {
	testCases := []struct {
		path      string
		want      string // fullpath
		wantError string
	}{
		{
			"",
			"data/log",
			"",
		},
		{
			"hello",
			"",
			"Path not found: hello, err: stat err: stat data/log/hello: no such file or directory",
		},
		{
			"node",
			"data/log/node",
			"",
		},
		{
			"pod",
			"data/log/pod",
			"",
		},
		{
			"pod/namespace01",
			"data/log/pod/namespace01",
			"",
		},
		{
			"pod/namespace01/hello.log",
			"",
			"Path not found: pod/namespace01/hello.log, err: stat err: stat data/log/pod/namespace01/hello.log: no such file or directory",
		},
		{
			"pod/namespace01/2009-11-10_21.log",
			"data/log/pod/namespace01/2009-11-10_21.log",
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i, tc.path), func(t *testing.T) {
			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			d := New(Params{RootDirectory: "data/log"})
			fi, err := d.Stat(tc.path)
			if tc.wantError != "" {
				require.EqualError(t, err, tc.wantError)
				return
			}
			require.NoError(t, err)
			got := fi.Fullpath()
			require.Equal(t, tc.want, got)
		})
	}
}

func TestList(t *testing.T) {
	testCases := []struct {
		path      string
		want      []string
		wantError string
	}{
		{
			"",
			[]string{"node", "pod"},
			"",
		},
		{
			"node",
			[]string{"node/node01", "node/node02"},
			"",
		},
		{
			"hello",
			nil,
			"Path not found: hello, err: open err: open data/log/hello: no such file or directory",
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
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i, tc.path), func(t *testing.T) {
			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			d := New(Params{RootDirectory: "data/log"})
			got, err := d.List(tc.path)
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

func TestMove(t *testing.T) {
	testCases := []struct {
		a         string
		b         string
		wantError string
	}{
		{
			"", "",
			"rename data/log data/log: file exists",
		},
		{
			"hello", "",
			"Path not found: hello, err: stat err: stat data/log/hello: no such file or directory",
		},
		{
			"", "hello",
			"rename data/log data/log/hello: invalid argument",
		},
		{
			"hello", "hello",
			"Path not found: hello, err: stat err: stat data/log/hello: no such file or directory",
		},
		{
			"pod/namespace01/hello.log", "pod/namespace01/hello.log",
			"Path not found: pod/namespace01/hello.log, err: stat err: stat data/log/pod/namespace01/hello.log: no such file or directory",
		},
		{
			"pod/namespace01/2009-11-10_21.log", "pod/namespace01/2009-11-10_00.log", // move
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i, tc.a, tc.b), func(t *testing.T) {
			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			d := New(Params{RootDirectory: "data/log"})
			err := d.Move(tc.a, tc.b)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		path      string
		wantError string
	}{
		{
			"",
			"deleting 0-1 depth directory is not allowed",
		},
		{
			"node",
			"deleting 0-1 depth directory is not allowed",
		},
		{
			"pod/namespace02",
			"",
		},
		{
			"pod/namespace01/2009-11-10_21.log", // delete
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i, tc.path), func(t *testing.T) {
			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			d := New(Params{RootDirectory: "data/log"})
			err := d.Delete(tc.path)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
		})
	}
}

func TestWalk(t *testing.T) {
	testCases := []struct {
		path      string
		want      []string
		wantError string
	}{
		// error
		{
			"hello",
			nil,
			"walk err: walkFunc err: lstat data/log/hello: no such file or directory",
		},
		{
			"data/log",
			nil,
			"walk err: walkFunc err: lstat data/log/data/log: no such file or directory",
		},
		// ok
		{
			"",
			[]string{"data/log/node/node01/2009-11-10_21.log", "data/log/node/node01/2009-11-10_22.log", "data/log/node/node02/2009-11-01_00.log", "data/log/node/node02/2009-11-10_21.log", "data/log/pod/namespace01/2000-01-01_00.log", "data/log/pod/namespace01/2009-11-10_21.log", "data/log/pod/namespace01/2009-11-10_22.log", "data/log/pod/namespace01/2029-11-10_23.log", "data/log/pod/namespace02/0000-00-00_00.log", "data/log/pod/namespace02/2009-11-10_22.log"},
			"",
		},
		{
			"node",
			[]string{"data/log/node/node01/2009-11-10_21.log", "data/log/node/node01/2009-11-10_22.log", "data/log/node/node02/2009-11-01_00.log", "data/log/node/node02/2009-11-10_21.log"},
			"",
		},
		{
			"pod",
			[]string{"data/log/pod/namespace01/2000-01-01_00.log", "data/log/pod/namespace01/2009-11-10_21.log", "data/log/pod/namespace01/2009-11-10_22.log", "data/log/pod/namespace01/2029-11-10_23.log", "data/log/pod/namespace02/0000-00-00_00.log", "data/log/pod/namespace02/2009-11-10_22.log"},
			"",
		},
		{
			"pod/namespace01",
			[]string{"data/log/pod/namespace01/2000-01-01_00.log", "data/log/pod/namespace01/2009-11-10_21.log", "data/log/pod/namespace01/2009-11-10_22.log", "data/log/pod/namespace01/2029-11-10_23.log"},
			"",
		},
		{
			"pod/namespace01/2009-11-10_21.log",
			[]string{"data/log/pod/namespace01/2009-11-10_21.log"},
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i, tc.path), func(t *testing.T) {
			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			d := New(Params{RootDirectory: "data/log"})
			infos, err := d.Walk(tc.path)
			if tc.wantError == "" {
				require.NoError(t, err)
				paths := []string{}
				for _, info := range infos {
					paths = append(paths, info.Fullpath())
				}
				require.Equal(t, tc.want, paths)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
		})
	}
}

func TestWalkDir(t *testing.T) {
	testCases := []struct {
		path      string
		want      []string
		wantError string
	}{
		// error
		{
			"hello",
			nil,
			"walkdir err: walkDirFunc err: lstat data/log/hello: no such file or directory",
		},
		{
			"data/log",
			nil,
			"walkdir err: walkDirFunc err: lstat data/log/data/log: no such file or directory",
		},
		// ok
		{
			"",
			[]string{".", "node", "node/node01", "node/node02", "pod", "pod/namespace01", "pod/namespace02"},
			"",
		},
		{
			"node",
			[]string{".", "node01", "node02"},
			"",
		},
		{
			"pod",
			[]string{".", "namespace01", "namespace02"},
			"",
		},
		{
			"pod/namespace01",
			[]string{"."},
			"",
		},
		{
			"pod/namespace01/2009-11-10_21.log",
			[]string{},
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i, tc.path), func(t *testing.T) {
			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			d := New(Params{RootDirectory: "data/log"})
			got, err := d.WalkDir(tc.path)
			if tc.wantError == "" {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
		})
	}
}
