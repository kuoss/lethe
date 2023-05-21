package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/lethe/router/handler"
	"github.com/kuoss/lethe/storage/fileservice"
)

func New(fileService *fileservice.FileService) *gin.Engine {

	h := handler.New(fileService)

	r := gin.Default()
	r.GET("/-/healthy", h.Healthy)
	r.GET("/-/ready", h.Ready)

	r.GET("/api/v1/metadata", h.Metadata)
	r.GET("/api/v1/query", h.Query)
	r.GET("/api/v1/query_range", h.QueryRange)
	r.GET("/api/v1/targets", h.Target)
	return r
}
