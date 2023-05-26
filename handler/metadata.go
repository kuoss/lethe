package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Metadata(c *gin.Context) {
	targets := []string{}
	dirs := h.fileService.ListTargets()

	for _, d := range dirs {
		var key string
		switch d.LogType {
		case "pod":
			key = "namespace"
		case "node":
			key = "node"
		}
		value := filepath.Base(d.Subpath)
		targets = append(targets, fmt.Sprintf(`%s{%s="%s"}`, d.LogType, key, value))
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"targets": targets,
		},
	})
}
