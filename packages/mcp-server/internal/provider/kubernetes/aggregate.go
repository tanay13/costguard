package kubernetes

import (
	"fmt"
	"math"
	"sort"

	"github.com/tanay13/costguard/packages/mcp-server/internal/types"
	"github.com/tanay13/costguard/packages/mcp-server/internal/utils"
)

func Aggregate(
	resources map[string][]types.MetricCollection,
	actualReqs map[string]types.Requests,
) []types.AggregatedMetrics {
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
			P50: utils.ComputePercentile(cpuVals, 50),
			P95: utils.ComputePercentile(cpuVals, 95),
			Avg: utils.CalculateAvg(cpuVals),
		}

		memStat := types.MetricStat{
			P50: utils.ComputePercentile(memVals, 50),
			P95: utils.ComputePercentile(memVals, 95),
			Avg: utils.CalculateAvg(memVals),
		}

		reqCPU, reqMem := ResolveRequests(name, cpuStat, memStat, actualReqs)

		optimalCPU, optimalMem := OptimalRequests(cpuStat, memStat)

		currentCost := ComputeCostFromRequests(reqCPU, reqMem)
		fmt.Println(reqCPU, reqMem, currentCost)
		optimalCost := ComputeCostFromRequests(optimalCPU, optimalMem)
		fmt.Println(optimalCPU, optimalMem, optimalCost)
		out = append(out, types.AggregatedMetrics{
			Provider: types.ProviderKubernetes,
			Resource: name,
			Metrics: map[string]types.MetricStat{
				"cpu_milli": cpuStat,
				"memory_gb": memStat,
			},
			RequestedCpuMilli: reqCPU,
			RequestedMemoryGB: reqMem,
			OptimalCpuMilli:   optimalCPU,
			OptimalMemoryGB:   optimalMem,
			CostCurrentUSD:    currentCost,
			CostOptimalUSD:    optimalCost,
			CostSavingsUSD:    currentCost - optimalCost,
			DataPoints:        len(pts),
		})
	}

	return out
}

func ResolveRequests(
	resource string,
	cpu types.MetricStat,
	mem types.MetricStat,
	actual map[string]types.Requests,
) (float64, float64) {
	if actual != nil {
		if r, ok := actual[resource]; ok {
			if r.CpuMilli > 0 && r.MemoryGB > 0 {
				return r.CpuMilli, r.MemoryGB
			}
		}
	}

	reqCPU := math.Max(cpu.P95*2, 50)
	reqMem := math.Max(mem.P95*2, 0.1)

	return reqCPU, reqMem
}
