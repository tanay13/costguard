package kubernetes

import (
	"fmt"
	"math"

	"github.com/tanay13/costguard/packages/mcp-server/internal/types"
)

func GenerateK8sFixActions(agg types.AggregatedMetrics) []types.FixAction {
	out := []types.FixAction{}

	reqCPU := agg.RequestedCpuMilli
	reqMem := agg.RequestedMemoryGB

	optCPU := agg.OptimalCpuMilli
	optMem := agg.OptimalMemoryGB

	cpuPercent := ((optCPU - reqCPU) / reqCPU) * 100
	if math.Abs(cpuPercent) > 5 {
		out = append(out, types.FixAction{
			Provider: types.ProviderKubernetes,
			Resource: agg.Resource,
			Intent:   "rightsize_cpu_request",
			Description: fmt.Sprintf(
				"CPU request %.0fm → %.0fm (%.1f%% change)",
				reqCPU, optCPU, cpuPercent,
			),
			Action: types.FixOperation{
				Field:     "resources.requests.cpu",
				Operation: "set_to",
				Value:     optCPU,
				Unit:      "m",
			},
			AIGuidance: fmt.Sprintf(
				"Update the Kubernetes manifest for '%s'. Set CPU request to %.0fm.",
				agg.Resource, optCPU,
			),
		})
	}

	memPercent := ((optMem - reqMem) / reqMem) * 100
	if math.Abs(memPercent) > 5 {
		out = append(out, types.FixAction{
			Provider: types.ProviderKubernetes,
			Resource: agg.Resource,
			Intent:   "rightsize_memory_request",
			Description: fmt.Sprintf(
				"Memory request %.2fGB → %.2fGB (%.1f%% change)",
				reqMem, optMem, memPercent,
			),
			Action: types.FixOperation{
				Field:     "resources.requests.memory",
				Operation: "set_to",
				Value:     optMem,
				Unit:      "GB",
			},
			AIGuidance: fmt.Sprintf(
				"Update the Kubernetes manifest for '%s'. Set memory request to %.2fGB.",
				agg.Resource, optMem,
			),
		})
	}

	return out
}
