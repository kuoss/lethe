package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) Target(c *gin.Context) {
	activeTargets := []gin.H{}
	for _, target := range r.fileService.ListTargets() {
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
