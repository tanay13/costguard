package types

type FixAction struct {
	Provider    Provider `json:"provider"`
	Resource    string   `json:"resource"`
	Intent      string   `json:"intent"` // e.g. rightsize_cpu_request
	Description string   `json:"description"`

	// AI operation fields
	Action FixOperation `json:"action"`

	CurrentCost  float64 `json:"current_cost_usd"`
	ExpectedCost float64 `json:"expected_cost_usd"`
	Savings      float64 `json:"savings_usd"`

	Risk     string `json:"risk"` // low/medium/high
	Priority int    `json:"priority"`

	AIGuidance  string   `json:"ai_guidance"` // Natural language instructions
	FilesToEdit []string `json:"files_to_edit,omitempty"`
}

type FixOperation struct {
	Field     string  `json:"field"`     // e.g. resources.requests.cpu
	Operation string  `json:"operation"` // e.g. scale_by_percentage
	Value     float64 `json:"value"`     // e.g. -40
	Unit      string  `json:"unit"`      // "percentage"
}
