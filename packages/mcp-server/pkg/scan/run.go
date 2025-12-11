package scan

import "github.com/tanay13/costguard/packages/mcp-server/pkg/types"

func RunScan(req types.ScanRequest) (types.ScanResponse, error) {
	agg := DataPointAggregator(req.Metrics, req.ActualRequests)
	resp := BuildScanResponse(agg)
	return resp, nil
}
