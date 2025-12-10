package provider

import (
	"fmt"
	"math"

	"github.com/tanay13/costguard/packages/mcp-server/internal/types"
)

const (
	minCPU = 50.0
	minMem = 0.1
	under  = 0.50
	over   = 0.90
)

func generateK8sFixActions(metric types.AggregatedMetrics, priority int) []types.FixAction {
	actions := []types.FixAction{}

	cpu, hasCPU := metric.Metrics["cpu_milli"]
	mem, hasMem := metric.Metrics["memory_gb"]

	inferredReqCPU := math.Max(cpu.P95*2, minCPU)
	inferredReqMem := math.Max(mem.P95*2, minMem)

	if hasCPU && cpu.P95 < inferredReqCPU*under {
		actions = append(actions, makeK8sFixAction(
			metric.Resource,
			"rightsize_cpu_request",
			"resources.requests.cpu",
			-40,
			"CPU is underutilized compared to inferred request",
			metric.CostUSD,
			priority,
		))
	}

	if hasCPU && cpu.P95 > inferredReqCPU*over {
		actions = append(actions, makeK8sFixAction(
			metric.Resource,
			"rightsize_cpu_request",
			"resources.requests.cpu",
			+25,
			"CPU usage is close to saturation (p95 > 90% of inferred request)",
			metric.CostUSD,
			priority,
		))
	}

	if hasMem && mem.P95 < inferredReqMem*under {
		actions = append(actions, makeK8sFixAction(
			metric.Resource,
			"rightsize_memory_request",
			"resources.requests.memory",
			-35,
			"Memory is underutilized compared to inferred request",
			metric.CostUSD,
			priority,
		))
	}

	if hasMem && mem.P95 > inferredReqMem*over {
		actions = append(actions, makeK8sFixAction(
			metric.Resource,
			"rightsize_memory_request",
			"resources.requests.memory",
			+20,
			"Memory usage is close to saturation",
			metric.CostUSD,
			priority,
		))
	}

	return actions
}

func makeK8sFixAction(
	resource string,
	intent string,
	field string,
	percent float64,
	explanation string,
	currentCost float64,
	priority int,
) types.FixAction {
	expected := currentCost * (1 - (percent / 100.0))
	if percent > 0 {
		expected = currentCost * (1 + (percent / 100.0))
	}

	return types.FixAction{
		Provider:     types.ProviderKubernetes,
		Resource:     resource,
		Intent:       intent,
		Description:  explanation,
		CurrentCost:  currentCost,
		ExpectedCost: expected,
		Savings:      currentCost - expected,
		Risk:         "low",
		Priority:     priority,

		Action: types.FixOperation{
			Field:     field,
			Operation: "scale_by_percentage",
			Value:     percent,
			Unit:      "percentage",
		},

		AIGuidance: fmt.Sprintf(`
Locate the Kubernetes manifest defining resource "%s".
Inspect its "%s" field.
Scale it by %+g%%.
Preserve units (m, Mi, Gi).
If field is missing, create it.
Return only modified files.
`, resource, field, percent),
	}
}
