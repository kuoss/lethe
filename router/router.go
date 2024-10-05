package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/storage/fileservice"
	"github.com/kuoss/lethe/storage/queryservice"
)

type Router struct {
	config       *config.Config
	fileService  *fileservice.FileService
	queryService *queryservice.QueryService
	ginRouter    *gin.Engine
}

func New(cfg *config.Config, fileService *fileservice.FileService, queryService *queryservice.QueryService) *Router {
	router := &Router{
		config:       cfg,
		fileService:  fileService,
		queryService: queryService,
	}
	router.setupRouter()
	return router
}

func (r *Router) setupRouter() {
	// gin.SetMode(r.config.Web.GinMode)
	ginRouter := gin.Default()

	ginRouter.GET("/-/healthy", r.Healthy)
	ginRouter.GET("/-/ready", r.Ready)

	ginRouter.GET("/api/v1/metadata", r.Metadata)
	ginRouter.GET("/api/v1/query", r.Query)
	ginRouter.GET("/api/v1/query_range", r.QueryRange)
	ginRouter.GET("/api/v1/targets", r.Target)

	r.ginRouter = ginRouter
}

func (r *Router) Run() error {
	return r.ginRouter.Run(r.config.Web.ListenAddress)
}
