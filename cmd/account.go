package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/improvmx/cli/internal/api"
	"github.com/improvmx/cli/internal/config"
	"github.com/improvmx/cli/internal/output"
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "View account information",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())

		resp, err := client.Get("/account")
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if output.JSONOutput {
			fmt.Println(string(resp))
			return
		}

		var result struct {
			Account api.Account `json:"account"`
		}
		if err := parseResponse(resp, &result); err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		a := result.Account
		fmt.Printf("Email:        %s\n", a.Email)
		fmt.Printf("Plan:         %s\n", a.Plan.Display)
		fmt.Printf("Premium:      %v\n", a.Premium)
		fmt.Printf("\nLimits:\n")
		fmt.Printf("  Domains:    %s\n", limitDisplay(a.Limits.Domains))
		fmt.Printf("  Aliases:    %s\n", limitDisplay(a.Limits.Aliases))
		fmt.Printf("  Daily Quota: %s\n", limitDisplay(a.Limits.DailyQuota))
	},
}

func limitDisplay(n int) string {
	if n == 0 {
		return "unlimited"
	}
	return strconv.Itoa(n)
}

func init() {
	rootCmd.AddCommand(accountCmd)
}
