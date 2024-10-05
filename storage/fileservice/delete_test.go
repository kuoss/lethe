package fileservice

import (
	"testing"
	"time"

	"github.com/kuoss/common/tester"
	"github.com/kuoss/lethe/clock"
	"github.com/kuoss/lethe/config"
	"github.com/stretchr/testify/require"
)

func TestDeleteByAge(t *testing.T) {
	testCases := []struct {
		retentionTime time.Duration
		want          []string
	}{
		{
			0, // disabled
			[]string{
				"node/node01/2009-11-10_21.log",
				"node/node01/2009-11-10_22.log",
				"node/node02/2009-11-01_00.log",
				"node/node02/2009-11-10_21.log",
				"pod/namespace01/2000-01-01_00.log",
				"pod/namespace01/2009-11-10_21.log",
				"pod/namespace01/2009-11-10_22.log",
				"pod/namespace01/2029-11-10_23.log",
				"pod/namespace02/0000-00-00_00.log",
				"pod/namespace02/2009-11-10_22.log",
			},
		},
		{
			1 * time.Second, // 1s
			[]string{
				"node/node01/2009-11-10_22.log",
				"pod/namespace01/2009-11-10_22.log",
				"pod/namespace01/2029-11-10_23.log",
				"pod/namespace02/2009-11-10_22.log",
			},
		},
		{
			1 * time.Hour, // 1h
			[]string{
				"node/node01/2009-11-10_22.log",
				"pod/namespace01/2009-11-10_22.log",
				"pod/namespace01/2029-11-10_23.log",
				"pod/namespace02/2009-11-10_22.log",
			},
		},
		{
			2 * time.Hour, // 2h
			[]string{
				"node/node01/2009-11-10_21.log",
				"node/node01/2009-11-10_22.log",
				"node/node02/2009-11-10_21.log",
				"pod/namespace01/2009-11-10_21.log",
				"pod/namespace01/2009-11-10_22.log",
				"pod/namespace01/2029-11-10_23.log",
				"pod/namespace02/2009-11-10_22.log",
			},
		},
		{
			1 * 24 * time.Hour, // 1d
			[]string{
				"node/node01/2009-11-10_21.log",
				"node/node01/2009-11-10_22.log",
				"node/node02/2009-11-10_21.log",
				"pod/namespace01/2009-11-10_21.log",
				"pod/namespace01/2009-11-10_22.log",
				"pod/namespace01/2029-11-10_23.log",
				"pod/namespace02/2009-11-10_22.log",
			},
		},
		{
			10 * 24 * time.Hour, // 10d
			[]string{
				"node/node01/2009-11-10_21.log",
				"node/node01/2009-11-10_22.log",
				"node/node02/2009-11-01_00.log",
				"node/node02/2009-11-10_21.log",
				"pod/namespace01/2009-11-10_21.log",
				"pod/namespace01/2009-11-10_22.log",
				"pod/namespace01/2029-11-10_23.log",
				"pod/namespace02/2009-11-10_22.log",
			},
		},
		{
			100 * 365 * 24 * time.Hour, // 100y
			[]string{
				"node/node01/2009-11-10_21.log",
				"node/node01/2009-11-10_22.log",
				"node/node02/2009-11-01_00.log",
				"node/node02/2009-11-10_21.log",
				"pod/namespace01/2000-01-01_00.log",
				"pod/namespace01/2009-11-10_21.log",
				"pod/namespace01/2009-11-10_22.log",
				"pod/namespace01/2029-11-10_23.log",
				"pod/namespace02/2009-11-10_22.log",
			},
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i, tc.retentionTime), func(t *testing.T) {
			clock.SetPlaygroundMode(true)
			defer clock.SetPlaygroundMode(false)

			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			cfg, err := config.New("test")
			require.NoError(t, err)
			cfg.Retention.Time = tc.retentionTime

			fileService, err := New(cfg)
			require.NoError(t, err)

			err = fileService.DeleteByAge()
			require.NoError(t, err)

			listFiles, err := fileService.ListFiles()
			require.NoError(t, err)

			subPaths := []string{}
			for _, f := range listFiles {
				subPaths = append(subPaths, f.Subpath)
			}
			require.Equal(t, tc.want, subPaths)
		})
	}
}

func TestDeleteBySize(t *testing.T) {
	testCases := []struct {
		retentionSize int64
		want          []string
	}{
		{
			0, // disabled
			[]string{
				"node/node01/2009-11-10_21.log",
				"node/node01/2009-11-10_22.log",
				"node/node02/2009-11-01_00.log",
				"node/node02/2009-11-10_21.log",
				"pod/namespace01/2000-01-01_00.log",
				"pod/namespace01/2009-11-10_21.log",
				"pod/namespace01/2009-11-10_22.log",
				"pod/namespace01/2029-11-10_23.log",
				"pod/namespace02/0000-00-00_00.log",
				"pod/namespace02/2009-11-10_22.log",
			},
		},
		{
			10, // 10 bytes
			[]string{},
		},
		{
			1 * 1024, // 1 KiB
			[]string{"pod/namespace01/2029-11-10_23.log"},
		},
		{
			2 * 1024, // 2 KiB
			[]string{
				"pod/namespace01/2029-11-10_23.log",
				"pod/namespace02/2009-11-10_22.log",
			},
		},
		{
			3 * 1024, // 3 KiB
			[]string{
				"node/node01/2009-11-10_22.log",
				"pod/namespace01/2009-11-10_21.log",
				"pod/namespace01/2009-11-10_22.log",
				"pod/namespace01/2029-11-10_23.log",
				"pod/namespace02/2009-11-10_22.log",
			},
		},
		{
			1 * 1024 * 1024, // 1 MiB
			[]string{
				"node/node01/2009-11-10_21.log",
				"node/node01/2009-11-10_22.log",
				"node/node02/2009-11-01_00.log",
				"node/node02/2009-11-10_21.log",
				"pod/namespace01/2000-01-01_00.log",
				"pod/namespace01/2009-11-10_21.log",
				"pod/namespace01/2009-11-10_22.log",
				"pod/namespace01/2029-11-10_23.log",
				"pod/namespace02/0000-00-00_00.log",
				"pod/namespace02/2009-11-10_22.log",
			},
		},
		{
			999999999 * 1024 * 1024 * 1024, // 999,999,999 GiB
			[]string{
				"node/node01/2009-11-10_21.log",
				"node/node01/2009-11-10_22.log",
				"node/node02/2009-11-01_00.log",
				"node/node02/2009-11-10_21.log",
				"pod/namespace01/2000-01-01_00.log",
				"pod/namespace01/2009-11-10_21.log",
				"pod/namespace01/2009-11-10_22.log",
				"pod/namespace01/2029-11-10_23.log",
				"pod/namespace02/0000-00-00_00.log",
				"pod/namespace02/2009-11-10_22.log",
			},
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i, tc.retentionSize), func(t *testing.T) {
			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			cfg, err := config.New("test")
			require.NoError(t, err)
			cfg.Retention.Size = tc.retentionSize

			fileService, err := New(cfg)
			require.NoError(t, err)

			err = fileService.DeleteBySize()
			require.NoError(t, err)

			listFiles, err := fileService.ListFiles()
			require.NoError(t, err)
			subPaths := []string{}
			for _, f := range listFiles {
				subPaths = append(subPaths, f.Subpath)
			}
			require.Equal(t, tc.want, subPaths)
		})
	}
}
