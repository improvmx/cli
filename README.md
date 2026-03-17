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

## Building

Requires [Go 1.23+](https://go.dev/dl/).

### Build for current platform

```bash
go build -o improvmx .
```

### Cross-compile for all platforms

```bash
GOOS=darwin  GOARCH=arm64 go build -o dist/improvmx-darwin-arm64 .
GOOS=darwin  GOARCH=amd64 go build -o dist/improvmx-darwin-amd64 .
GOOS=linux   GOARCH=amd64 go build -o dist/improvmx-linux-amd64 .
GOOS=linux   GOARCH=arm64 go build -o dist/improvmx-linux-arm64 .
GOOS=windows GOARCH=amd64 go build -o dist/improvmx-windows-amd64.exe .
```

## Signing & Notarization (macOS)

To distribute the macOS binaries without Gatekeeper warnings, you need to sign and notarize them with an Apple Developer ID.

### Prerequisites

- An [Apple Developer Program](https://developer.apple.com/programs/enroll/) membership (Organization, requires a DUNS number)
- A **Developer ID Application** certificate — create one at [Certificates, Identifiers & Profiles](https://developer.apple.com/account/resources/certificates/list):
  - Choose **Developer ID Application**
  - Select **G2 Sub-CA (Xcode 11.4.1 or later)** as the intermediary
  - Generate a CSR via **Keychain Access → Certificate Assistant → Request a Certificate From a Certificate Authority** (save to disk)
  - Upload the CSR and download the resulting `.cer` file
  - Double-click the `.cer` file to install it into your Keychain
- An **app-specific password** from [Apple ID account](https://appleid.apple.com/account/manage) (Sign-In and Security → App-Specific Passwords)

### Verify your certificate

```bash
security find-identity -v
```

You should see something like:

```
1) XXXXXXXX "Developer ID Application: ImprovMX Incorporated (2TMRXZB6JT)"
```

### Sign the binaries

```bash
codesign --sign "Developer ID Application: ImprovMX Incorporated (2TMRXZB6JT)" \
  --options runtime \
  --timestamp \
  dist/improvmx-darwin-arm64

codesign --sign "Developer ID Application: ImprovMX Incorporated (2TMRXZB6JT)" \
  --options runtime \
  --timestamp \
  dist/improvmx-darwin-amd64
```

### Verify signatures

```bash
codesign --verify --verbose dist/improvmx-darwin-arm64
codesign --verify --verbose dist/improvmx-darwin-amd64
```

### Store notarization credentials

Store your credentials in the Keychain to avoid passing them on the command line:

```bash
xcrun notarytool store-credentials "improvmx" \
  --apple-id "YOUR_APPLE_ID" \
  --team-id "2TMRXZB6JT" \
  --password "YOUR_APP_SPECIFIC_PASSWORD"
```

### Notarize the binaries

```bash
# Zip the binaries
ditto -c -k --keepParent dist/improvmx-darwin-arm64 dist/improvmx-darwin-arm64.zip
ditto -c -k --keepParent dist/improvmx-darwin-amd64 dist/improvmx-darwin-amd64.zip

# Submit for notarization
xcrun notarytool submit dist/improvmx-darwin-arm64.zip \
  --keychain-profile "improvmx" \
  --wait

xcrun notarytool submit dist/improvmx-darwin-amd64.zip \
  --keychain-profile "improvmx" \
  --wait
```

Notarization typically takes under 5 minutes. If a submission gets stuck, check [Apple's system status](https://developer.apple.com/system-status/) and resubmit.

### Check notarization status

```bash
# List all submissions
xcrun notarytool history --keychain-profile "improvmx"

# Get details for a specific submission
xcrun notarytool info SUBMISSION_ID --keychain-profile "improvmx"

# View the log if rejected
xcrun notarytool log SUBMISSION_ID --keychain-profile "improvmx"
```

## Releasing

After building, signing, and notarizing, create a GitHub release:

```bash
# Tag the release
git tag v0.2.0
git push origin v0.2.0

# Create the release with all binaries
gh release create v0.2.0 \
  dist/improvmx-darwin-arm64 \
  dist/improvmx-darwin-amd64 \
  dist/improvmx-linux-amd64 \
  dist/improvmx-linux-arm64 \
  dist/improvmx-windows-amd64.exe \
  --title "v0.2.0" \
  --notes "Release notes here"
```

Make sure to upload the **signed and notarized** macOS binaries, not the unsigned ones.

## License

MIT
