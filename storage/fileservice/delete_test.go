package fileservice

import (
	"fmt"
	"testing"
	"time"

	"github.com/kuoss/lethe/clock"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/util/testutil"
	"github.com/stretchr/testify/assert"
)

var (
	fileService2 *FileService
)

func init() {
	testutil.ChdirRoot()
	testutil.ResetLogData()
	clock.SetPlaygroundMode(true)

	cfg, err := config.New("test")
	if err != nil {
		panic(err)
	}
	cfg.SetLogDataPath("tmp/storage_fileservice_delete_test")
	fileService2, err = New(cfg)
	if err != nil {
		panic(err)
	}
}

func TestDeleteByAge(t *testing.T) {
	testCases := []struct {
		retentionTime time.Duration
		want          []LogFile
	}{
		{
			0, // disabled
			[]LogFile{
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_21.log", Subpath: "node/node01/2009-11-10_21.log", LogType: "node", Target: "node01", Name: "2009-11-10_21.log", Extension: ".log", Size: 1057},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_22.log", Subpath: "node/node01/2009-11-10_22.log", LogType: "node", Target: "node01", Name: "2009-11-10_22.log", Extension: ".log", Size: 177},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node02/2009-11-01_00.log", Subpath: "node/node02/2009-11-01_00.log", LogType: "node", Target: "node02", Name: "2009-11-01_00.log", Extension: ".log", Size: 0},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node02/2009-11-10_21.log", Subpath: "node/node02/2009-11-10_21.log", LogType: "node", Target: "node02", Name: "2009-11-10_21.log", Extension: ".log", Size: 1116},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2000-01-01_00.log", Subpath: "pod/namespace01/2000-01-01_00.log", LogType: "pod", Target: "namespace01", Name: "2000-01-01_00.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_21.log", Subpath: "pod/namespace01/2009-11-10_21.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_21.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_22.log", Subpath: "pod/namespace01/2009-11-10_22.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_22.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2029-11-10_23.log", Subpath: "pod/namespace01/2029-11-10_23.log", LogType: "pod", Target: "namespace01", Name: "2029-11-10_23.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace02/0000-00-00_00.log", Subpath: "pod/namespace02/0000-00-00_00.log", LogType: "pod", Target: "namespace02", Name: "0000-00-00_00.log", Extension: ".log", Size: 12},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace02/2009-11-10_22.log", Subpath: "pod/namespace02/2009-11-10_22.log", LogType: "pod", Target: "namespace02", Name: "2009-11-10_22.log", Extension: ".log", Size: 1125}},
		},
		{
			1 * time.Second, // 1 second
			[]LogFile{
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_22.log", Subpath: "node/node01/2009-11-10_22.log", LogType: "node", Target: "node01", Name: "2009-11-10_22.log", Extension: ".log", Size: 177},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_22.log", Subpath: "pod/namespace01/2009-11-10_22.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_22.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2029-11-10_23.log", Subpath: "pod/namespace01/2029-11-10_23.log", LogType: "pod", Target: "namespace01", Name: "2029-11-10_23.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace02/2009-11-10_22.log", Subpath: "pod/namespace02/2009-11-10_22.log", LogType: "pod", Target: "namespace02", Name: "2009-11-10_22.log", Extension: ".log", Size: 1125}},
		},
		{
			1 * time.Hour, // 1 hour
			[]LogFile{
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_22.log", Subpath: "node/node01/2009-11-10_22.log", LogType: "node", Target: "node01", Name: "2009-11-10_22.log", Extension: ".log", Size: 177},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_22.log", Subpath: "pod/namespace01/2009-11-10_22.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_22.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2029-11-10_23.log", Subpath: "pod/namespace01/2029-11-10_23.log", LogType: "pod", Target: "namespace01", Name: "2029-11-10_23.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace02/2009-11-10_22.log", Subpath: "pod/namespace02/2009-11-10_22.log", LogType: "pod", Target: "namespace02", Name: "2009-11-10_22.log", Extension: ".log", Size: 1125}},
		},
		{
			2 * time.Hour, // 2 hours
			[]LogFile{
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_21.log", Subpath: "node/node01/2009-11-10_21.log", LogType: "node", Target: "node01", Name: "2009-11-10_21.log", Extension: ".log", Size: 1057},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_22.log", Subpath: "node/node01/2009-11-10_22.log", LogType: "node", Target: "node01", Name: "2009-11-10_22.log", Extension: ".log", Size: 177},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node02/2009-11-10_21.log", Subpath: "node/node02/2009-11-10_21.log", LogType: "node", Target: "node02", Name: "2009-11-10_21.log", Extension: ".log", Size: 1116},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_21.log", Subpath: "pod/namespace01/2009-11-10_21.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_21.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_22.log", Subpath: "pod/namespace01/2009-11-10_22.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_22.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2029-11-10_23.log", Subpath: "pod/namespace01/2029-11-10_23.log", LogType: "pod", Target: "namespace01", Name: "2029-11-10_23.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace02/2009-11-10_22.log", Subpath: "pod/namespace02/2009-11-10_22.log", LogType: "pod", Target: "namespace02", Name: "2009-11-10_22.log", Extension: ".log", Size: 1125}},
		},
		{
			1 * 24 * time.Hour, // 1 day
			[]LogFile{
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_21.log", Subpath: "node/node01/2009-11-10_21.log", LogType: "node", Target: "node01", Name: "2009-11-10_21.log", Extension: ".log", Size: 1057},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_22.log", Subpath: "node/node01/2009-11-10_22.log", LogType: "node", Target: "node01", Name: "2009-11-10_22.log", Extension: ".log", Size: 177},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node02/2009-11-10_21.log", Subpath: "node/node02/2009-11-10_21.log", LogType: "node", Target: "node02", Name: "2009-11-10_21.log", Extension: ".log", Size: 1116},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_21.log", Subpath: "pod/namespace01/2009-11-10_21.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_21.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_22.log", Subpath: "pod/namespace01/2009-11-10_22.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_22.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2029-11-10_23.log", Subpath: "pod/namespace01/2029-11-10_23.log", LogType: "pod", Target: "namespace01", Name: "2029-11-10_23.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace02/2009-11-10_22.log", Subpath: "pod/namespace02/2009-11-10_22.log", LogType: "pod", Target: "namespace02", Name: "2009-11-10_22.log", Extension: ".log", Size: 1125}},
		},
		{
			10 * 24 * time.Hour, // 10 days
			[]LogFile{
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_21.log", Subpath: "node/node01/2009-11-10_21.log", LogType: "node", Target: "node01", Name: "2009-11-10_21.log", Extension: ".log", Size: 1057},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_22.log", Subpath: "node/node01/2009-11-10_22.log", LogType: "node", Target: "node01", Name: "2009-11-10_22.log", Extension: ".log", Size: 177},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node02/2009-11-01_00.log", Subpath: "node/node02/2009-11-01_00.log", LogType: "node", Target: "node02", Name: "2009-11-01_00.log", Extension: ".log", Size: 0},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node02/2009-11-10_21.log", Subpath: "node/node02/2009-11-10_21.log", LogType: "node", Target: "node02", Name: "2009-11-10_21.log", Extension: ".log", Size: 1116},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_21.log", Subpath: "pod/namespace01/2009-11-10_21.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_21.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_22.log", Subpath: "pod/namespace01/2009-11-10_22.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_22.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2029-11-10_23.log", Subpath: "pod/namespace01/2029-11-10_23.log", LogType: "pod", Target: "namespace01", Name: "2029-11-10_23.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace02/2009-11-10_22.log", Subpath: "pod/namespace02/2009-11-10_22.log", LogType: "pod", Target: "namespace02", Name: "2009-11-10_22.log", Extension: ".log", Size: 1125}},
		},
		{
			100 * 365 * 24 * time.Hour, // 100 years
			[]LogFile{
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_21.log", Subpath: "node/node01/2009-11-10_21.log", LogType: "node", Target: "node01", Name: "2009-11-10_21.log", Extension: ".log", Size: 1057},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_22.log", Subpath: "node/node01/2009-11-10_22.log", LogType: "node", Target: "node01", Name: "2009-11-10_22.log", Extension: ".log", Size: 177},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node02/2009-11-01_00.log", Subpath: "node/node02/2009-11-01_00.log", LogType: "node", Target: "node02", Name: "2009-11-01_00.log", Extension: ".log", Size: 0},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node02/2009-11-10_21.log", Subpath: "node/node02/2009-11-10_21.log", LogType: "node", Target: "node02", Name: "2009-11-10_21.log", Extension: ".log", Size: 1116},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2000-01-01_00.log", Subpath: "pod/namespace01/2000-01-01_00.log", LogType: "pod", Target: "namespace01", Name: "2000-01-01_00.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_21.log", Subpath: "pod/namespace01/2009-11-10_21.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_21.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_22.log", Subpath: "pod/namespace01/2009-11-10_22.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_22.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2029-11-10_23.log", Subpath: "pod/namespace01/2029-11-10_23.log", LogType: "pod", Target: "namespace01", Name: "2029-11-10_23.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace02/2009-11-10_22.log", Subpath: "pod/namespace02/2009-11-10_22.log", LogType: "pod", Target: "namespace02", Name: "2009-11-10_22.log", Extension: ".log", Size: 1125}},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d__%s", i, tc.retentionTime), func(t *testing.T) {
			fileService2.config.SetRetentionTime(tc.retentionTime)
			err := fileService2.DeleteByAge()
			assert.NoError(t, err)

			got, err := fileService2.ListFiles()
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
			testutil.ResetLogData()
		})
	}
}

func TestDeleteBySize(t *testing.T) {
	testCases := []struct {
		retentionSize int
		want          []LogFile
	}{
		{
			0, // disabled
			[]LogFile{
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_21.log", Subpath: "node/node01/2009-11-10_21.log", LogType: "node", Target: "node01", Name: "2009-11-10_21.log", Extension: ".log", Size: 1057},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_22.log", Subpath: "node/node01/2009-11-10_22.log", LogType: "node", Target: "node01", Name: "2009-11-10_22.log", Extension: ".log", Size: 177},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node02/2009-11-01_00.log", Subpath: "node/node02/2009-11-01_00.log", LogType: "node", Target: "node02", Name: "2009-11-01_00.log", Extension: ".log", Size: 0},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node02/2009-11-10_21.log", Subpath: "node/node02/2009-11-10_21.log", LogType: "node", Target: "node02", Name: "2009-11-10_21.log", Extension: ".log", Size: 1116},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2000-01-01_00.log", Subpath: "pod/namespace01/2000-01-01_00.log", LogType: "pod", Target: "namespace01", Name: "2000-01-01_00.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_21.log", Subpath: "pod/namespace01/2009-11-10_21.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_21.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_22.log", Subpath: "pod/namespace01/2009-11-10_22.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_22.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2029-11-10_23.log", Subpath: "pod/namespace01/2029-11-10_23.log", LogType: "pod", Target: "namespace01", Name: "2029-11-10_23.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace02/0000-00-00_00.log", Subpath: "pod/namespace02/0000-00-00_00.log", LogType: "pod", Target: "namespace02", Name: "0000-00-00_00.log", Extension: ".log", Size: 12},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace02/2009-11-10_22.log", Subpath: "pod/namespace02/2009-11-10_22.log", LogType: "pod", Target: "namespace02", Name: "2009-11-10_22.log", Extension: ".log", Size: 1125}},
		},
		{
			10, // 10 bytes
			[]LogFile{},
		},
		{
			1 * 1024, // 1 KiB
			[]LogFile{
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2029-11-10_23.log", Subpath: "pod/namespace01/2029-11-10_23.log", LogType: "pod", Target: "namespace01", Name: "2029-11-10_23.log", Extension: ".log", Size: 279}},
		},
		{
			2 * 1024, // 2 KiB
			[]LogFile{
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2029-11-10_23.log", Subpath: "pod/namespace01/2029-11-10_23.log", LogType: "pod", Target: "namespace01", Name: "2029-11-10_23.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace02/2009-11-10_22.log", Subpath: "pod/namespace02/2009-11-10_22.log", LogType: "pod", Target: "namespace02", Name: "2009-11-10_22.log", Extension: ".log", Size: 1125}},
		},
		{
			3 * 1024, // 3 KiB
			[]LogFile{
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_22.log", Subpath: "node/node01/2009-11-10_22.log", LogType: "node", Target: "node01", Name: "2009-11-10_22.log", Extension: ".log", Size: 177},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_21.log", Subpath: "pod/namespace01/2009-11-10_21.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_21.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_22.log", Subpath: "pod/namespace01/2009-11-10_22.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_22.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2029-11-10_23.log", Subpath: "pod/namespace01/2029-11-10_23.log", LogType: "pod", Target: "namespace01", Name: "2029-11-10_23.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace02/2009-11-10_22.log", Subpath: "pod/namespace02/2009-11-10_22.log", LogType: "pod", Target: "namespace02", Name: "2009-11-10_22.log", Extension: ".log", Size: 1125}},
		},
		{
			1 * 1024 * 1024, // 1 MiB
			[]LogFile{
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_21.log", Subpath: "node/node01/2009-11-10_21.log", LogType: "node", Target: "node01", Name: "2009-11-10_21.log", Extension: ".log", Size: 1057},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_22.log", Subpath: "node/node01/2009-11-10_22.log", LogType: "node", Target: "node01", Name: "2009-11-10_22.log", Extension: ".log", Size: 177},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node02/2009-11-01_00.log", Subpath: "node/node02/2009-11-01_00.log", LogType: "node", Target: "node02", Name: "2009-11-01_00.log", Extension: ".log", Size: 0},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node02/2009-11-10_21.log", Subpath: "node/node02/2009-11-10_21.log", LogType: "node", Target: "node02", Name: "2009-11-10_21.log", Extension: ".log", Size: 1116},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2000-01-01_00.log", Subpath: "pod/namespace01/2000-01-01_00.log", LogType: "pod", Target: "namespace01", Name: "2000-01-01_00.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_21.log", Subpath: "pod/namespace01/2009-11-10_21.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_21.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_22.log", Subpath: "pod/namespace01/2009-11-10_22.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_22.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2029-11-10_23.log", Subpath: "pod/namespace01/2029-11-10_23.log", LogType: "pod", Target: "namespace01", Name: "2029-11-10_23.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace02/0000-00-00_00.log", Subpath: "pod/namespace02/0000-00-00_00.log", LogType: "pod", Target: "namespace02", Name: "0000-00-00_00.log", Extension: ".log", Size: 12},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace02/2009-11-10_22.log", Subpath: "pod/namespace02/2009-11-10_22.log", LogType: "pod", Target: "namespace02", Name: "2009-11-10_22.log", Extension: ".log", Size: 1125}},
		},
		{
			999999999 * 1024 * 1024 * 1024, // 999,999,999 GiB
			[]LogFile{
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_21.log", Subpath: "node/node01/2009-11-10_21.log", LogType: "node", Target: "node01", Name: "2009-11-10_21.log", Extension: ".log", Size: 1057},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node01/2009-11-10_22.log", Subpath: "node/node01/2009-11-10_22.log", LogType: "node", Target: "node01", Name: "2009-11-10_22.log", Extension: ".log", Size: 177},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node02/2009-11-01_00.log", Subpath: "node/node02/2009-11-01_00.log", LogType: "node", Target: "node02", Name: "2009-11-01_00.log", Extension: ".log", Size: 0},
				{Fullpath: "tmp/storage_fileservice_delete_test/node/node02/2009-11-10_21.log", Subpath: "node/node02/2009-11-10_21.log", LogType: "node", Target: "node02", Name: "2009-11-10_21.log", Extension: ".log", Size: 1116},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2000-01-01_00.log", Subpath: "pod/namespace01/2000-01-01_00.log", LogType: "pod", Target: "namespace01", Name: "2000-01-01_00.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_21.log", Subpath: "pod/namespace01/2009-11-10_21.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_21.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2009-11-10_22.log", Subpath: "pod/namespace01/2009-11-10_22.log", LogType: "pod", Target: "namespace01", Name: "2009-11-10_22.log", Extension: ".log", Size: 1031},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace01/2029-11-10_23.log", Subpath: "pod/namespace01/2029-11-10_23.log", LogType: "pod", Target: "namespace01", Name: "2029-11-10_23.log", Extension: ".log", Size: 279},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace02/0000-00-00_00.log", Subpath: "pod/namespace02/0000-00-00_00.log", LogType: "pod", Target: "namespace02", Name: "0000-00-00_00.log", Extension: ".log", Size: 12},
				{Fullpath: "tmp/storage_fileservice_delete_test/pod/namespace02/2009-11-10_22.log", Subpath: "pod/namespace02/2009-11-10_22.log", LogType: "pod", Target: "namespace02", Name: "2009-11-10_22.log", Extension: ".log", Size: 1125}},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d__%d", i, tc.retentionSize), func(t *testing.T) {
			fileService2.config.SetRetentionSize(tc.retentionSize)
			err := fileService2.DeleteBySize()
			assert.NoError(t, err)

			got, err := fileService2.ListFiles()
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
			testutil.ResetLogData()
		})
	}
}
