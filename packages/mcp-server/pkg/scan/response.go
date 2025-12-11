package scan

import (
	"sort"

	"github.com/tanay13/costguard/packages/mcp-server/pkg/types"
)

func BuildScanResponse(agg []types.AggregatedMetrics) types.ScanResponse {
	resp := types.ScanResponse{}
	totalCurrent := 0.0
	totalOptimal := 0.0

	temp := []struct {
		name    string
		savings float64
		data    types.ScanResource
	}{}

	for _, a := range agg {
		res := types.ScanResource{
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

		res.Requested.CpuMilli = a.RequestedCpuMilli
		res.Requested.MemoryGB = a.RequestedMemoryGB

		totalCurrent += a.CostCurrentUSD
		totalOptimal += a.CostOptimalUSD

		temp = append(temp, struct {
			name    string
			savings float64
			data    types.ScanResource
		}{
			name:    a.Resource,
			savings: a.CostSavingsUSD,
			data:    res,
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
