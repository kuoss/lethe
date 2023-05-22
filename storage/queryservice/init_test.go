package queryservice

import (
	"github.com/kuoss/lethe/clock"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/storage/fileservice"
	"github.com/kuoss/lethe/storage/logservice"
	"github.com/kuoss/lethe/util/testutil"
)

var (
	queryService *QueryService
)

func init() {
	testutil.ChdirRoot()
	testutil.ResetLogData()
	clock.SetPlaygroundMode(true)

	cfg, err := config.New("test")
	if err != nil {
		panic(err)
	}
	cfg.SetLogDataPath("tmp/init")
	fileService, err := fileservice.New(cfg)
	if err != nil {
		panic(err)
	}
	logService := logservice.New(fileService)
	queryService = New(logService)
}
