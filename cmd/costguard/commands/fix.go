package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/tanay13/costguard/packages/mcp-server/pkg/dashboard"
	"github.com/tanay13/costguard/packages/mcp-server/pkg/fixplan"
	"github.com/tanay13/costguard/packages/mcp-server/pkg/github"
	"github.com/tanay13/costguard/packages/mcp-server/pkg/types"
	"github.com/tanay13/costguard/packages/mcp-server/pkg/utils"
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Generate fix plans and optionally apply them automatically",
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

		// Get dashboard URL from flag or environment variable
		dashboardURL, _ := cmd.Flags().GetString("dashboard-url")
		if dashboardURL == "" {
			dashboardURL = os.Getenv("COSTGUARD_DASHBOARD_URL")
		}

		// Get base branch from flag or default to "main"
		baseBranch, _ := cmd.Flags().GetString("base")
		if baseBranch == "" {
			baseBranch = "main"
		}

		fmt.Print("\nApply fixes and create PR? (y/n): ")
		var choice string
		fmt.Scanln(&choice)

		if choice != "y" {
			color.Yellow("\nNo fixes applied.")
			// Still send to dashboard if URL is provided (for tracking)
			if dashboardURL != "" {
				sendToDashboard(dashboardURL, scanRes, plan, "", 0, false)
			}
			return nil
		}

		// Create decision summary (all actions will be applied)
		decisionSummary := createDecisionSummary(plan, true)

		// Get repo info for PR creation
		repoOwner, repoName, _, err := dashboard.GetRepoInfo()
		if err != nil {
			color.Yellow("⚠️  Could not determine repo info: %v", err)
			repoOwner = "unknown"
			repoName = "unknown"
		}

		// Apply fixes and create PR using the github package
		// CreatePR will apply fixes internally
		color.Cyan("\n🔧 Applying fixes and creating Pull Request...")
		prConfig := github.PRConfig{
			BaseBranch: baseBranch,
			RepoOwner:  repoOwner,
			RepoName:   repoName,
		}

		prResult, err := github.CreatePR(prConfig, plan.Actions, decisionSummary)
		if err != nil {
			return fmt.Errorf("failed to create PR: %w", err)
		}

		if prResult.Success {
			color.Green("\n🎉 PR created successfully!")
			color.Green("PR: %s", prResult.PRURL)
			color.Green("Branch: %s", prResult.BranchName)
		}

		// Send update to dashboard
		if dashboardURL != "" {
			sendToDashboard(dashboardURL, scanRes, plan, prResult.PRURL, prResult.PRNumber, true)
		}

		return nil
	},
}

func init() {
	fixCmd.Flags().String("dashboard-url", "", "Dashboard URL to send scan results and decisions to (can also use COSTGUARD_DASHBOARD_URL env var)")
	fixCmd.Flags().String("base", "main", "Base branch for the PR")
	rootCmd.AddCommand(fixCmd)
}

func createDecisionSummary(plan types.FixPlanResponse, applied bool) types.AIDecisionSummary {
	decisions := make([]types.AIDecision, len(plan.Actions))
	totalSavings := 0.0
	actionsToApply := 0
	actionsSkipped := 0

	for i, action := range plan.Actions {
		decision := "apply"
		if !applied {
			decision = "skip"
			actionsSkipped++
		} else {
			actionsToApply++
			totalSavings += action.EstimatedSavingsUSD
		}

		decisions[i] = types.AIDecision{
			ActionID:            fmt.Sprintf("action-%d", i),
			Decision:            decision,
			Rationale:           action.Description,
			Priority:            i + 1,
			RiskLevel:           "low",
			EstimatedSavingsUSD: action.EstimatedSavingsUSD,
		}
	}

	summary := plan.Summary
	if !applied {
		summary = "Fix plan generated but not applied by user"
	}

	return types.AIDecisionSummary{
		TotalActions:    len(plan.Actions),
		ActionsToApply:  actionsToApply,
		ActionsDeferred: 0,
		ActionsSkipped:  actionsSkipped,
		TotalSavingsUSD: totalSavings,
		Decisions:       decisions,
		Summary:         summary,
	}
}

func sendToDashboard(dashboardURL string, scanRes types.ScanResponse, plan types.FixPlanResponse, prURL string, prNumber int, applied bool) {
	client := dashboard.NewClient(dashboardURL)
	decisionSummary := createDecisionSummary(plan, applied)

	err := client.SendUpdate(scanRes, decisionSummary, prURL, prNumber)
	if err != nil {
		color.Yellow("\n⚠️  Failed to send update to dashboard: %v", err)
	} else {
		color.Green("\n✔ Dashboard updated successfully")
	}
}
