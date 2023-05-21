package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Healthy(c *gin.Context) {
	c.String(http.StatusOK, "Venti is Healthy.\n")
}

func (h *Handler) Ready(c *gin.Context) {
	c.String(http.StatusOK, "Venti is Ready.\n")
}
