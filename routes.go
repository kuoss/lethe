package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/lethe/handlers"
	"github.com/kuoss/lethe/logs/rotator"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	routesRootGroup(r)
	routesAPIV1Group(r)
	return r
}

func routesRootGroup(r *gin.Engine) {

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

}

func routesAPIV1Group(r *gin.Engine) {

	v1 := r.Group("api/v1")
	h := handlers.LetheHandler{Rotator: rotator.NewRotator()}

	v1.GET("query", h.Query)
	v1.GET("query_range", h.QueryRange)
	v1.GET("metadata", h.Metadata)
	v1.GET("targets", h.Target)
}
