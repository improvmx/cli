package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/improvmx/cli/internal/api"
	"github.com/improvmx/cli/internal/config"
	"github.com/improvmx/cli/internal/output"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs <domain>",
	Short: "View email logs for a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())
		domain := args[0]

		resp, err := client.Get("/domains/" + domain + "/logs")
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if output.JSONOutput {
			fmt.Println(string(resp))
			return
		}

		var result struct {
			Logs []api.LogEntry `json:"logs"`
		}
		if err := parseResponse(resp, &result); err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if len(result.Logs) == 0 {
			fmt.Printf("No logs found for %s\n", domain)
			return
		}

		table := output.NewTable("TIME", "FROM", "TO", "SUBJECT", "STATUS")
		for _, l := range result.Logs {
			t := time.Unix(l.Created, 0).Format("Jan 02 15:04")
			status := "pending"
			if len(l.Events) > 0 {
				last := l.Events[len(l.Events)-1]
				status = last.Status
			}
			subject := l.Subject
			if len(subject) > 40 {
				subject = subject[:37] + "..."
			}
			table.AddRow(t, l.Sender.Email, l.Recipient.Email, subject, status)
		}
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
