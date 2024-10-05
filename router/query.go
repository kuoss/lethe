package router

import (
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/lethe/letheql/model"
	"github.com/spf13/cast"
)

func (r *Router) Query(c *gin.Context) {
	qs := c.Query("query")
	r.query(c, qs, model.TimeRange{})
}

func (r *Router) QueryRange(c *gin.Context) {
	qs := c.Query("query")
	start := c.Query("start")
	end := c.Query("end")
	r.query(c, qs, model.TimeRange{Start: floatStringToTime(start), End: floatStringToTime(end)})
}

func (r *Router) query(c *gin.Context, qs string, tr model.TimeRange) {
	result := r.queryService.Query(c.Request.Context(), qs, tr)

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

func floatStringToTime(timeFloat string) time.Time {
	sec, dec := math.Modf(cast.ToFloat64(timeFloat))
	return time.Unix(int64(sec), int64(dec*(1e9)))
}
