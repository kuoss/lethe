package filesystem

import (
	"fmt"
	"sort"
	"testing"

	storagedriver "github.com/kuoss/lethe/storage/driver"
	"github.com/kuoss/lethe/util/testutil"
	"github.com/stretchr/testify/assert"
)

var (
	driver1            storagedriver.Driver
	logDataPath_driver = "tmp/storage_driver_filesystem_driver_test"
)

func init() {
	testutil.ChdirRoot()
	testutil.ResetLogData()
	driver1 = New(Params{RootDirectory: logDataPath_driver})
}

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
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestRootDirectory(t *testing.T) {
	got := driver1.RootDirectory()
	assert.Equal(t, logDataPath_driver, got)
}

func TestName(t *testing.T) {
	got := driver1.Name()
	assert.Equal(t, "filesystem", got)
}

func TestGetContent(t *testing.T) {
	testCases := []struct {
		path      string
		want      string
		wantError string
	}{
		{
			"",
			"",
			"read tmp/storage_driver_filesystem_driver_test: is a directory",
		},
		{
			"node",
			"",
			"read tmp/storage_driver_filesystem_driver_test/node: is a directory",
		},
		{
			"pod/namespace01/2009-11-10_21.log",
			"2009-11-10T21:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\n2009-11-10T21:01:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\n2009-11-10T21:02:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\n",
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got, err := driver1.GetContent(tc.path)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			assert.Equal(t, tc.want, string(got))
		})
	}
}

func TestPutContent(t *testing.T) {
	testCases := []struct {
		path       string
		content    string
		wantError  string
		wantError2 string
		want       string
	}{
		{
			"", "",
			"open tmp/storage_driver_filesystem_driver_test: is a directory",
			"",
			"",
		},
		{
			"node", "",
			"open tmp/storage_driver_filesystem_driver_test/node: is a directory",
			"",
			"",
		},
		{
			"pod/namespace01/test1.log", "hello",
			"",
			"",
			"hello",
		},
		{
			"pod/namespace01/2009-11-10_21.log", "hello",
			"",
			"",
			"hello11-10T21:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\n2009-11-10T21:01:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\n2009-11-10T21:02:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world\n",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			err := driver1.PutContent(tc.path, ([]byte)(tc.content))
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
				return
			}
			got, err := driver1.GetContent(tc.path)
			if tc.wantError2 == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError2)
				return
			}
			assert.Equal(t, tc.want, string(got))
		})
	}
	testutil.ResetLogData()
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
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got, err := driver1.Reader(tc.path)
			if tc.wantError == "" {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)
			} else {
				assert.EqualError(t, err, tc.wantError)
				assert.Nil(t, got)
			}
		})
	}
}

func TestWriter(t *testing.T) {
	testCases := []struct {
		path      string
		wantError string
	}{
		{
			"",
			"open tmp/storage_driver_filesystem_driver_test: is a directory",
		},
		{
			"hello",
			"",
		},
		{
			"node",
			"open tmp/storage_driver_filesystem_driver_test/node: is a directory",
		},
		{
			"pod/namespace01/2009-11-10_21.log",
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got, err := driver1.Writer(tc.path)
			if tc.wantError == "" {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)
			} else {
				assert.EqualError(t, err, tc.wantError)
				assert.Nil(t, got)
			}
		})
	}
	testutil.ResetLogData()
}

func TestStat(t *testing.T) {
	testCases := []struct {
		path      string
		want      string // fullpath
		wantError string
	}{
		{
			"",
			"tmp/storage_driver_filesystem_driver_test",
			"",
		},
		{
			"hello",
			"",
			"Path not found: hello, err: stat err: stat tmp/storage_driver_filesystem_driver_test/hello: no such file or directory",
		},
		{
			"node",
			"tmp/storage_driver_filesystem_driver_test/node",
			"",
		},
		{
			"pod",
			"tmp/storage_driver_filesystem_driver_test/pod",
			"",
		},
		{
			"pod/namespace01",
			"tmp/storage_driver_filesystem_driver_test/pod/namespace01",
			"",
		},
		{
			"pod/namespace01/hello.log",
			"",
			"Path not found: pod/namespace01/hello.log, err: stat err: stat tmp/storage_driver_filesystem_driver_test/pod/namespace01/hello.log: no such file or directory",
		},
		{
			"pod/namespace01/2009-11-10_21.log",
			"tmp/storage_driver_filesystem_driver_test/pod/namespace01/2009-11-10_21.log",
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			fi, err := driver1.Stat(tc.path)
			if tc.wantError != "" {
				assert.EqualError(t, err, tc.wantError)
				return
			}
			assert.NoError(t, err)
			got := fi.Fullpath()
			assert.Equal(t, tc.want, got)
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
			"Path not found: hello, err: open err: open tmp/storage_driver_filesystem_driver_test/hello: no such file or directory",
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
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got, err := driver1.List(tc.path)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			sort.Strings(got)
			assert.Equal(t, tc.want, got)
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
			"rename tmp/storage_driver_filesystem_driver_test tmp/storage_driver_filesystem_driver_test: file exists",
		},
		{
			"hello", "",
			"Path not found: hello, err: stat err: stat tmp/storage_driver_filesystem_driver_test/hello: no such file or directory",
		},
		{
			"", "hello",
			"rename tmp/storage_driver_filesystem_driver_test tmp/storage_driver_filesystem_driver_test/hello: invalid argument",
		},
		{
			"hello", "hello",
			"Path not found: hello, err: stat err: stat tmp/storage_driver_filesystem_driver_test/hello: no such file or directory",
		},
		{
			"pod/namespace01/hello.log", "pod/namespace01/hello.log",
			"Path not found: pod/namespace01/hello.log, err: stat err: stat tmp/storage_driver_filesystem_driver_test/pod/namespace01/hello.log: no such file or directory",
		},
		{
			"pod/namespace01/2009-11-10_21.log", "pod/namespace01/2009-11-10_00.log", // move
			"",
		},
		{
			"pod/namespace01/2009-11-10_21.log", "pod/namespace01/2009-11-10_00.log", // duplicate
			"Path not found: pod/namespace01/2009-11-10_21.log, err: stat err: stat tmp/storage_driver_filesystem_driver_test/pod/namespace01/2009-11-10_21.log: no such file or directory",
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			err := driver1.Move(tc.a, tc.b)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
		})
	}
	testutil.ResetLogData()
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
		{
			"pod/namespace01/2009-11-10_21.log", // duplicate
			"Path not found: pod/namespace01/2009-11-10_21.log, err: stat err: stat tmp/storage_driver_filesystem_driver_test/pod/namespace01/2009-11-10_21.log: no such file or directory",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			err := driver1.Delete(tc.path)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
		})
	}
	testutil.ResetLogData()
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
			"walk err: walkFunc err: lstat tmp/storage_driver_filesystem_driver_test/hello: no such file or directory",
		},
		{
			"tmp/storage_driver_filesystem_driver_test",
			nil,
			"walk err: walkFunc err: lstat tmp/storage_driver_filesystem_driver_test/tmp/storage_driver_filesystem_driver_test: no such file or directory",
		},
		// ok
		{
			"",
			[]string{"tmp/storage_driver_filesystem_driver_test/node/node01/2009-11-10_21.log", "tmp/storage_driver_filesystem_driver_test/node/node01/2009-11-10_22.log", "tmp/storage_driver_filesystem_driver_test/node/node02/2009-11-01_00.log", "tmp/storage_driver_filesystem_driver_test/node/node02/2009-11-10_21.log", "tmp/storage_driver_filesystem_driver_test/pod/namespace01/2000-01-01_00.log", "tmp/storage_driver_filesystem_driver_test/pod/namespace01/2009-11-10_21.log", "tmp/storage_driver_filesystem_driver_test/pod/namespace01/2009-11-10_22.log", "tmp/storage_driver_filesystem_driver_test/pod/namespace01/2029-11-10_23.log", "tmp/storage_driver_filesystem_driver_test/pod/namespace02/0000-00-00_00.log", "tmp/storage_driver_filesystem_driver_test/pod/namespace02/2009-11-10_22.log"},
			"",
		},
		{
			"node",
			[]string{"tmp/storage_driver_filesystem_driver_test/node/node01/2009-11-10_21.log", "tmp/storage_driver_filesystem_driver_test/node/node01/2009-11-10_22.log", "tmp/storage_driver_filesystem_driver_test/node/node02/2009-11-01_00.log", "tmp/storage_driver_filesystem_driver_test/node/node02/2009-11-10_21.log"},
			"",
		},
		{
			"pod",
			[]string{"tmp/storage_driver_filesystem_driver_test/pod/namespace01/2000-01-01_00.log", "tmp/storage_driver_filesystem_driver_test/pod/namespace01/2009-11-10_21.log", "tmp/storage_driver_filesystem_driver_test/pod/namespace01/2009-11-10_22.log", "tmp/storage_driver_filesystem_driver_test/pod/namespace01/2029-11-10_23.log", "tmp/storage_driver_filesystem_driver_test/pod/namespace02/0000-00-00_00.log", "tmp/storage_driver_filesystem_driver_test/pod/namespace02/2009-11-10_22.log"},
			"",
		},
		{
			"pod/namespace01",
			[]string{"tmp/storage_driver_filesystem_driver_test/pod/namespace01/2000-01-01_00.log", "tmp/storage_driver_filesystem_driver_test/pod/namespace01/2009-11-10_21.log", "tmp/storage_driver_filesystem_driver_test/pod/namespace01/2009-11-10_22.log", "tmp/storage_driver_filesystem_driver_test/pod/namespace01/2029-11-10_23.log"},
			"",
		},
		{
			"pod/namespace01/2009-11-10_21.log",
			[]string{"tmp/storage_driver_filesystem_driver_test/pod/namespace01/2009-11-10_21.log"},
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			infos, err := driver1.Walk(tc.path)
			if tc.wantError == "" {
				assert.NoError(t, err)
				paths := []string{}
				for _, info := range infos {
					paths = append(paths, info.Fullpath())
				}
				assert.Equal(t, tc.want, paths)
			} else {
				assert.EqualError(t, err, tc.wantError)
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
			"walkdir err: walkDirFunc err: lstat tmp/storage_driver_filesystem_driver_test/hello: no such file or directory",
		},
		{
			"tmp/storage_driver_filesystem_driver_test",
			nil,
			"walkdir err: walkDirFunc err: lstat tmp/storage_driver_filesystem_driver_test/tmp/storage_driver_filesystem_driver_test: no such file or directory",
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
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got, err := driver1.WalkDir(tc.path)
			if tc.wantError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
		})
	}
}
