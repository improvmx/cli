package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/improvmx/cli/internal/api"
	"github.com/improvmx/cli/internal/config"
	"github.com/improvmx/cli/internal/output"
	"github.com/spf13/cobra"
)

var ruleCmd = &cobra.Command{
	Use:     "rule",
	Aliases: []string{"rules"},
	Short:   "Manage rules for a domain",
}

var ruleListCmd = &cobra.Command{
	Use:     "list <domain>",
	Short:   "List rules for a domain",
	Aliases: []string{"ls"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())
		domain := args[0]

		params := map[string]string{}
		if search, _ := cmd.Flags().GetString("search"); search != "" {
			params["search"] = search
		}
		if page, _ := cmd.Flags().GetInt("page"); page > 0 {
			params["page"] = strconv.Itoa(page)
		}

		resp, err := client.Get("/domains/" + domain + "/rules" + api.QueryEncode(params))
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if output.JSONOutput {
			fmt.Println(string(resp))
			return
		}

		var result struct {
			Rules []api.Rule `json:"rules"`
		}
		if err := parseResponse(resp, &result); err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if len(result.Rules) == 0 {
			fmt.Printf("No rules for %s. Add one with: improvmx rule add %s --type alias --alias hello --forward user@gmail.com\n", domain, domain)
			return
		}

		table := output.NewTable("ID", "TYPE", "ACTIVE", "RANK", "CONFIG", "CREATED")
		for _, r := range result.Rules {
			active := "no"
			if r.Active {
				active = "yes"
			}
			created := time.Unix(r.Created, 0).Format("2006-01-02")
			table.AddRow(r.ID, r.Type, active, fmt.Sprintf("%.1f", r.Rank), formatConfig(r), created)
		}
		table.Render()
	},
}

var ruleGetCmd = &cobra.Command{
	Use:   "get <domain> <rule-id>",
	Short: "Get rule details",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())
		domain, ruleID := args[0], args[1]

		resp, err := client.Get("/domains/" + domain + "/rules/" + ruleID)
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if output.JSONOutput {
			fmt.Println(string(resp))
			return
		}

		var r api.Rule
		if err := parseResponse(resp, &r); err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}
		active := "no"
		if r.Active {
			active = "yes"
		}
		fmt.Printf("ID:      %s\n", r.ID)
		fmt.Printf("Type:    %s\n", r.Type)
		fmt.Printf("Active:  %s\n", active)
		fmt.Printf("Rank:    %.1f\n", r.Rank)
		fmt.Printf("Created: %s\n", time.Unix(r.Created, 0).Format("2006-01-02 15:04:05"))
		fmt.Printf("Config:\n")

		configJSON, _ := json.MarshalIndent(r.Config, "  ", "  ")
		fmt.Printf("  %s\n", string(configJSON))
	},
}

var ruleAddCmd = &cobra.Command{
	Use:   "add <domain>",
	Short: "Add a rule to a domain",
	Long: `Add a rule to a domain. Three rule types are supported:

Alias rule (forward emails for a specific alias):
  improvmx rule add example.com --type alias --alias hello --forward user@gmail.com

Regex rule (match against email fields):
  improvmx rule add example.com --type regex --regex ".*invoice.*" --scopes subject,body --forward user@gmail.com

CEL rule (use CEL expressions):
  improvmx rule add example.com --type cel --expression "subject.contains('finance')" --forward user@gmail.com`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())
		domain := args[0]

		ruleType, _ := cmd.Flags().GetString("type")
		forward, _ := cmd.Flags().GetString("forward")
		rank, _ := cmd.Flags().GetFloat64("rank")
		active, _ := cmd.Flags().GetBool("active")

		if ruleType == "" {
			output.Error("--type is required (alias, regex, or cel)")
			os.Exit(1)
		}
		if forward == "" {
			output.Error("--forward is required")
			os.Exit(1)
		}

		cfg := map[string]interface{}{
			"forward": forward,
		}

		switch ruleType {
		case "alias":
			alias, _ := cmd.Flags().GetString("alias")
			if alias == "" {
				output.Error("--alias is required for alias rules")
				os.Exit(1)
			}
			cfg["alias"] = alias
		case "regex":
			regex, _ := cmd.Flags().GetString("regex")
			scopes, _ := cmd.Flags().GetString("scopes")
			if regex == "" {
				output.Error("--regex is required for regex rules")
				os.Exit(1)
			}
			if scopes == "" {
				output.Error("--scopes is required for regex rules (comma-separated: sender,recipient,subject,body)")
				os.Exit(1)
			}
			cfg["regex"] = regex
			cfg["scopes"] = strings.Split(scopes, ",")
		case "cel":
			expression, _ := cmd.Flags().GetString("expression")
			if expression == "" {
				output.Error("--expression is required for cel rules")
				os.Exit(1)
			}
			cfg["expression"] = expression
		default:
			output.Error(fmt.Sprintf("Unknown rule type: %s (must be alias, regex, or cel)", ruleType))
			os.Exit(1)
		}

		body := map[string]interface{}{
			"type":   ruleType,
			"config": cfg,
			"active": active,
		}
		if rank > 0 {
			body["rank"] = rank
		}

		resp, err := client.Post("/domains/"+domain+"/rules", body)
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if output.JSONOutput {
			fmt.Println(string(resp))
			return
		}

		var rule api.Rule
		if err := json.Unmarshal(resp, &rule); err == nil && rule.ID != "" {
			output.Success(fmt.Sprintf("Rule %s (%s) added to %s", rule.ID, ruleType, domain))
		} else {
			output.Success(fmt.Sprintf("Rule (%s) added to %s", ruleType, domain))
		}
	},
}

var ruleUpdateCmd = &cobra.Command{
	Use:   "update <domain> <rule-id>",
	Short: "Update a rule",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())
		domain, ruleID := args[0], args[1]

		body := map[string]interface{}{}

		if cmd.Flags().Changed("active") {
			active, _ := cmd.Flags().GetBool("active")
			body["active"] = active
		}
		if cmd.Flags().Changed("rank") {
			rank, _ := cmd.Flags().GetFloat64("rank")
			body["rank"] = rank
		}

		// Build config from any provided flags
		cfg := map[string]interface{}{}
		if forward, _ := cmd.Flags().GetString("forward"); forward != "" {
			cfg["forward"] = forward
		}
		if alias, _ := cmd.Flags().GetString("alias"); alias != "" {
			cfg["alias"] = alias
		}
		if regex, _ := cmd.Flags().GetString("regex"); regex != "" {
			cfg["regex"] = regex
		}
		if scopes, _ := cmd.Flags().GetString("scopes"); scopes != "" {
			cfg["scopes"] = strings.Split(scopes, ",")
		}
		if expression, _ := cmd.Flags().GetString("expression"); expression != "" {
			cfg["expression"] = expression
		}
		if len(cfg) > 0 {
			body["config"] = cfg
		}

		if len(body) == 0 {
			output.Error("No updates specified. Use flags like --forward, --active, --rank, etc.")
			os.Exit(1)
		}

		resp, err := client.Put("/domains/"+domain+"/rules/"+ruleID, body)
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		if output.JSONOutput {
			fmt.Println(string(resp))
			return
		}

		output.Success(fmt.Sprintf("Rule %s updated", ruleID))
	},
}

var ruleDeleteCmd = &cobra.Command{
	Use:     "delete <domain> <rule-id>",
	Short:   "Delete a rule",
	Aliases: []string{"rm", "remove"},
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())
		domain, ruleID := args[0], args[1]

		_, err := client.Delete("/domains/" + domain + "/rules/" + ruleID)
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		output.Success(fmt.Sprintf("Rule %s deleted", ruleID))
	},
}

var ruleDeleteAllCmd = &cobra.Command{
	Use:   "delete-all <domain>",
	Short: "Delete all rules for a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(config.GetAPIKey())
		domain := args[0]

		_, err := client.Delete("/domains/" + domain + "/rules-all")
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}

		output.Success(fmt.Sprintf("All rules deleted for %s", domain))
	},
}

func init() {
	rootCmd.AddCommand(ruleCmd)
	ruleCmd.AddCommand(ruleListCmd)
	ruleCmd.AddCommand(ruleGetCmd)
	ruleCmd.AddCommand(ruleAddCmd)
	ruleCmd.AddCommand(ruleUpdateCmd)
	ruleCmd.AddCommand(ruleDeleteCmd)
	ruleCmd.AddCommand(ruleDeleteAllCmd)

	// List flags
	ruleListCmd.Flags().String("search", "", "Filter rules by string match")
	ruleListCmd.Flags().Int("page", 0, "Page number")

	// Add flags
	ruleAddCmd.Flags().String("type", "", "Rule type: alias, regex, or cel (required)")
	ruleAddCmd.Flags().String("forward", "", "Forward destination email (required)")
	ruleAddCmd.Flags().String("alias", "", "Alias name (for alias rules)")
	ruleAddCmd.Flags().String("regex", "", "Regex pattern (for regex rules)")
	ruleAddCmd.Flags().String("scopes", "", "Comma-separated scopes: sender,recipient,subject,body (for regex rules)")
	ruleAddCmd.Flags().String("expression", "", "CEL expression (for cel rules)")
	ruleAddCmd.Flags().Float64("rank", 0, "Priority ranking")
	ruleAddCmd.Flags().Bool("active", true, "Whether rule is active")

	// Update flags
	ruleUpdateCmd.Flags().String("forward", "", "Forward destination email")
	ruleUpdateCmd.Flags().String("alias", "", "Alias name")
	ruleUpdateCmd.Flags().String("regex", "", "Regex pattern")
	ruleUpdateCmd.Flags().String("scopes", "", "Comma-separated scopes")
	ruleUpdateCmd.Flags().String("expression", "", "CEL expression")
	ruleUpdateCmd.Flags().Float64("rank", 0, "Priority ranking")
	ruleUpdateCmd.Flags().Bool("active", true, "Whether rule is active")
}

func shortID(id string) string {
	if len(id) > 8 {
		return id[:8]
	}
	return id
}

func formatConfig(r api.Rule) string {
	switch r.Type {
	case "alias":
		alias, _ := r.Config["alias"].(string)
		forward, _ := r.Config["forward"].(string)
		return fmt.Sprintf("%s -> %s", alias, forward)
	case "regex":
		regex, _ := r.Config["regex"].(string)
		forward, _ := r.Config["forward"].(string)
		return fmt.Sprintf("/%s/ -> %s", regex, forward)
	case "cel":
		expr, _ := r.Config["expression"].(string)
		forward, _ := r.Config["forward"].(string)
		if len(expr) > 30 {
			expr = expr[:27] + "..."
		}
		return fmt.Sprintf("%s -> %s", expr, forward)
	default:
		return fmt.Sprintf("%v", r.Config)
	}
}
