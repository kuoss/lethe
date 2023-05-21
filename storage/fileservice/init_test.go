package fileservice

import (
	"github.com/kuoss/lethe/config"
	_ "github.com/kuoss/lethe/storage/driver/filesystem"
	"github.com/kuoss/lethe/util/testutil"
)

var (
	fileService *FileService
)

func init() {
	testutil.ChdirRoot()
	testutil.ResetLogData()

	cfg, err := config.New("test")
	if err != nil {
		panic(err)
	}
	cfg.SetLogDataPath("tmp/init")
	fileService, err = New(cfg)
	if err != nil {
		panic(err)
	}
}
