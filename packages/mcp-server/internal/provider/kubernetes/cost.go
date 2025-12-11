package kubernetes

import "github.com/tanay13/costguard/packages/mcp-server/internal/types"

const (
	cpuRatePerMilli = 0.00001
	memRatePerGB    = 0.00002
)

func ComputeCostFromRequests(cpuMilli float64, memGB float64) float64 {
	return cpuMilli*cpuRatePerMilli*24*30 +
		memGB*memRatePerGB*24*30
}

func OptimalRequests(cpu types.MetricStat, mem types.MetricStat) (float64, float64) {
	return cpu.P50 * 1.2, mem.P50 * 1.2
}
