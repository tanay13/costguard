package provider

import (
	"github.com/tanay13/costguard/packages/mcp-server/internal/types"
	"github.com/tanay13/costguard/packages/mcp-server/internal/utils"
)

const (
	CpuUsageCost    = 0.01
	MemoryUsageCost = 0.05
)

func KubernetesMetricAggregator(metrics []types.K8sResourceMetrics) []types.AggregatedMetrics {
	cpuUsage := make([]float64, 0)
	memoryUsage := make([]float64, 0)

	aggregatedMetrics := make([]types.AggregatedMetrics, 0)

	for _, r := range metrics {
		cpuUsage = append(cpuUsage, r.CpuMilli)
		memoryUsage = append(memoryUsage, r.MemoryGB)
	}

	p50_cpu_usage := utils.ComputePercentile(cpuUsage, 50)
	p95_cpu_usage := utils.ComputePercentile(cpuUsage, 95)
	average_cpu_usage := utils.CalculateAvg(cpuUsage)
	monthly_cost_cpu := average_cpu_usage * CpuUsageCost

	cpu_metrics := types.AggregatedMetrics{
		Provider:       types.ProviderKubernetes,
		Resource:       "CpuMilli",
		P50:            p50_cpu_usage,
		P95:            p95_cpu_usage,
		Avg:            average_cpu_usage,
		MonthlyCostUSD: monthly_cost_cpu,
	}

	p50_memory_usage := utils.ComputePercentile(memoryUsage, 50)
	p95_memory_usage := utils.ComputePercentile(memoryUsage, 95)
	average_memory_usage := utils.CalculateAvg(memoryUsage)
	monthly_cost_memory := average_memory_usage * MemoryUsageCost

	memory_metrics := types.AggregatedMetrics{
		Provider:       types.ProviderKubernetes,
		Resource:       "memoryGB",
		P50:            p50_memory_usage,
		P95:            p95_memory_usage,
		Avg:            average_memory_usage,
		MonthlyCostUSD: monthly_cost_memory,
	}
	aggregatedMetrics = append(aggregatedMetrics, cpu_metrics)
	aggregatedMetrics = append(aggregatedMetrics, memory_metrics)

	return aggregatedMetrics
}
