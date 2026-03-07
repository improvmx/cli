package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/improvmx/cli/internal/api"
	"github.com/improvmx/cli/internal/config"
	"github.com/improvmx/cli/internal/output"
	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:     "domain",
	Aliases: []string{"domains"},
	Short:   "Manage domains",
}

var domainListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all domains",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())

		params := map[string]string{}
		if isActive, _ := cmd.Flags().GetBool("active"); isActive {
			params["is_active"] = "1"
		}

		resp, err := client.Get("/domains" + api.QueryEncode(params))
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if output.JSONOutput {
			fmt.Println(string(resp))
			return
		}

		var result struct {
			Domains []api.Domain `json:"domains"`
		}
		if err := parseResponse(resp, &result); err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if len(result.Domains) == 0 {
			fmt.Println("No domains found. Add one with: improvmx domain add <domain>")
			return
		}

		table := output.NewTable("DOMAIN", "ACTIVE", "ALIASES", "ADDED")
		for _, d := range result.Domains {
			active := "no"
			if d.Active {
				active = "yes"
			}
			added := time.Unix(d.Added, 0).Format("2006-01-02")
			table.AddRow(d.Name, active, strconv.Itoa(len(d.Aliases)), added)
		}
		table.Render()
	},
}

var domainGetCmd = &cobra.Command{
	Use:   "get <domain>",
	Short: "Get domain details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())
		domain := args[0]

		resp, err := client.Get("/domains/" + domain)
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if output.JSONOutput {
			fmt.Println(string(resp))
			return
		}

		var result struct {
			Domain api.Domain `json:"domain"`
		}
		if err := parseResponse(resp, &result); err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		d := result.Domain
		active := "no"
		if d.Active {
			active = "yes"
		}
		fmt.Printf("Domain:       %s\n", d.Name)
		fmt.Printf("Active:       %s\n", active)
		fmt.Printf("Notification: %s\n", d.NotificationEmail)
		fmt.Printf("Whitelabel:   %s\n", d.Whitelabel)
		fmt.Printf("Added:        %s\n", time.Unix(d.Added, 0).Format("2006-01-02 15:04:05"))

		if len(d.Aliases) > 0 {
			fmt.Printf("\nAliases (%d):\n", len(d.Aliases))
			table := output.NewTable("ALIAS", "FORWARDS TO")
			for _, a := range d.Aliases {
				alias := a.Alias
				if alias == "" {
					alias = "*"
				}
				table.AddRow(alias, a.Forward)
			}
			table.Render()
		}
	},
}

var domainAddCmd = &cobra.Command{
	Use:   "add <domain>",
	Short: "Add a new domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())
		domain := args[0]

		body := map[string]interface{}{
			"domain": domain,
		}

		if email, _ := cmd.Flags().GetString("notification-email"); email != "" {
			body["notification_email"] = email
		}
		if whitelabel, _ := cmd.Flags().GetString("whitelabel"); whitelabel != "" {
			body["whitelabel"] = whitelabel
		}

		resp, err := client.Post("/domains", body)
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if output.JSONOutput {
			fmt.Println(string(resp))
			return
		}

		output.Success(fmt.Sprintf("Domain %s added", domain))
		fmt.Println("\nNext steps:")
		fmt.Println("  1. Add MX records pointing to ImprovMX")
		fmt.Printf("  2. Run 'improvmx domain check %s' to verify DNS\n", domain)
	},
}

var domainDeleteCmd = &cobra.Command{
	Use:   "delete <domain>",
	Short: "Delete a domain",
	Aliases: []string{"rm", "remove"},
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())
		domain := args[0]

		_, err := client.Delete("/domains/" + domain)
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		output.Success(fmt.Sprintf("Domain %s deleted", domain))
	},
}

var domainCheckCmd = &cobra.Command{
	Use:   "check <domain>",
	Short: "Check domain DNS configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())
		domain := args[0]

		resp, err := client.Get("/domains/" + domain + "/check")
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if output.JSONOutput {
			fmt.Println(string(resp))
			return
		}

		var result struct {
			Valid  bool `json:"valid"`
			MX     struct {
				Valid bool `json:"valid"`
			} `json:"mx"`
			SPF struct {
				Valid bool `json:"valid"`
			} `json:"spf"`
		}
		if err := json.Unmarshal(resp, &result); err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		check := func(valid bool) string {
			if valid {
				return "OK"
			}
			return "MISSING"
		}

		fmt.Printf("Domain: %s\n", domain)
		fmt.Printf("MX:     %s\n", check(result.MX.Valid))
		fmt.Printf("SPF:    %s\n", check(result.SPF.Valid))
		fmt.Printf("Overall: %s\n", check(result.Valid))
	},
}

func init() {
	rootCmd.AddCommand(domainCmd)
	domainCmd.AddCommand(domainListCmd)
	domainCmd.AddCommand(domainGetCmd)
	domainCmd.AddCommand(domainAddCmd)
	domainCmd.AddCommand(domainDeleteCmd)
	domainCmd.AddCommand(domainCheckCmd)

	domainListCmd.Flags().Bool("active", false, "Only show active domains")
	domainAddCmd.Flags().String("notification-email", "", "Notification email address")
	domainAddCmd.Flags().String("whitelabel", "", "Whitelabel domain")
}
