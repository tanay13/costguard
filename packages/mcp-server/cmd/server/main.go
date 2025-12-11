package main

import (
	"net/http"
	"os"
	"sort"

	"github.com/gin-gonic/gin"
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
	var req struct {
		Metrics        []types.MetricCollection  `json:"metrics"`
		ActualRequests map[string]types.Requests `json:"actual_requests"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	agg := scan.DataPointAggregator(req.Metrics, req.ActualRequests)

	response := buildScanResponse(agg)

	c.JSON(200, response)
}

func buildScanResponse(agg []types.AggregatedMetrics) types.ScanResponse {
	resp := types.ScanResponse{}
	totalCurrent := 0.0
	totalOptimal := 0.0

	// sort by savings later
	temp := []struct {
		name    string
		savings float64
		data    types.ScanResource
	}{}

	for _, a := range agg {

		resource := types.ScanResource{
			Provider: a.Provider,
			Resource: a.Resource,
			Usage:    a.Metrics,
			Costs: types.ScanResourceCost{
				CurrentCostUSD:      a.CostCurrentUSD,
				OptimalCostUSD:      a.CostOptimalUSD,
				PotentialSavingsUSD: a.CostSavingsUSD,
				WastePercentage:     (a.CostSavingsUSD / a.CostCurrentUSD) * 100,
			},
		}

		resource.Requested.CpuMilli = a.RequestedCpuMilli
		resource.Requested.MemoryGB = a.RequestedMemoryGB

		totalCurrent += a.CostCurrentUSD
		totalOptimal += a.CostOptimalUSD

		temp = append(temp, struct {
			name    string
			savings float64
			data    types.ScanResource
		}{
			name:    a.Resource,
			savings: a.CostSavingsUSD,
			data:    resource,
		})
	}

	sort.Slice(temp, func(i, j int) bool {
		return temp[i].savings > temp[j].savings
	})

	top := []string{}
	for i := 0; i < len(temp) && i < 3; i++ {
		top = append(top, temp[i].name)
		resp.Resources = append(resp.Resources, temp[i].data)
	}

	resp.Summary = types.ScanSummary{
		TotalCurrentCostUSD:      totalCurrent,
		TotalOptimalCostUSD:      totalOptimal,
		TotalPotentialSavingsUSD: totalCurrent - totalOptimal,
		TopOffenders:             top,
	}

	return resp
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
