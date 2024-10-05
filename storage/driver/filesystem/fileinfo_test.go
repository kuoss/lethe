package filesystem

import (
	"os"
	"testing"

	"github.com/kuoss/common/tester"
	"github.com/stretchr/testify/require"
)

func TestFileInfo(t *testing.T) {
	testCases := []struct {
		path            string
		wantError       string
		wantFullpath    string
		wantSize        int64
		wantModTimeYear int
		wantIsDir       bool
	}{
		{
			"", "stat : no such file or directory",
			"", 0, 2023, false,
		},
		{
			"hello", "stat hello: no such file or directory",
			"", 0, 2023, false,
		},
		{
			"data/log/pod", "",
			"data/log/pod", 0, 2023, true,
		},
		{
			"data/log/pod/namespace01", "",
			"data/log/pod/namespace01", 0, 2023, true,
		},
		{
			"data/log/pod/namespace01/2009-11-10_21.log", "",
			"data/log/pod/namespace01/2009-11-10_21.log", 279, 2023, false,
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i, tc.path), func(t *testing.T) {
			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			fi, err := os.Stat(tc.path)
			if tc.wantError != "" {
				require.EqualError(t, err, tc.wantError)
			} else {
				require.NoError(t, err)

				got := FileInfo{fi, tc.path}
				require.Equal(t, tc.wantFullpath, got.Fullpath())
				require.Equal(t, tc.wantSize, got.Size())
				require.GreaterOrEqual(t, got.ModTime().Year(), tc.wantModTimeYear)
				require.Equal(t, tc.wantIsDir, got.IsDir())
			}
		})
	}
}
