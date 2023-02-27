package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/lethe/letheql"
	"github.com/kuoss/lethe/util"
)

func (lh *LetheHandler) Query(c *gin.Context) {
	query := c.Query("query")
	log.Println("QueryHandler", "query=", query)

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "empty query",
		})
		return
	}

	queryData, err := letheql.ProcQuery(query, letheql.TimeRange{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  fmt.Sprintf("%s", err),
		})
		return
	}

	if queryData.ResultType == letheql.ValueTypeLogs {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"resultType": "logs",
				"result":     queryData.Logs,
			},
		})
		return
	}

	if queryData.Scalar == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"resultType": "vector",
				"result":     []int{},
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"resultType": "vector",
			"result": []gin.H{
				{
					"value": queryData.Scalar,
				},
			},
		},
	})
}

func (lh *LetheHandler) QueryRange(c *gin.Context) {
	query := c.Query("query")
	start := c.Query("start")
	end := c.Query("end")

	log.Println("query_range...", query, start, end)
	if query == "" || start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "empty query",
		})
		return
	}
	startTime := util.FloatStringToTime(start)
	endTime := util.FloatStringToTime(end)
	queryData, err := letheql.ProcQuery(query, letheql.TimeRange{Start: startTime, End: endTime})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  fmt.Sprintf("%s", err),
		})
		return
	}
	fmt.Println("queryData.ResultType=", queryData.ResultType)
	if queryData.ResultType == letheql.ValueTypeLogs {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"resultType": "logs",
				"result":     queryData.Logs,
			},
		})
		return
	}
	if queryData.Scalar == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"resultType": "vector",
				"result":     []int{},
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"resultType": "vector",
			"result": []gin.H{
				{
					"value": queryData.Scalar,
				},
			},
		},
	})
}
