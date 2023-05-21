package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func testHandlerGet(relativePath string, handler gin.HandlerFunc) (code int, body string) {
	r := gin.New()
	r.GET(relativePath, handler)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", relativePath, nil)
	if err != nil {
		panic(err)
	}
	r.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}

func TestHealthy(t *testing.T) {
	code, body := testHandlerGet("/-/healthy", handler1.Healthy)
	assert.Equal(t, 200, code)
	assert.Equal(t, "Venti is Healthy.\n", body)
}

func TestReady(t *testing.T) {
	code, body := testHandlerGet("/-/ready", handler1.Ready)
	assert.Equal(t, 200, code)
	assert.Equal(t, "Venti is Ready.\n", body)
}
