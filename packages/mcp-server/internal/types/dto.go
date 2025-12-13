package types

type ScanResponse struct {
	Project  string              `json:"project"`
	Hotspots []AggregatedMetrics `json:"hotspots"`
}

type FixPlansAPIRequest struct {
	Metrics        []MetricCollection  `json:"metrics"`                   
	ActualRequests map[string]Requests `json:"actual_requests,omitempty"` 
	BudgetTarget   float64             `json:"budget_target_usd,omitempty"`
	AutoApprove    bool                `json:"auto_approve,omitempty"`
}
type FixPlanRequest struct {
	AggregatedMetrics []AggregatedMetrics `json:"aggregated_metrics"`
	BudgetTarget      float64             `json:"budget_target_usd,omitempty"`
	AutoApprove       bool                `json:"auto_approve,omitempty"`
}

type FixPlanResponse struct {
	TotalCurrentCost  float64     `json:"total_current_cost_usd"`
	TotalExpectedCost float64     `json:"total_expected_cost_usd"`
	TotalSavings      float64     `json:"total_savings_usd"`
	BudgetTarget      float64     `json:"budget_target_usd"`
	MeetsBudget       bool        `json:"meets_budget"`
	RequiresApproval  bool        `json:"requires_approval"`
	Actions           []FixAction `json:"actions"`
	Summary           string      `json:"summary"`
}
