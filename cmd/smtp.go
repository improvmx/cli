package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/improvmx/cli/internal/api"
	"github.com/improvmx/cli/internal/config"
	"github.com/improvmx/cli/internal/output"
	"github.com/spf13/cobra"
)

var smtpCmd = &cobra.Command{
	Use:   "smtp",
	Short: "Manage SMTP credentials",
}

var smtpListCmd = &cobra.Command{
	Use:     "list <domain>",
	Short:   "List SMTP credentials for a domain",
	Aliases: []string{"ls"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())
		domain := args[0]

		resp, err := client.Get("/domains/" + domain + "/credentials")
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if output.JSONOutput {
			fmt.Println(string(resp))
			return
		}

		var result struct {
			Credentials []api.SMTPCredential `json:"credentials"`
		}
		if err := parseResponse(resp, &result); err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if len(result.Credentials) == 0 {
			fmt.Printf("No SMTP credentials for %s\n", domain)
			return
		}

		table := output.NewTable("USERNAME", "USAGE", "CREATED")
		for _, c := range result.Credentials {
			created := time.Unix(c.Created/1000, 0).Format("2006-01-02")
			table.AddRow(c.Username, strconv.Itoa(c.Usage), created)
		}
		table.Render()
	},
}

var smtpAddCmd = &cobra.Command{
	Use:   "add <domain> <username> <password>",
	Short: "Add SMTP credentials",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())
		domain, username, password := args[0], args[1], args[2]

		body := map[string]string{
			"username": username,
			"password": password,
		}

		resp, err := client.Post("/domains/"+domain+"/credentials", body)
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if output.JSONOutput {
			fmt.Println(string(resp))
			return
		}

		output.Success(fmt.Sprintf("SMTP credentials created for %s@%s", username, domain))
	},
}

var smtpDeleteCmd = &cobra.Command{
	Use:     "delete <domain> <username>",
	Short:   "Delete SMTP credentials",
	Aliases: []string{"rm", "remove"},
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())
		domain, username := args[0], args[1]

		_, err := client.Delete("/domains/" + domain + "/credentials/" + username)
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		output.Success(fmt.Sprintf("SMTP credentials for %s@%s deleted", username, domain))
	},
}

func init() {
	rootCmd.AddCommand(smtpCmd)
	smtpCmd.AddCommand(smtpListCmd)
	smtpCmd.AddCommand(smtpAddCmd)
	smtpCmd.AddCommand(smtpDeleteCmd)
}
