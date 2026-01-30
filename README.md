# Basecamp CLI

A command-line interface for Basecamp written in Go with zero external dependencies.

## Installation

```bash
go install github.com/rzolkos/basecamp-cli/cmd/basecamp@latest
```

Or build from source:

```bash
make build
```

## Setup

1. Create a Basecamp OAuth app at https://launchpad.37signals.com/integrations
2. Run `basecamp init` to configure credentials
3. Run `basecamp auth` to authenticate

Configuration files (XDG Base Directory):
- `~/.config/basecamp/config.json` - client credentials
- `~/.local/share/basecamp/token.json` - OAuth token

## Usage

```bash
# List all projects
basecamp projects

# List card tables in a project
basecamp boards <project_id>

# List cards in a board
basecamp cards <project_id> <board_id>

# Filter cards by column
basecamp cards <project_id> <board_id> --column "In Progress"

# View card details
basecamp card <project_id> <card_id>

# View card with comments
basecamp card <project_id> <card_id> --comments

# Move a card to a different column
basecamp move <project_id> <board_id> <card_id> --to "Done"
```

## Project-specific config

Create `.basecamp.yml` in your project directory to set a default project_id:

```yaml
project_id: 12345678
```

Then omit project_id from commands:

```bash
basecamp boards              # uses project_id from .basecamp.yml
basecamp cards 87654321      # just need board_id
basecamp card 44444444       # just need card_id
```

The CLI searches current directory and parent directories for `.basecamp.yml`.

## Output

All commands output JSON for easy parsing with `jq`:

```bash
basecamp projects | jq '.[] | select(.status == "active") | .name'
```

Errors are output as JSON to stderr:

```json
{"error": "not authenticated, run 'basecamp auth' first"}
```

## Development

```bash
# Run tests
make test

# Build
make build

# Install to $GOPATH/bin
make install
```
