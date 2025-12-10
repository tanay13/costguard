package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tanay13/costguard/packages/mcp-server/internal/provider"
	"github.com/tanay13/costguard/packages/mcp-server/internal/scan"
	"github.com/tanay13/costguard/packages/mcp-server/internal/types"
)

func HealthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "codeguard is up..",
		"status":  "OK",
	})
}

func ScanHandler(c *gin.Context) {
	var requestData []types.MetricCollection

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not a valid JSON",
		})
		return
	}
	resp := scan.DataPointAggregator(requestData)

	c.JSON(http.StatusOK, resp)
}

func FixPlansHandler(c *gin.Context) {
	var raw []types.MetricCollection

	if err := c.BindJSON(&raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agg := scan.DataPointAggregator(raw)

	req := types.FixPlanRequest{
		AggregatedMetrics: agg,
		BudgetTarget:      500,
		AutoApprove:       false,
	}

	resp := provider.GenerateFixPlan(req)

	c.JSON(200, resp)
}

func EventsHandler(c *gin.Context) {
	jsonBytes, err := os.ReadFile("events.json")
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"message": "No file found at that path"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error reading from the file"})
		return
	}

	c.Data(200, "application/json", jsonBytes)
}

func main() {
	router := gin.Default()
	router.GET("/health", HealthHandler)
	router.POST("/v1/scan", ScanHandler)
	router.POST("/v1/fixplans", FixPlansHandler)
	router.GET("/v1/events", EventsHandler)
	router.Run()
}
