# Basecamp CLI

A command-line interface for Basecamp written in Go with zero external dependencies.

## Installation

### Homebrew

```bash
brew install rzolkos/tap/basecamp
```

### AUR (Arch Linux)

```bash
yay -S basecamp-cli
```

### Go

```bash
go install github.com/rzolkos/basecamp-cli/cmd/basecamp@latest
```

### Build from source

```bash
make build
```

## Setup

### Prerequisites

During OAuth authentication, Basecamp redirects to your computer on **port 3002**. Your machine must be accessible via a URL for this to work. Use a service like:

- [Tailscale](https://tailscale.com/) - Recommended for persistent access
- [ngrok](https://ngrok.com/) - Quick setup for temporary access
- Any reverse proxy that exposes localhost:3002

### Registration

1. Start your tunnel service and note the public URL (e.g., `https://myhost.tailscale.ts.net`)

2. Run the registration helper to generate your OAuth app values:

```bash
basecamp register
```

This will ask for your application details and output the exact values to enter in the Basecamp registration form, including the correct redirect URI.

3. Visit https://launchpad.37signals.com/integrations and register your app using the generated values

4. Run `basecamp init` to configure your credentials (Client ID, Client Secret, and the same Redirect URI)

5. Run `basecamp auth` to authenticate (ensure your tunnel is running on port 3002)

### Configuration Files

Configuration follows XDG Base Directory specification:
- `~/.config/basecamp/config.json` - client credentials
- `~/.local/share/basecamp/token.json` - OAuth token

## Usage

### Card Tables

```bash
# List all projects
basecamp projects

# List card tables in a project
basecamp boards <project_id>

# List columns in a board
basecamp columns <project_id> <board_id>

# List cards in a board
basecamp cards <project_id> <board_id>

# Filter cards by column
basecamp cards <project_id> <board_id> --column "In Progress"

# View card details
basecamp card <project_id> <card_id>

# View card with comments
basecamp card <project_id> <card_id> --comments

# Create a card
basecamp card-create <project_id> <board_id> --column <column_id> --title "Card title"

# Update a card
basecamp card-update <project_id> <card_id> --title "New title" --content "Description"

# Move a card to a different column
basecamp move <project_id> <board_id> <card_id> --to "Done"
```

### Card Steps

```bash
# Steps are shown when viewing a card
basecamp card <project_id> <card_id>

# Create a step on a card
basecamp step-create <project_id> <card_id> --title "Step description"

# Create a step with due date and assignees
basecamp step-create <project_id> <card_id> --title "Task" --due 2026-02-01 --assignees "123,456"

# Update a step
basecamp step-update <project_id> <step_id> --title "Updated title"

# Complete a step
basecamp step-complete <project_id> <step_id>

# Uncomplete a step
basecamp step-uncomplete <project_id> <step_id>

# Reposition a step (0-indexed)
basecamp step-reposition <project_id> <card_id> <step_id> --position 0
```

### Todos

```bash
# List todo lists in a project
basecamp todolists <project_id>

# List todos in a todo list
basecamp todos <project_id> <todolist_id>

# List completed todos
basecamp todos <project_id> <todolist_id> --completed

# View a todo
basecamp todo <project_id> <todo_id>

# Create a todo
basecamp todo-create <project_id> <todolist_id> --content "Task description"

# Complete a todo
basecamp todo-complete <project_id> <todo_id>

# Uncomplete a todo
basecamp todo-uncomplete <project_id> <todo_id>

# Reposition a todo within its list (1-indexed)
basecamp todo-reposition <project_id> <todo_id> --position 1
```

### Todo Groups

```bash
# List groups within a todolist
basecamp todolist-groups <project_id> <todolist_id>

# View a group
basecamp todolist-group <project_id> <group_id>

# Create a group within a todolist
basecamp todolist-group-create <project_id> <todolist_id> --name "Group Name"

# Create a group with color
basecamp todolist-group-create <project_id> <todolist_id> --name "Group Name" --color green
```

### Messages

```bash
# List messages in a project
basecamp messages <project_id>

# View a message
basecamp message <project_id> <message_id>

# View a message with comments
basecamp message <project_id> <message_id> --comments

# Create a message
basecamp message-create <project_id> --subject "Subject" --content "Body"
```

### Comments

```bash
# Add a comment to any recording (card, message, todo, etc.)
basecamp comment-add <project_id> <recording_id> --content "Comment text"
```

### Documents

```bash
# List documents in a project
basecamp docs <project_id>

# View a document
basecamp doc <project_id> <doc_id>

# View a document with comments
basecamp doc <project_id> <doc_id> --comments

# Create a document
basecamp doc-create <project_id> --title "Title" --content "Content"
```

### Schedule

```bash
# List schedule entries
basecamp schedule <project_id>

# View an event
basecamp event <project_id> <entry_id>

# View an event with comments
basecamp event <project_id> <entry_id> --comments

# Create an event
basecamp event-create <project_id> --summary "Meeting" --starts-at "2026-02-01T10:00:00Z" --ends-at "2026-02-01T11:00:00Z"

# Create an all-day event
basecamp event-create <project_id> --summary "Holiday" --starts-at "2026-02-01" --ends-at "2026-02-01" --all-day
```

### Campfire

```bash
# List campfire messages
basecamp campfire <project_id>

# Post to campfire
basecamp campfire-post <project_id> --content "Hello team!"
```

### Search

```bash
# Search across all projects
basecamp search "query"

# Search with type filter
basecamp search "query" --type Todo

# Search within a project
basecamp search "query" --project <project_id>
```

### People

```bash
# List all people
basecamp people

# View a person's details
basecamp person <person_id>

# List pingable people
basecamp people-pingable

# List people on a project
basecamp people-project <project_id>

# View your own profile
basecamp my-profile

# Update project access (grant/revoke people)
basecamp project-access <project_id> --grant "123,456" --revoke "789"
```

### Automatic Check-ins (Questions)

```bash
# Get questionnaire info for a project
basecamp questionnaire <project_id>

# List check-in questions
basecamp questions <project_id>

# View a question
basecamp question <project_id> <question_id>

# View a question with comments
basecamp question <project_id> <question_id> --comments

# List answers to a question
basecamp question-answers <project_id> <question_id>

# View an answer
basecamp question-answer <project_id> <answer_id>

# View an answer with comments
basecamp question-answer <project_id> <answer_id> --comments
```

### Uploads

```bash
# Upload a file (returns attachable_sgid for use in other API calls)
basecamp upload /path/to/file.pdf

# List uploads in a vault (get vault_id from 'docs' command)
basecamp uploads <project_id> <vault_id>

# View an upload
basecamp upload-view <project_id> <upload_id>

# View an upload with comments
basecamp upload-view <project_id> <upload_id> --comments
```

### Recordings Management

```bash
# Archive any recording (todo, message, card, etc.)
basecamp archive <project_id> <recording_id>

# Unarchive a recording (set back to active)
basecamp unarchive <project_id> <recording_id>

# Trash a recording
basecamp trash <project_id> <recording_id>
```

### Message Types

```bash
# List message types (categories)
basecamp message-types <project_id>

# View a message type
basecamp message-type <project_id> <type_id>

# Create a message type
basecamp message-type-create <project_id> --name "Announcement" --icon "📢"

# Update a message type
basecamp message-type-update <project_id> <type_id> --name "Update" --icon "✅"

# Delete a message type
basecamp message-type-delete <project_id> <type_id>
```

### Activity Events

```bash
# List all events (across all projects)
basecamp events

# List events for a project
basecamp events-project <project_id>

# List events for a specific recording
basecamp events-recording <project_id> <recording_id>
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

## Agent Skills

Install the Basecamp skill for AI coding agents (Claude Code, OpenCode, and others):

```bash
npx skills add robzolkos/basecamp-cli
```

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
