package rotator

import (
	"testing"
	"time"

	"github.com/kuoss/common/tester"
	"github.com/kuoss/lethe/clock"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/storage/fileservice"
	"github.com/stretchr/testify/require"
)

func TestRunOnce(t *testing.T) {
	testCases := []struct {
		retentionSize int64
		retentionTime time.Duration
		want          []string
	}{
		{
			retentionSize: 0,
			retentionTime: 0,
			want: []string{
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
			retentionSize: 10,              // 10 bytes
			retentionTime: 100 * time.Hour, // 100 hours
			want:          []string{},
		},
		{
			retentionSize: 2 * 1024,        // 2 KiB
			retentionTime: 100 * time.Hour, // 100 hours
			want:          []string{"pod/namespace01/2029-11-10_23.log", "pod/namespace02/2009-11-10_22.log"},
		},
		{
			retentionSize: 3 * 1024,        // 3 KiB
			retentionTime: 100 * time.Hour, // 100 hours
			want: []string{
				"node/node01/2009-11-10_22.log",
				"pod/namespace01/2009-11-10_21.log",
				"pod/namespace01/2009-11-10_22.log",
				"pod/namespace01/2029-11-10_23.log",
				"pod/namespace02/2009-11-10_22.log",
			},
		},
		{
			retentionSize: 100 * 1024 * 1024 * 1024, // 100 GiB
			retentionTime: 2 * 24 * time.Hour,       // 2 days
			want: []string{
				"node/node01/2009-11-10_21.log",
				"node/node01/2009-11-10_22.log",
				"node/node02/2009-11-10_21.log",
				"pod/namespace01/2009-11-10_21.log",
				"pod/namespace01/2009-11-10_22.log",
				"pod/namespace01/2029-11-10_23.log",
				"pod/namespace02/2009-11-10_22.log",
			},
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i, tc.retentionSize, tc.retentionTime), func(t *testing.T) {
			clock.SetPlaygroundMode(true)
			defer clock.SetPlaygroundMode(false)

			_, cleanup := tester.SetupDir(t, map[string]string{
				"@/testdata/log": "data/log",
			})
			defer cleanup()

			cfg, err := config.New("test")
			require.NoError(t, err)
			cfg.Retention.Size = tc.retentionSize
			cfg.Retention.Time = tc.retentionTime

			fileService, err := fileservice.New(cfg)
			require.NoError(t, err)
			rotator := New(cfg, fileService)

			rotator.RunOnce()

			list, err := rotator.fileService.ListFiles()
			subpaths := []string{}
			for _, item := range list {
				subpaths = append(subpaths, item.Subpath)
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, subpaths)
		})
	}
}
