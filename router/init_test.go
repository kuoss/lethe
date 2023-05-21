package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/lethe/config"
	_ "github.com/kuoss/lethe/storage/driver/filesystem"
	"github.com/kuoss/lethe/storage/fileservice"
	"github.com/kuoss/lethe/util/testutil"
)

var (
	router1 *gin.Engine
)

func init() {
	testutil.ChdirRoot()
	testutil.ResetLogData()

	cfg, err := config.New("test")
	if err != nil {
		panic(err)
	}

	fs, err := fileservice.New(cfg)
	if err != nil {
		panic(err)
	}
	router1 = New(fs)
}
