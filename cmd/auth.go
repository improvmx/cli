package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/improvmx/cli/internal/api"
	"github.com/improvmx/cli/internal/config"
	"github.com/improvmx/cli/internal/output"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with your ImprovMX API key",
	Long: `Authenticate with your ImprovMX API key.

You can find your API key at https://app.improvmx.com/api`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, _ := cmd.Flags().GetString("api-key")

		if apiKey == "" {
			fmt.Println("Get your API key from https://app.improvmx.com/api")
			fmt.Print("Enter your API key: ")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				output.Error("Failed to read input")
				os.Exit(1)
			}
			apiKey = strings.TrimSpace(input)
		}

		if apiKey == "" {
			output.Error("API key cannot be empty")
			os.Exit(1)
		}

		// Verify the key works
		client := api.NewClient(apiKey)
		_, err := client.Get("/account")
		if err != nil {
			output.Error(fmt.Sprintf("Invalid API key: %v", err))
			os.Exit(1)
		}

		if err := config.SaveAPIKey(apiKey); err != nil {
			output.Error(fmt.Sprintf("Failed to save API key: %v", err))
			os.Exit(1)
		}

		output.Success("Authenticated successfully")
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current authentication status",
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := config.GetAPIKey()
		client := api.NewClient(apiKey)

		resp, err := client.Get("/account")
		if err != nil {
			output.Error(fmt.Sprintf("Not authenticated: %v", err))
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

		fmt.Printf("Logged in as: %s\n", result.Account.Email)
		fmt.Printf("Plan: %s\n", result.Account.Plan.Display)
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored credentials",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.SaveAPIKey(""); err != nil {
			output.Error(fmt.Sprintf("Failed to remove credentials: %v", err))
			os.Exit(1)
		}
		output.Success("Logged out successfully")
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authLogoutCmd)

	authLoginCmd.Flags().String("api-key", "", "API key (or enter interactively)")
}
