package commands

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/tanay13/costguard/packages/mcp-server/pkg/types"
)

func printScanReport(resp types.ScanResponse) {
	
	cyan := color.New(color.FgCyan).Add(color.Bold)
	cyan.Println("\n┌──────────────────────────────────────────────────────────┐")
	cyan.Println("│                    COSTGUARD — SCAN REPORT               │")
	cyan.Println("└──────────────────────────────────────────────────────────┘")

	fmt.Println()

	
	fmt.Printf("📦 Resources Analyzed: %d\n", len(resp.Resources))
	fmt.Printf("💰 Total Current Cost:   $%.2f / month\n", resp.Summary.TotalCurrentCostUSD)
	fmt.Printf("🎯 Optimal Cost:         $%.2f / month\n", resp.Summary.TotalOptimalCostUSD)
	fmt.Printf("💡 Potential Savings:    $%.2f / month (%.1f%%)\n\n",
		resp.Summary.TotalPotentialSavingsUSD,
		(resp.Summary.TotalPotentialSavingsUSD/resp.Summary.TotalCurrentCostUSD)*100,
	)

	
	if len(resp.Summary.TopOffenders) > 0 {
		fmt.Println("🔥 Top Offenders:")
		for i, name := range resp.Summary.TopOffenders {
			fmt.Printf("   %d. %s\n", i+1, name)
		}
		fmt.Println()
	}

	
	for _, r := range resp.Resources {
		printResourceDetail(r)
	}
}

func printResourceDetail(r types.ScanResource) {
	divider := strings.Repeat("─", 60)
	color.New(color.FgHiBlue, color.Bold).Printf("\n%s\nRESOURCE: %s (%s)\n%s\n",
		divider, r.Resource, r.Provider, divider)

	
	fmt.Printf("CPU Usage (milli):\n")
	fmt.Printf("   P50:      %.0fm\n", r.Usage["cpu_milli"].P50)
	fmt.Printf("   P95:      %.0fm\n", r.Usage["cpu_milli"].P95)
	fmt.Printf("   Average:  %.0fm\n", r.Usage["cpu_milli"].Avg)
	fmt.Printf("   Requested: %.0fm\n", r.Requested.CpuMilli)
	fmt.Printf("   Waste:     %.1f%%\n\n", r.Costs.WastePercentage)

	
	fmt.Printf("Memory Usage (GB):\n")
	fmt.Printf("   P50:      %.2f GB\n", r.Usage["memory_gb"].P50)
	fmt.Printf("   P95:      %.2f GB\n", r.Usage["memory_gb"].P95)
	fmt.Printf("   Requested: %.2f GB\n", r.Requested.MemoryGB)
	fmt.Printf("   Waste:     %.1f%%\n\n", r.Costs.WastePercentage)

	
	fmt.Printf("Cost:\n")
	fmt.Printf("   Current:  $%.2f\n", r.Costs.CurrentCostUSD)
	fmt.Printf("   Optimal:  $%.2f\n", r.Costs.OptimalCostUSD)
	fmt.Printf("   Savings:  $%.2f\n\n", r.Costs.PotentialSavingsUSD)
}
