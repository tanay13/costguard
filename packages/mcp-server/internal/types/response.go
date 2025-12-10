package types

type ScanResponse struct {
	Project  string              `json:"project"`
	Hotspots []AggregatedMetrics `json:"hotspots"`
}
