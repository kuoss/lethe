package filesystem

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileInfo(t *testing.T) {
	testCases := []struct {
		fullpath        string
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
			"tmp/init/pod", "",
			"tmp/init/pod", 0, 2023, true,
		},
		{
			"tmp/init/pod/namespace01", "",
			"tmp/init/pod/namespace01", 0, 2023, true,
		},
		{
			"tmp/init/pod/namespace01/2009-11-10_21.log", "",
			"tmp/init/pod/namespace01/2009-11-10_21.log", 279, 2023, false,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			fi, err := os.Stat(tc.fullpath)
			if tc.wantError != "" {
				assert.EqualError(t, err, tc.wantError)
			} else {
				assert.NoError(t, err)

				got := FileInfo{fi, tc.fullpath}
				assert.Equal(t, tc.wantFullpath, got.Fullpath())
				assert.Equal(t, tc.wantSize, got.Size())
				assert.GreaterOrEqual(t, got.ModTime().Year(), tc.wantModTimeYear)
				assert.Equal(t, tc.wantIsDir, got.IsDir())
			}
		})
	}
}
