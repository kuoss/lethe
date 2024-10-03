package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) Healthy(c *gin.Context) {
	c.String(http.StatusOK, "Venti is Healthy.\n")
}

func (r *Router) Ready(c *gin.Context) {
	c.String(http.StatusOK, "Venti is Ready.\n")
}
