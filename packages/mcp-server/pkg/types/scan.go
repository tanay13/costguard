package types

type ScanResourceCost struct {
	CurrentCostUSD      float64 `json:"current_cost_usd"`
	OptimalCostUSD      float64 `json:"optimal_cost_usd"`
	PotentialSavingsUSD float64 `json:"potential_savings_usd"`
	WastePercentage     float64 `json:"waste_percentage"`
}

type ScanResource struct {
	Provider  Provider              `json:"provider"`
	Resource  string                `json:"resource"`
	Usage     map[string]MetricStat `json:"usage"`
	Requested struct {
		CpuMilli float64 `json:"cpu_milli"`
		MemoryGB float64 `json:"memory_gb"`
	} `json:"requested"`
	Costs ScanResourceCost `json:"costs"`
}

type ScanSummary struct {
	TotalCurrentCostUSD      float64  `json:"total_current_cost_usd"`
	TotalOptimalCostUSD      float64  `json:"total_optimal_cost_usd"`
	TotalPotentialSavingsUSD float64  `json:"total_potential_savings_usd"`
	TopOffenders             []string `json:"top_offenders"`
}

type ScanResponse struct {
	Resources []ScanResource `json:"resources"`
	Summary   ScanSummary    `json:"summary"`
}
