package provider

import (
	"math"
	"sort"

	"github.com/tanay13/costguard/packages/mcp-server/internal/types"
)

func KubernetesMetricAggregator(resources map[string][]types.MetricCollection) []types.AggregatedMetrics {
	out := []types.AggregatedMetrics{}

	for name, pts := range resources {
		cpuVals := make([]float64, 0, len(pts))
		memVals := make([]float64, 0, len(pts))

		for _, p := range pts {
			cpuVals = append(cpuVals, p.Metrics.K8sResourceMetrics.CpuMilli)
			memVals = append(memVals, p.Metrics.K8sResourceMetrics.MemoryGB)
		}

		sort.Float64s(cpuVals)
		sort.Float64s(memVals)

		cpuStat := types.MetricStat{
			P50: computePercentile(cpuVals, 50),
			P95: computePercentile(cpuVals, 95),
			Avg: calculateAvg(cpuVals),
		}

		memStat := types.MetricStat{
			P50: computePercentile(memVals, 50),
			P95: computePercentile(memVals, 95),
			Avg: calculateAvg(memVals),
		}

		// Calculate costs based on P95 usage
		requestedCPU := math.Max(cpuStat.P95*2.0, 50.0)
		requestedMem := math.Max(memStat.P95*2.0, 0.1)

		cpuRatePerMilli := 0.00001
		memRatePerGB := 0.00002

		currentCost := (requestedCPU * cpuRatePerMilli) + (requestedMem * memRatePerGB)

		out = append(out, types.AggregatedMetrics{
			Provider:          types.ProviderKubernetes,
			Resource:          name,
			Metrics:           map[string]types.MetricStat{"cpu_milli": cpuStat, "memory_gb": memStat},
			RequestedCpuMilli: requestedCPU,
			RequestedMemoryGB: requestedMem,
			CostUSD:           currentCost,
			DataPoints:        len(pts),
		})
	}

	return out
}

func computePercentile(data []float64, percentile float64) float64 {
	if len(data) == 0 || percentile < 0 || percentile > 100 {
		return 0
	}

	sort.Float64s(data)

	N := float64(len(data))
	P := float64(percentile)

	rank := P / 100.0 * (N - 1)
	finalRank := int(math.Round(rank))

	if finalRank >= len(data) {
		finalRank = len(data) - 1
	}

	return data[finalRank]
}

func calculateAvg(data []float64) float64 {
	if len(data) == 0 {
		return 0.0
	}

	var sum float64 = 0
	for _, d := range data {
		sum += d
	}
	return sum / float64(len(data))
}
