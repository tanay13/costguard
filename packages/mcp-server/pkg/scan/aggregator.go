package scan

import (
	"github.com/tanay13/costguard/packages/mcp-server/pkg/provider/kubernetes"
	"github.com/tanay13/costguard/packages/mcp-server/pkg/types"
)

func DataPointAggregator(
	points []types.MetricCollection,
	actualRequests map[string]types.Requests,
) []types.AggregatedMetrics {

	grouped := make(map[types.Provider]map[string][]types.MetricCollection)

	for _, p := range points {
		if grouped[p.Provider] == nil {
			grouped[p.Provider] = make(map[string][]types.MetricCollection)
		}
		grouped[p.Provider][p.Resource] = append(grouped[p.Provider][p.Resource], p)
	}

	out := []types.AggregatedMetrics{}

	if km, ok := grouped[types.ProviderKubernetes]; ok {
		out = append(out, kubernetes.Aggregate(km, actualRequests)...)
	}

	// TODO: lambda, vercel etc.

	return out
}
