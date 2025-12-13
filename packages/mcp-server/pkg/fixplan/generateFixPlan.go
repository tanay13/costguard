package fixplan

import (
	"fmt"

	"github.com/tanay13/costguard/packages/mcp-server/pkg/provider/kubernetes"
	"github.com/tanay13/costguard/packages/mcp-server/pkg/types"
)

func GenerateFixPlan(req types.FixPlanRequest) types.FixPlanResponse {
	totalCurrent := 0.0
	totalOptimal := 0.0
	actions := []types.FixAction{}

	for _, agg := range req.AggregatedMetrics {

		totalCurrent += agg.CostCurrentUSD
		totalOptimal += agg.CostOptimalUSD

		switch agg.Provider {
		case types.ProviderKubernetes:
			kActions := kubernetes.GenerateK8sFixActions(agg)

			for i := range kActions {

				if len(kActions) > 1 {
					kActions[i].EstimatedSavingsUSD = agg.CostSavingsUSD / float64(len(kActions))
				} else {
					kActions[i].EstimatedSavingsUSD = agg.CostSavingsUSD
				}
			}
			actions = append(actions, kActions...)
		}

	}

	totalSavings := totalCurrent - totalOptimal

	summary := fmt.Sprintf(
		"Total cost: $%.2f → optimized: $%.2f (savings: $%.2f)",
		totalCurrent, totalOptimal, totalSavings,
	)

	return types.FixPlanResponse{
		TotalCurrentCost: totalCurrent,
		TotalOptimalCost: totalOptimal,
		TotalSavings:     totalSavings,
		BudgetTarget:     req.BudgetTarget,
		MeetsBudget:      req.BudgetTarget > 0 && totalOptimal <= req.BudgetTarget,
		RequiresApproval: !req.AutoApprove,
		Actions:          actions,
		Summary:          summary,
	}
}
