package scan

import (
	"github.com/tanay13/costguard/packages/mcp-server/internal/provider"
	"github.com/tanay13/costguard/packages/mcp-server/internal/types"
)

func DataPointAggregator(metricCollection []types.MetricCollection) []types.AggregatedMetrics {
	aggregatedMetrics := make([]types.AggregatedMetrics, 0)

	providerMetrics := make(map[types.Provider]map[string][]types.MetricCollection)

	for _, c := range metricCollection {
		if providerMetrics[c.Provider] == nil {
			providerMetrics[c.Provider] = make(map[string][]types.MetricCollection)
		}

		resourceName := c.Resource
		if resourceName == "" {
			resourceName = "unknown"
		}

		providerMetrics[c.Provider][resourceName] =
			append(providerMetrics[c.Provider][resourceName], c)
	}

	if k8sGroups, ok := providerMetrics[types.ProviderKubernetes]; ok {
		result := provider.KubernetesMetricAggregator(k8sGroups)
		aggregatedMetrics = append(aggregatedMetrics, result...)
	}

	// For Future providers similar code can be used:
	// if lambdaGroups, ok := providerMetrics[types.ProviderAWSLambda]; ok {
	//     result := provider.LambdaMetricAggregator(lambdaGroups)
	//     aggregatedMetrics = append(aggregatedMetrics, result...)
	// }

	return aggregatedMetrics
}
