# ImprovMX CLI

Manage your [ImprovMX](https://improvmx.com) email forwarding directly from the terminal. Add domains, create aliases, manage SMTP credentials, view logs, and more. Perfect for interacting with ImprovMX from Claude Code, OpenAI Codex, or OpenClaw!

## Installation

### From source

Requires [Go 1.21+](https://go.dev/dl/).

```bash
go install github.com/improvmx/cli@latest
```

### From release binaries

Download the latest binary for your platform from [Releases](https://github.com/improvmx/cli/releases).

## Authentication

Get your API key from your [ImprovMX dashboard](https://app.improvmx.com/api), then:

```bash
improvmx auth login
```

Or set the environment variable:

```bash
export IMPROVMX_API_KEY=your-api-key
```

## Usage

### Domains

```bash
improvmx domain list                  # List all domains
improvmx domain add example.com       # Add a domain
improvmx domain get example.com       # Get domain details
improvmx domain check example.com     # Check DNS configuration
improvmx domain delete example.com    # Delete a domain
```

### Aliases

```bash
improvmx alias list example.com                       # List aliases
improvmx alias add example.com hello user@gmail.com   # Add an alias
improvmx alias add example.com "*" user@gmail.com     # Add a catch-all
improvmx alias update example.com hello new@gmail.com # Update an alias
improvmx alias delete example.com hello               # Delete an alias
```

### Email Logs

```bash
improvmx logs example.com             # View recent email logs
```

### Rules

```bash
improvmx rule list example.com                                                                        # List rules
improvmx rule get example.com <rule-id>                                                               # Get rule details
improvmx rule add example.com --type alias --alias hello --forward user@gmail.com                     # Add alias rule
improvmx rule add example.com --type regex --regex ".*invoice.*" --scopes subject,body --forward user@gmail.com  # Add regex rule
improvmx rule add example.com --type cel --expression "subject.contains('finance')" --forward user@gmail.com     # Add CEL rule
improvmx rule update example.com <rule-id> --forward new@gmail.com                                    # Update a rule
improvmx rule delete example.com <rule-id>                                                            # Delete a rule
improvmx rule delete-all example.com                                                                  # Delete all rules
```

### SMTP Credentials

```bash
improvmx smtp list example.com                    # List credentials
improvmx smtp add example.com user password        # Add credentials
improvmx smtp delete example.com user              # Delete credentials
```

### Account

```bash
improvmx account                       # View account info
```

## Options

| Flag     | Description              |
|----------|--------------------------|
| `--json` | Output in JSON format    |
| `--help` | Show help for any command |

### JSON output

Every command supports `--json` for scripting and piping:

```bash
improvmx domain list --json | jq '.domains[].domain'
```

### Shell completions

```bash
# Bash
improvmx completion bash > /etc/bash_completion.d/improvmx

# Zsh
improvmx completion zsh > "${fpath[1]}/_improvmx"

# Fish
improvmx completion fish > ~/.config/fish/completions/improvmx.fish
```

## Configuration

Credentials are stored at:

- **macOS**: `~/Library/Application Support/improvmx/config.yaml`
- **Linux**: `~/.config/improvmx/config.yaml`

## License

MIT
