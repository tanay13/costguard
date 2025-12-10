package scan

import (
	"github.com/tanay13/costguard/packages/mcp-server/internal/provider"
	"github.com/tanay13/costguard/packages/mcp-server/internal/types"
)

func DataPointAggregator(metricCollection []types.MetricCollection) []types.AggregatedMetrics {
	aggregatedMetrics := make([]types.AggregatedMetrics, 0)

	providerMetrics := make(map[types.Provider]interface{}, 0)

	for _, c := range metricCollection {
		switch c.Provider {
		case types.ProviderKubernetes:
			metricsSlice, _ := providerMetrics[types.ProviderKubernetes].([]types.K8sResourceMetrics)
			metricsSlice = append(metricsSlice, c.Metrics.K8sResourceMetrics)
			providerMetrics[types.ProviderKubernetes] = metricsSlice
		}
	}

	if k8sMetrics, ok := providerMetrics[types.ProviderKubernetes].([]types.K8sResourceMetrics); ok {
		result := provider.KubernetesMetricAggregator(k8sMetrics)
		aggregatedMetrics = append(aggregatedMetrics, result...)
	}

	return aggregatedMetrics
}
