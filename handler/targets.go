package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Target(c *gin.Context) {
	activeTargets := []gin.H{}
	for _, target := range h.fileService.ListTargets() {
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
