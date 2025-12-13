package types

type FixOperation struct {
	Field     string  `json:"field"`
	Operation string  `json:"operation"` 
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
}

type FixAction struct {
	Provider    Provider     `json:"provider"`
	Resource    string       `json:"resource"`
	Intent      string       `json:"intent"`
	Description string       `json:"description"`
	Action      FixOperation `json:"action"`

	CurrentCost  float64 `json:"current_cost_usd"`
	ExpectedCost float64 `json:"expected_cost_usd"`
	Savings      float64 `json:"savings_usd"`

	Risk     string `json:"risk"`
	Priority int    `json:"priority"`

	AIGuidance  string   `json:"ai_guidance"`
	FilesToEdit []string `json:"files_to_edit,omitempty"`
}
