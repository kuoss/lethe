package rotator

import (
	"fmt"
	"testing"
	"time"

	"github.com/kuoss/lethe/clock"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/storage/fileservice"
	"github.com/kuoss/lethe/util/testutil"
	"github.com/stretchr/testify/assert"
)

var (
	rotator1 *Rotator
)

func init() {
	testutil.ChdirRoot()
	testutil.ResetLogData()
	clock.SetPlaygroundMode(true)

	cfg, err := config.New("test")
	if err != nil {
		panic(err)
	}
	cfg.SetLogDataPath("tmp/rotator_rotator_test")
	fileService, err := fileservice.New(cfg)
	if err != nil {
		panic(err)
	}
	rotator1 = New(cfg, fileService)
}

func TestRunOnce(t *testing.T) {
	testCases := []struct {
		retentionSize int
		retentionTime time.Duration
		want          []fileservice.LogFile
	}{
		{
			0,
			0,
			[]fileservice.LogFile{
				{Fullpath: "tmp/rotator_rotator_test/node/node01/2009-11-10_21.log", Subpath: "node/node01/2009-11-10_21.log", LogType: "node", Target: "node01", Name: "2009-11-10_21.log", Extension: ".log", Size: 1057},
				{Fullpath: "tmp/rotator_rotator_test/node/node01/2009-11-10_22.log", Subpath: "node/node01/2009-11-10_22.log", LogType: "node", Target: "node01", Name: "2009-11-10_22.log", Extension: ".log", Size: 177},
				{Fullpath: "tmp/rotator_rotator_test/node/node02/2009-11-01_00.log", Subpath: "node/node02/2009-11-01_00.log", LogType: "node", Target: "node02", Name: "2009-11-01_00.log", Extension: ".log", Size: 0},
				{Fullpath: "tmp/rotator_rotator_test/node/node02/2009-11-10_21.log", Subpath: "node/node02/2009-11-10_21.log", LogType: "node", Target: "node02", Name: "2009-11-10_21.log", Extension: ".log", Size: 1116},
				{Fullpath: "tmp/rotator_rotator_test/pod/namespace01/2000-01-01_00.log", Subpath: "pod/namespace01/2000-01-01_00.log", LogType: "pod", Target: "namespace01", Name: "2000-01-01_00.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/rotator_rotator_test/pod/namespace01/2009-11-10_21.log", Subpath: "pod/namespace01/2009-11-10_21.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_21.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/rotator_rotator_test/pod/namespace01/2009-11-10_22.log", Subpath: "pod/namespace01/2009-11-10_22.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_22.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/rotator_rotator_test/pod/namespace01/2029-11-10_23.log", Subpath: "pod/namespace01/2029-11-10_23.log", LogType: "pod", Target: "namespace01", Name: "2029-11-10_23.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/rotator_rotator_test/pod/namespace02/0000-00-00_00.log", Subpath: "pod/namespace02/0000-00-00_00.log", LogType: "pod", Target: "namespace02", Name: "0000-00-00_00.log", Extension: ".log", Size: 12},
				{Fullpath: "tmp/rotator_rotator_test/pod/namespace02/2009-11-10_22.log", Subpath: "pod/namespace02/2009-11-10_22.log", LogType: "pod", Target: "namespace02", Name: "2009-11-10_22.log", Extension: ".log", Size: 1125}},
		},
		{
			10,              // 10 bytes
			100 * time.Hour, // 100 hours
			[]fileservice.LogFile{},
		},
		{
			2 * 1024,        // 2 KiB
			100 * time.Hour, // 100 hours
			[]fileservice.LogFile{
				{Fullpath: "tmp/rotator_rotator_test/pod/namespace01/2029-11-10_23.log", Subpath: "pod/namespace01/2029-11-10_23.log", LogType: "pod", Target: "namespace01", Name: "2029-11-10_23.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/rotator_rotator_test/pod/namespace02/2009-11-10_22.log", Subpath: "pod/namespace02/2009-11-10_22.log", LogType: "pod", Target: "namespace02", Name: "2009-11-10_22.log", Extension: ".log", Size: 1125}},
		},
		{
			3 * 1024,        // 3 KiB
			100 * time.Hour, // 100 hours
			[]fileservice.LogFile{
				{Fullpath: "tmp/rotator_rotator_test/node/node01/2009-11-10_22.log", Subpath: "node/node01/2009-11-10_22.log", LogType: "node", Target: "node01", Name: "2009-11-10_22.log", Extension: ".log", Size: 177},
				{Fullpath: "tmp/rotator_rotator_test/pod/namespace01/2009-11-10_21.log", Subpath: "pod/namespace01/2009-11-10_21.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_21.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/rotator_rotator_test/pod/namespace01/2009-11-10_22.log", Subpath: "pod/namespace01/2009-11-10_22.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_22.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/rotator_rotator_test/pod/namespace01/2029-11-10_23.log", Subpath: "pod/namespace01/2029-11-10_23.log", LogType: "pod", Target: "namespace01", Name: "2029-11-10_23.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/rotator_rotator_test/pod/namespace02/2009-11-10_22.log", Subpath: "pod/namespace02/2009-11-10_22.log", LogType: "pod", Target: "namespace02", Name: "2009-11-10_22.log", Extension: ".log", Size: 1125}},
		},
		{
			100 * 1024 * 1024 * 1024, // 100 GiB
			2 * 24 * time.Hour,       // 2 days
			[]fileservice.LogFile{
				{Fullpath: "tmp/rotator_rotator_test/node/node01/2009-11-10_21.log", Subpath: "node/node01/2009-11-10_21.log", LogType: "node", Target: "node01", Name: "2009-11-10_21.log", Extension: ".log", Size: 1057},
				{Fullpath: "tmp/rotator_rotator_test/node/node01/2009-11-10_22.log", Subpath: "node/node01/2009-11-10_22.log", LogType: "node", Target: "node01", Name: "2009-11-10_22.log", Extension: ".log", Size: 177},
				{Fullpath: "tmp/rotator_rotator_test/node/node02/2009-11-10_21.log", Subpath: "node/node02/2009-11-10_21.log", LogType: "node", Target: "node02", Name: "2009-11-10_21.log", Extension: ".log", Size: 1116},
				{Fullpath: "tmp/rotator_rotator_test/pod/namespace01/2009-11-10_21.log", Subpath: "pod/namespace01/2009-11-10_21.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_21.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/rotator_rotator_test/pod/namespace01/2009-11-10_22.log", Subpath: "pod/namespace01/2009-11-10_22.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_22.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/rotator_rotator_test/pod/namespace01/2029-11-10_23.log", Subpath: "pod/namespace01/2029-11-10_23.log", LogType: "pod", Target: "namespace01", Name: "2029-11-10_23.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/rotator_rotator_test/pod/namespace02/2009-11-10_22.log", Subpath: "pod/namespace02/2009-11-10_22.log", LogType: "pod", Target: "namespace02", Name: "2009-11-10_22.log", Extension: ".log", Size: 1125}},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			rotator1.config.SetRetentionSize(tc.retentionSize)
			rotator1.config.SetRetentionTime(tc.retentionTime)
			rotator1.RunOnce()
			got, err := rotator1.fileService.ListFiles()
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
			testutil.ResetLogData()
		})
	}
}
