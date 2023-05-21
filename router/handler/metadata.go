package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/lethe/util"
	"github.com/thoas/go-funk"
)

func (h *Handler) Metadata(c *gin.Context) {
	dirs := h.fileService.ListDirs()
	targets := funk.Map(dirs, func(x string) string {
		typ := util.SubstrBefore(x, "/")
		value := util.SubstrAfter(x, "/")
		var key string
		switch typ {
		case "pod":
			key = "namespace"
		case "node":
			key = "node"
		}
		return fmt.Sprintf(`%s{%s="%s"}`, typ, key, value)
	})
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"targets": targets,
		},
	})
}
