package provider

import (
	"fmt"
	"math"

	"github.com/tanay13/costguard/packages/mcp-server/internal/types"
)

const (
	_cpuRatePerMilli = 0.00001
	_memRatePerGB    = 0.00002
)

func generateK8sFixActions(agg types.AggregatedMetrics, priority int) []types.FixAction {
	out := []types.FixAction{}
	cpuStat, hasCPU := agg.Metrics["cpu_milli"]
	memStat, hasMem := agg.Metrics["memory_gb"]

	requestedCpu := agg.RequestedCpuMilli
	requestedMem := agg.RequestedMemoryGB
	if requestedCpu <= 0 && hasCPU {
		requestedCpu = math.Max(cpuStat.P95*2.0, 50.0)
	}
	if requestedMem <= 0 && hasMem {
		requestedMem = math.Max(memStat.P95*2.0, 0.1)
	}

	if hasCPU && cpuStat.Avg < requestedCpu*0.5 {
		percent := -40.0
		out = append(out, buildK8sAction(agg, "rightsize_cpu_request", "resources.requests.cpu", percent, fmt.Sprintf("CPU avg %.2fm < 50%% of requested %.2fm", cpuStat.Avg, requestedCpu), agg.CostUSD, priority))
	}

	if hasCPU && cpuStat.P95 > requestedCpu*0.9 {
		percent := 25.0
		out = append(out, buildK8sAction(agg, "rightsize_cpu_request", "resources.requests.cpu", percent, fmt.Sprintf("CPU p95 %.2fm > 90%% of requested %.2fm", cpuStat.P95, requestedCpu), agg.CostUSD, priority))
	}

	if hasMem && memStat.Avg < requestedMem*0.5 {
		percent := -35.0
		out = append(out, buildK8sAction(agg, "rightsize_memory_request", "resources.requests.memory", percent, fmt.Sprintf("Memory avg %.3fGB < 50%% of requested %.3fGB", memStat.Avg, requestedMem), agg.CostUSD, priority))
	}
	if hasMem && memStat.P95 > requestedMem*0.9 {
		percent := 20.0
		out = append(out, buildK8sAction(agg, "rightsize_memory_request", "resources.requests.memory", percent, fmt.Sprintf("Memory p95 %.3fGB > 90%% of requested %.3fGB", memStat.P95, requestedMem), agg.CostUSD, priority))
	}

	return out
}

func buildK8sAction(agg types.AggregatedMetrics, intent, field string, percent float64, explanation string, currentCost float64, priority int) types.FixAction {
	expected := currentCost * (1.0 + (percent / 100.0))
	op := types.FixOperation{
		Field:     field,
		Operation: "scale_by_percentage",
		Value:     percent,
		Unit:      "percentage",
	}

	guidance := fmt.Sprintf("Locate manifests for '%s' and apply: %s %+g%%. Preserve units. If missing, create field.", agg.Resource, field, percent)

	return types.FixAction{
		Provider:     types.ProviderKubernetes,
		Resource:     agg.Resource,
		Intent:       intent,
		Description:  explanation,
		Action:       op,
		CurrentCost:  currentCost,
		ExpectedCost: expected,
		Savings:      currentCost - expected,
		Risk:         "low",
		Priority:     priority,
		AIGuidance:   guidance,
		FilesToEdit:  []string{},
	}
}
