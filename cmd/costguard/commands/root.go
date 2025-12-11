package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "costguard",
	Short: "CostGuard CLI — local cost analysis & optimization",
	Long: `
CostGuard CLI
--------------

Run cost analysis:
    costguard scan --metrics metrics.json

Generate and apply fix plans:
    costguard fix
`,
}

func Execute() error {
	return rootCmd.Execute()
}
