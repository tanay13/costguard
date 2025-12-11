package kubernetes

import (
	"fmt"

	"github.com/tanay13/costguard/packages/mcp-server/internal/types"
)

func GenerateFixPlan(req types.FixPlanRequest) types.FixPlanResponse {
	actions := []types.FixAction{}
	totalCurrent := 0.0

	for i, metric := range req.AggregatedMetrics {
		totalCurrent += metric.CostCurrentUSD

		switch metric.Provider {
		case types.ProviderKubernetes:
			actions = append(actions, generateK8sFixActions(metric, i+1)...)
		}
	}

	expected := 0.0
	for _, a := range actions {
		expected += a.ExpectedCost
	}

	savings := totalCurrent - expected
	meetsBudget := expected <= req.BudgetTarget
	needsApproval := !req.AutoApprove || savings > 50

	summary := fmt.Sprintf("Generated %d optimization actions with estimated savings $%.2f/mo", len(actions), savings)

	return types.FixPlanResponse{
		TotalCurrentCost:  totalCurrent,
		TotalExpectedCost: expected,
		TotalSavings:      savings,
		BudgetTarget:      req.BudgetTarget,
		MeetsBudget:       meetsBudget,
		RequiresApproval:  needsApproval,
		Actions:           actions,
		Summary:           summary,
	}
}
