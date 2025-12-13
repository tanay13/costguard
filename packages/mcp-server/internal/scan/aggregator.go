package scan

import (
	"github.com/tanay13/costguard/packages/mcp-server/internal/provider"
	"github.com/tanay13/costguard/packages/mcp-server/internal/types"
)

func DataPointAggregator(points []types.MetricCollection) []types.AggregatedMetrics {
	out := make([]types.AggregatedMetrics, 0)

	
	groups := make(map[types.Provider]map[string][]types.MetricCollection)

	for _, p := range points {
		if groups[p.Provider] == nil {
			groups[p.Provider] = make(map[string][]types.MetricCollection)
		}
		name := p.Resource
		if name == "" {
			name = "unknown"
		}
		groups[p.Provider][name] = append(groups[p.Provider][name], p)
	}

	
	if k, ok := groups[types.ProviderKubernetes]; ok {
		kagg := provider.KubernetesMetricAggregator(k)
		out = append(out, kagg...)
	}

	
	return out
}
