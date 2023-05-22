package handler

import (
	"net/http"
	"net/http/httptest"

	"github.com/kuoss/lethe/clock"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/storage/fileservice"
	"github.com/kuoss/lethe/storage/logservice"
	"github.com/kuoss/lethe/storage/queryservice"
	"github.com/kuoss/lethe/util/testutil"
)

var (
	handler1 *Handler
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
	queryService := queryservice.New(logService)
	handler1 = New(cfg, fileService, queryService)
}

func testGET(url string) (code int, body string) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	handler1.router.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}
