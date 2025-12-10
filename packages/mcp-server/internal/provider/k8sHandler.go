package provider

import (
	"github.com/tanay13/costguard/packages/mcp-server/internal/types"
	"github.com/tanay13/costguard/packages/mcp-server/internal/utils"
)

const (
	CpuUsageCost    = 0.01
	MemoryUsageCost = 0.05
)

func KubernetesMetricAggregator(groups map[string][]types.MetricCollection) []types.AggregatedMetrics {
	out := make([]types.AggregatedMetrics, 0)

	for resourceName, points := range groups {
		cpuValues := make([]float64, 0, len(points))
		memValues := make([]float64, 0, len(points))

		for _, p := range points {
			k := p.Metrics.K8sResourceMetrics
			cpuValues = append(cpuValues, k.CpuMilli)
			memValues = append(memValues, k.MemoryGB)
		}

		p50cpu := utils.ComputePercentile(cpuValues, 50)
		p95cpu := utils.ComputePercentile(cpuValues, 95)
		avgcpu := utils.CalculateAvg(cpuValues)

		p50mem := utils.ComputePercentile(memValues, 50)
		p95mem := utils.ComputePercentile(memValues, 95)
		avgmem := utils.CalculateAvg(memValues)

		cost := (avgcpu * CpuUsageCost) + (avgmem * MemoryUsageCost)

		metrics := map[string]types.MetricStat{
			"cpu_milli": {P50: p50cpu, P95: p95cpu, Avg: avgcpu},
			"memory_gb": {P50: p50mem, P95: p95mem, Avg: avgmem},
		}

		agg := types.AggregatedMetrics{
			Provider:   types.ProviderKubernetes,
			Resource:   resourceName,
			Metrics:    metrics,
			CostUSD:    cost,
			DataPoints: len(points),
		}

		out = append(out, agg)
	}

	return out
}
