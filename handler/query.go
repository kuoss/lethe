package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/lethe/letheql/model"
	"github.com/kuoss/lethe/util"
)

func (h *Handler) Query(c *gin.Context) {
	qs := c.Query("query")
	h.query(c, qs, model.TimeRange{})
}

func (h *Handler) query(c *gin.Context, qs string, tr model.TimeRange) {
	result := h.queryService.Query(c.Request.Context(), qs, tr)

	// error
	if result.Err != nil {
		obj := gin.H{
			"status":    "error",
			"errorType": "queryError",
			"error":     result.Err.Error(),
		}
		if len(result.Warnings) > 0 {
			obj["warnings"] = result.Warnings
		}
		c.JSON(http.StatusInternalServerError, obj)
		return
	}

	// log
	if result.Value.Type() == model.ValueTypeLog {
		log, ok := result.Value.(model.Log)
		if ok {
			obj := gin.H{
				"status": "success",
				"data": gin.H{
					"resultType": "logs",
					"result":     log.Lines,
				},
			}
			if len(result.Warnings) > 0 {
				obj["warnings"] = result.Warnings
			}
			c.JSON(http.StatusOK, obj)
			return
		}
	}

	// any
	obj := gin.H{
		"status": "success",
		"data": gin.H{
			"resultType": result.Value.Type(),
			"result":     result.Value,
		},
	}
	if len(result.Warnings) > 0 {
		obj["warnings"] = result.Warnings
	}
	c.JSON(http.StatusOK, obj)
}

func (h *Handler) QueryRange(c *gin.Context) {
	qs := c.Query("query")
	start := c.Query("start")
	end := c.Query("end")
	h.query(c, qs, model.TimeRange{
		Start: util.FloatStringToTime(start),
		End:   util.FloatStringToTime(end),
	})
}
