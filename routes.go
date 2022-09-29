package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/lethe/file"
	"github.com/kuoss/lethe/letheql"
	"github.com/kuoss/lethe/util"
	"github.com/thoas/go-funk"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	routesRootGroup(r)
	routesAPIV1Group(r)
	return r
}

func routesRootGroup(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
}

func routesAPIV1Group(r *gin.Engine) {
	v1 := r.Group("api/v1")
	v1.GET("query", func(c *gin.Context) {
		// log.Println("QueryHandler...")
		query := c.Query("query")
		log.Println("query=", query)
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  "empty query",
			})
			return
		}
		queryData, err := letheql.ProcQuery(query, letheql.TimeRange{})
		// log.Println("queryData=", queryData)
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
	})

	v1.GET("query_range", func(c *gin.Context) {
		query := c.Query("query")
		start := c.Query("start")
		end := c.Query("end")
		// log.Println("query_range...", query, start, end)
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
	})

	v1.GET("metadata", func(c *gin.Context) {
		dirs := file.ListDirs()
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
	})

	v1.GET("targets", func(c *gin.Context) {
		activeTargets := []gin.H{}
		for _, target := range file.ListTargets() {
			labelKey := "__meta_kubernetes_node_name"
			if target.Typ == "pod" {
				labelKey = "__meta_kubernetes_namespace"
			}
			activeTargets = append(activeTargets, gin.H{
				"lastScrape": target.LastForward,
				"health":     "up",
				"discoveredLabels": gin.H{
					"job":    target.Typ,
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
	})
}
