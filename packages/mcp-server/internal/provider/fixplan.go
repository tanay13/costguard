package provider

import (
	"fmt"

	"github.com/tanay13/costguard/packages/mcp-server/internal/types"
)

func GenerateFixPlan(req types.FixPlanRequest) types.FixPlanResponse {
	actions := []types.FixAction{}
	totalCurrent := 0.0

	for i, agg := range req.AggregatedMetrics {
		totalCurrent += agg.CostUSD

		if agg.Provider == types.ProviderKubernetes {
			acts := generateK8sFixActions(agg, i+1)
			actions = append(actions, acts...)
		}
	}

	expected := 0.0
	for _, a := range actions {
		expected += a.ExpectedCost
	}
	savings := totalCurrent - expected
	meets := expected <= req.BudgetTarget
	needApproval := !req.AutoApprove || savings > 50.0

	summary := fmt.Sprintf("Generated %d actions, estimated savings $%.2f/mo", len(actions), savings)

	return types.FixPlanResponse{
		TotalCurrentCost:  totalCurrent,
		TotalExpectedCost: expected,
		TotalSavings:      savings,
		BudgetTarget:      req.BudgetTarget,
		MeetsBudget:       meets,
		RequiresApproval:  needApproval,
		Actions:           actions,
		Summary:           summary,
	}
}
