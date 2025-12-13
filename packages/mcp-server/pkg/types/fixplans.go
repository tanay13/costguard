package types

type FixPlansAPIRequest struct {
	Metrics        []MetricCollection  `json:"metrics"`
	ActualRequests map[string]Requests `json:"actual_requests,omitempty"`
	BudgetTarget   float64             `json:"budget_target_usd,omitempty"`
	AutoApprove    bool                `json:"auto_approve,omitempty"`
}

type FixAction struct {
	Provider    Provider `json:"provider"`
	Resource    string   `json:"resource"`
	Intent      string   `json:"intent"`
	Description string   `json:"description"`

	Action FixOperation `json:"action"`

	FilesToEdit []string `json:"files_to_edit,omitempty"`
	AIGuidance  string   `json:"ai_guidance"`

	EstimatedSavingsUSD float64 `json:"estimated_savings_usd"`
}

type FixOperation struct {
	Field     string  `json:"field"`
	Operation string  `json:"operation"`
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
}

type FixPlanRequest struct {
	AggregatedMetrics []AggregatedMetrics `json:"aggregated_metrics"`
	BudgetTarget      float64             `json:"budget_target_usd,omitempty"`
	AutoApprove       bool                `json:"auto_approve,omitempty"`
}

type FixPlanResponse struct {
	TotalCurrentCost float64     `json:"total_current_cost_usd"`
	TotalOptimalCost float64     `json:"total_optimal_cost_usd"`
	TotalSavings     float64     `json:"total_savings_usd"`
	BudgetTarget     float64     `json:"budget_target_usd"`
	MeetsBudget      bool        `json:"meets_budget"`
	RequiresApproval bool        `json:"requires_approval"`
	Actions          []FixAction `json:"actions"`
	Summary          string      `json:"summary"`
}

type AIDecision struct {
	ActionID            string  `json:"action_id"`
	Decision            string  `json:"decision"`
	Rationale           string  `json:"rationale"`
	Priority            int     `json:"priority"`
	RiskLevel           string  `json:"risk_level"`
	EstimatedSavingsUSD float64 `json:"estimated_savings_usd"`
}

type AIDecisionSummary struct {
	TotalActions    int          `json:"total_actions"`
	ActionsToApply  int          `json:"actions_to_apply"`
	ActionsDeferred int          `json:"actions_deferred"`
	ActionsSkipped  int          `json:"actions_skipped"`
	TotalSavingsUSD float64      `json:"total_savings_usd"`
	Decisions       []AIDecision `json:"decisions"`
	Summary         string       `json:"summary"`
}
