package letheql

import (
	"context"
	"time"

	"github.com/kuoss/lethe/clock"
	"github.com/kuoss/lethe/config"
	_ "github.com/kuoss/lethe/storage/driver/filesystem"
	"github.com/kuoss/lethe/storage/fileservice"
	"github.com/kuoss/lethe/storage/logservice"
	"github.com/kuoss/lethe/util/testutil"
)

var (
	engine1    *Engine
	evaluator1 *evaluator
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
	engine1 = NewEngine(logService)

	now := clock.Now()
	evaluator1 = &evaluator{
		logService:     logService,
		ctx:            context.TODO(),
		start:          now.Add(-4 * time.Hour),
		end:            now,
		startTimestamp: 0,
		endTimestamp:   0,
		interval:       0,
	}
}
