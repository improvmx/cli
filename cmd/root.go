package cmd

import (
	"fmt"
	"os"

	"github.com/improvmx/cli/internal/config"
	"github.com/improvmx/cli/internal/output"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "improvmx",
	Short: "ImprovMX CLI - Manage your email forwarding from the terminal",
	Long: `ImprovMX CLI lets you manage domains, aliases, SMTP credentials,
and view email logs directly from your terminal.

Get started by authenticating:
  improvmx auth login

Then add a domain:
  improvmx domain add example.com`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.Init)
	rootCmd.PersistentFlags().BoolVar(&output.JSONOutput, "json", false, "Output in JSON format")
}
