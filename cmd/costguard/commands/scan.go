package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tanay13/costguard/packages/mcp-server/pkg/scan"
	"github.com/tanay13/costguard/packages/mcp-server/pkg/types"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run cost analysis and display a detailed report",
	RunE: func(cmd *cobra.Command, args []string) error {
		metricsPath, _ := cmd.Flags().GetString("metrics")

		if metricsPath == "" {
			return fmt.Errorf("missing required flag: --metrics")
		}

		raw, err := os.ReadFile(metricsPath)
		if err != nil {
			return err
		}

		var req types.ScanRequest
		if err := json.Unmarshal(raw, &req); err != nil {
			return err
		}

		resp, err := scan.RunScan(req)
		if err != nil {
			return err
		}

		printScanReport(resp)

		out, _ := json.MarshalIndent(resp, "", "  ")
		os.MkdirAll(".costguard", 0755)
		os.WriteFile(".costguard/scan.json", out, 0644)

		color.Green("\n✔ Scan complete → .costguard/scan.json")

		return nil
	},
}

func init() {
	scanCmd.Flags().String("metrics", "", "Path to metrics JSON")
	rootCmd.AddCommand(scanCmd)
}
