package utils

import (
	"math"
	"sort"

	"github.com/tanay13/costguard/packages/mcp-server/pkg/types"
)

func ConvertScanToAggregated(res types.ScanResponse) []types.AggregatedMetrics {
	out := []types.AggregatedMetrics{}

	for _, r := range res.Resources {

		usage := map[string]types.MetricStat{
			"cpu":    r.Usage["cpu_milli"],
			"memory": r.Usage["memory_gb"],
		}

		agg := types.AggregatedMetrics{
			Provider: r.Provider,
			Resource: r.Resource,
			Metrics:  usage,

			RequestedCpuMilli: r.Requested.CpuMilli,
			RequestedMemoryGB: r.Requested.MemoryGB,

			CostCurrentUSD: r.Costs.CurrentCostUSD,
			CostOptimalUSD: r.Costs.OptimalCostUSD,
			CostSavingsUSD: r.Costs.PotentialSavingsUSD,

			DataPoints: 1,
		}

		out = append(out, agg)
	}

	return out
}

func ComputePercentile(data []float64, percentile float64) float64 {
	if len(data) == 0 || percentile < 0 || percentile > 100 {
		return 0
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i] < data[j]
	})

	N := float64(len(data))
	P := float64(percentile)

	rank := P / 100.0 * (N - 1)

	finalRank := int(math.Round(rank))

	return float64(data[finalRank])
}

func CalculateAvg(data []float64) float64 {
	if len(data) == 0 {
		return 0.0
	}

	var sum float64 = 0

	for _, d := range data {
		sum += d
	}
	return float64(sum) / float64(len(data))
}
