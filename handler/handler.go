package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/storage/fileservice"
	"github.com/kuoss/lethe/storage/queryservice"
)

type Handler struct {
	config       *config.Config
	fileService  *fileservice.FileService
	queryService *queryservice.QueryService
	router       *gin.Engine
}

func New(cfg *config.Config, fileService *fileservice.FileService, queryService *queryservice.QueryService) *Handler {
	handler := &Handler{
		config:       cfg,
		fileService:  fileService,
		queryService: queryService,
	}
	handler.setupRouter()
	return handler
}

func (h *Handler) setupRouter() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })

	r.GET("/-/healthy", h.Healthy)
	r.GET("/-/ready", h.Ready)

	r.GET("/api/v1/metadata", h.Metadata)
	r.GET("/api/v1/query", h.Query)
	r.GET("/api/v1/query_range", h.QueryRange)
	r.GET("/api/v1/targets", h.Target)
	h.router = r
}

func (h *Handler) Run() error {
	return h.router.Run(h.config.WebListenAddress())
}
