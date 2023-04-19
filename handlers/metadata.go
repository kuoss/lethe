package handlers

import (
	"fmt"
	"github.com/kuoss/lethe/logs/rotator"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/lethe/util"
	"github.com/thoas/go-funk"
)

type LetheHandler struct {
	*rotator.Rotator
}

func (lh *LetheHandler) Metadata(c *gin.Context) {
	dirs := lh.ListDirs()
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

func (lh *LetheHandler) Target(c *gin.Context) {
	activeTargets := []gin.H{}
	for _, target := range lh.ListTargets() {
		labelKey := "__meta_kubernetes_node_name"
		if target.LogType == "pod" {
			labelKey = "__meta_kubernetes_namespace"
		}
		activeTargets = append(activeTargets, gin.H{
			"lastScrape": target.LastForward,
			"health":     "up",
			"discoveredLabels": gin.H{
				"job":    target.LogType,
				labelKey: target.Target,
			},
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"activeTargets": activeTargets,
		},
	})
}
