package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/tanay13/costguard/packages/mcp-server/pkg/fixplan"
	"github.com/tanay13/costguard/packages/mcp-server/pkg/types"
	"github.com/tanay13/costguard/packages/mcp-server/pkg/utils"
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Generate fix plans and optionally apply them using Cline AI",
	RunE: func(cmd *cobra.Command, args []string) error {

		data, err := os.ReadFile(".costguard/scan.json")
		if err != nil {
			return fmt.Errorf("missing scan.json — run `costguard scan` first")
		}

		var scanRes types.ScanResponse
		if err := json.Unmarshal(data, &scanRes); err != nil {
			return fmt.Errorf("invalid scan.json: %v", err)
		}

		agg := utils.ConvertScanToAggregated(scanRes)

		req := types.FixPlanRequest{
			AggregatedMetrics: agg,
			BudgetTarget:      0,
			AutoApprove:       false,
		}

		plan := fixplan.GenerateFixPlan(req)

		printFixPlan(plan)

		fmt.Print("\nApply fixes using Cline? (y/n): ")
		var choice string
		fmt.Scanln(&choice)

		if choice != "y" {
			color.Yellow("\nNo fixes applied.")
			return nil
		}

		for _, action := range plan.Actions {
			color.Cyan("\nApplying fix for %s (%s)…", action.Resource, action.Intent)
			color.Yellow(action.AIGuidance)

			cmd := exec.Command("cline", "ai", action.AIGuidance)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}

		color.Green("\n✔ Fixes applied successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
