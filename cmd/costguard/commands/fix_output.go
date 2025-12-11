package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/tanay13/costguard/packages/mcp-server/pkg/types"
)

func printFixPlan(plan types.FixPlanResponse) {
	divider := "────────────────────────────────────────────────────────────"

	color.Cyan("\n%s\nCOSTGUARD — FIX PLAN\n%s\n", divider, divider)

	fmt.Printf("\n💰 Current Monthly Cost: $%.2f\n", plan.TotalCurrentCost)
	fmt.Printf("🎯 Optimal Cost:         $%.2f\n", plan.TotalOptimalCost)
	fmt.Printf("💡 Potential Savings:    $%.2f\n", plan.TotalSavings)
	if plan.MeetsBudget {
		color.Green("✔ Meets target budget\n")
	} else {
		color.Red("✘ Does not meet target budget\n")
	}

	color.Yellow("\nActions Recommended: %d\n", len(plan.Actions))

	for _, a := range plan.Actions {
		color.New(color.FgHiBlue, color.Bold).
			Printf("\n%s\nRESOURCE: %s (%s)\n%s\n",
				divider, a.Resource, a.Provider, divider)

		fmt.Printf("Intent:       %s\n", a.Intent)
		fmt.Printf("Description:  %s\n", a.Description)
		fmt.Printf("Fix Type:     %s %v %s\n",
			a.Action.Field, a.Action.Value, a.Action.Unit)
		fmt.Printf("Files:        %v\n", a.FilesToEdit)
		fmt.Printf("Savings:      $%.2f\n", a.Action.Value)

		color.Green("\nAI Guidance:\n")
		color.White("%s\n", a.AIGuidance)
	}
}
