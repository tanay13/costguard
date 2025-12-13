package ai

import (
	"fmt"
	"sort"

	"github.com/tanay13/costguard/packages/mcp-server/pkg/types"
)

func MakeDecisions(plan types.FixPlanResponse) types.AIDecisionSummary {
	decisions := []types.AIDecision{}

	actionSavings := make(map[int]float64)
	totalSavings := plan.TotalSavings
	numActions := len(plan.Actions)

	avgSavings := totalSavings / float64(numActions)

	for i, action := range plan.Actions {
		savings := avgSavings
		if action.EstimatedSavingsUSD > 0 {
			savings = action.EstimatedSavingsUSD
		}
		actionSavings[i] = savings
	}

	type actionWithIndex struct {
		index   int
		action  types.FixAction
		savings float64
	}

	actionsSorted := make([]actionWithIndex, 0, len(plan.Actions))
	for i, action := range plan.Actions {
		actionsSorted = append(actionsSorted, actionWithIndex{
			index:   i,
			action:  action,
			savings: actionSavings[i],
		})
	}

	sort.Slice(actionsSorted, func(i, j int) bool {
		return actionsSorted[i].savings > actionsSorted[j].savings
	})

	actionsToApply := 0
	actionsDeferred := 0
	actionsSkipped := 0
	totalSavingsToApply := 0.0

	for i, item := range actionsSorted {
		action := item.action
		savings := item.savings

		riskLevel := "low"
		decision := "apply"
		priority := 10 - i

		changePercent := 0.0
		if action.Action.Operation == "set_to" {

			if action.Intent == "rightsize_cpu_request" || action.Intent == "rightsize_memory_request" {
				changePercent = 20.0
			}
		}

		if changePercent > 50 {
			riskLevel = "high"
			decision = "defer"
			actionsDeferred++
		} else if changePercent > 30 {
			riskLevel = "medium"
			decision = "defer"
			actionsDeferred++
		} else if savings < 1.0 {
			decision = "skip"
			actionsSkipped++
			riskLevel = "low"
		} else {
			decision = "apply"
			actionsToApply++
			totalSavingsToApply += savings
		}

		rationale := fmt.Sprintf(
			"Savings: $%.2f/month, Risk: %s, Change: %.1f%%. %s",
			savings, riskLevel, changePercent,
			action.Description,
		)

		decisions = append(decisions, types.AIDecision{
			ActionID:            fmt.Sprintf("action-%d", item.index),
			Decision:            decision,
			Rationale:           rationale,
			Priority:            priority,
			RiskLevel:           riskLevel,
			EstimatedSavingsUSD: savings,
		})
	}

	summary := fmt.Sprintf(
		"AI Decision Summary: %d actions to apply (savings: $%.2f/month), %d deferred, %d skipped",
		actionsToApply, totalSavingsToApply, actionsDeferred, actionsSkipped,
	)

	return types.AIDecisionSummary{
		TotalActions:    len(plan.Actions),
		ActionsToApply:  actionsToApply,
		ActionsDeferred: actionsDeferred,
		ActionsSkipped:  actionsSkipped,
		TotalSavingsUSD: totalSavingsToApply,
		Decisions:       decisions,
		Summary:         summary,
	}
}
