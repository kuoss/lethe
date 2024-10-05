package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kuoss/common/tester"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/storage/fileservice"
	"github.com/kuoss/lethe/storage/logservice"
	"github.com/kuoss/lethe/storage/queryservice"
	"github.com/stretchr/testify/require"
)

func newRouter(t *testing.T) (*Router, func()) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/testdata/log": "data/log",
	})
	cfg, err := config.New("test")
	require.NoError(t, err)
	fileService, err := fileservice.New(cfg)
	require.NoError(t, err)
	logService := logservice.New(cfg, fileService)
	queryService := queryservice.New(cfg, logService)
	return New(cfg, fileService, queryService), cleanup
}

func testGET(t *testing.T, url string) (code int, body string, cleanup func()) {
	r, cleanup := newRouter(t)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)
	r.ginRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), cleanup
}

func TestNew(t *testing.T) {
	want := &Router{}
	got, cleanup := newRouter(t)
	defer cleanup()
	require.NotEqual(t, want, got)
}
