package provider

import (
	"fmt"
	"math"

	"github.com/tanay13/costguard/packages/mcp-server/internal/types"
)

const (
	minCPU = 50.0 // 50m minimum inferred request
	minMem = 0.1  // 0.1 GB minimum inferred request

	under = 0.50 // underutilized: p95 < 50% of inferred request
	over  = 0.90 // overutilized: p95 > 90% of inferred request
)

// ======== KUBERNETES FIXES ========

func generateK8sFixActions(metric types.AggregatedMetrics, priority int) []types.FixAction {
	actions := []types.FixAction{}

	cpu, hasCPU := metric.Metrics["cpu_milli"]
	mem, hasMem := metric.Metrics["memory_gb"]

	// 1. Infer requests
	inferredReqCPU := math.Max(cpu.P95*2, minCPU)
	inferredReqMem := math.Max(mem.P95*2, minMem)

	// 2. CPU underutilized → scale down
	if hasCPU && cpu.P95 < inferredReqCPU*under {
		actions = append(actions, makeK8sFixAction(
			metric.Resource,
			"rightsize_cpu_request",
			"resources.requests.cpu",
			-40, // reduce CPU by 40%
			"CPU is underutilized compared to inferred request",
			metric.CostUSD,
			priority,
		))
	}

	// 3. CPU overutilized → scale up
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

	// 4. Memory underutilized → scale down
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

	// 5. Memory overutilized → scale up
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
		// scaling up increases cost
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
