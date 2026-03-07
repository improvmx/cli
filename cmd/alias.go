package cmd

import (
	"fmt"
	"os"

	"github.com/improvmx/cli/internal/api"
	"github.com/improvmx/cli/internal/config"
	"github.com/improvmx/cli/internal/output"
	"github.com/spf13/cobra"
)

var aliasCmd = &cobra.Command{
	Use:     "alias",
	Aliases: []string{"aliases"},
	Short:   "Manage aliases for a domain",
}

var aliasListCmd = &cobra.Command{
	Use:     "list <domain>",
	Short:   "List aliases for a domain",
	Aliases: []string{"ls"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())
		domain := args[0]

		resp, err := client.Get("/domains/" + domain + "/aliases")
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if output.JSONOutput {
			fmt.Println(string(resp))
			return
		}

		var result struct {
			Aliases []api.Alias `json:"aliases"`
		}
		if err := parseResponse(resp, &result); err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if len(result.Aliases) == 0 {
			fmt.Printf("No aliases for %s. Add one with: improvmx alias add %s <alias> <forward-to>\n", domain, domain)
			return
		}

		table := output.NewTable("ALIAS", "FORWARDS TO", "ID")
		for _, a := range result.Aliases {
			alias := a.Alias
			if alias == "" {
				alias = "*"
			}
			table.AddRow(alias, a.Forward, fmt.Sprintf("%d", a.ID))
		}
		table.Render()
	},
}

var aliasAddCmd = &cobra.Command{
	Use:   "add <domain> <alias> <forward-to>",
	Short: "Add an alias to a domain",
	Long: `Add an alias to a domain.

Use "*" as the alias for a catch-all.

Examples:
  improvmx alias add example.com hello user@gmail.com
  improvmx alias add example.com "*" user@gmail.com`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())
		domain, alias, forward := args[0], args[1], args[2]

		body := map[string]string{
			"alias":   alias,
			"forward": forward,
		}

		resp, err := client.Post("/domains/"+domain+"/aliases", body)
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if output.JSONOutput {
			fmt.Println(string(resp))
			return
		}

		displayAlias := alias
		if alias == "" || alias == "*" {
			displayAlias = "*"
		}
		output.Success(fmt.Sprintf("Alias %s@%s -> %s created", displayAlias, domain, forward))
	},
}

var aliasUpdateCmd = &cobra.Command{
	Use:   "update <domain> <alias> <forward-to>",
	Short: "Update an alias",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())
		domain, alias, forward := args[0], args[1], args[2]

		body := map[string]string{
			"forward": forward,
		}

		resp, err := client.Put("/domains/"+domain+"/aliases/"+alias, body)
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if output.JSONOutput {
			fmt.Println(string(resp))
			return
		}

		output.Success(fmt.Sprintf("Alias %s@%s updated -> %s", alias, domain, forward))
	},
}

var aliasDeleteCmd = &cobra.Command{
	Use:     "delete <domain> <alias>",
	Short:   "Delete an alias",
	Aliases: []string{"rm", "remove"},
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())
		domain, alias := args[0], args[1]

		_, err := client.Delete("/domains/" + domain + "/aliases/" + alias)
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		output.Success(fmt.Sprintf("Alias %s@%s deleted", alias, domain))
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
	aliasCmd.AddCommand(aliasListCmd)
	aliasCmd.AddCommand(aliasAddCmd)
	aliasCmd.AddCommand(aliasUpdateCmd)
	aliasCmd.AddCommand(aliasDeleteCmd)
}
