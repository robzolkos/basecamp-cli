---
name: basecamp
description: Interact with Basecamp projects and card tables via CLI. List projects, browse Kanban boards, view cards with comments, and move cards between columns. Use when the user mentions Basecamp, card tables, Kanban boards, or wants to check project status.
---

# Basecamp CLI

Interact with Basecamp projects and card tables using the `basecamp` command. All output is JSON for easy parsing with `jq`.

## Prerequisites

The CLI must be installed and authenticated:

```bash
basecamp version  # Verify installation
basecamp projects  # Verify authentication
```

If not authenticated, the user needs to run `basecamp auth`.

## Commands

### List projects

```bash
basecamp projects
```

Output:
```json
[
  {
    "id": 12345678,
    "name": "Website Redesign",
    "description": "Project description here",
    "status": "active"
  },
  {
    "id": 23456789,
    "name": "Mobile App",
    "status": "active"
  }
]
```

### List boards in a project

```bash
basecamp boards <project_id>
```

Output:
```json
{
  "project_id": 12345678,
  "project_name": "Website Redesign",
  "board_id": 87654321,
  "board_title": "Development Tasks",
  "columns": [
    {"title": "Backlog", "cards_count": 12},
    {"title": "In Progress", "cards_count": 3},
    {"title": "Done", "cards_count": 45}
  ]
}
```

### List cards

```bash
basecamp cards <project_id> <board_id>
basecamp cards <project_id> <board_id> --column "In Progress"
```

Output:
```json
{
  "board_id": 87654321,
  "board_title": "Development Tasks",
  "columns": [
    {
      "column": "In Progress",
      "cards": [
        {"id": 44444444, "title": "Implement dark mode", "creator": "John Doe"},
        {"id": 55555555, "title": "Refactor authentication", "creator": "Jane Smith"}
      ]
    }
  ]
}
```

### View card details

```bash
basecamp card <project_id> <card_id>
basecamp card <project_id> <card_id> --comments
```

Output:
```json
{
  "id": 44444444,
  "title": "Implement dark mode",
  "creator": "John Doe",
  "created_at": "2025-01-15T09:30:00.000Z",
  "updated_at": "2025-01-20T14:22:00.000Z",
  "url": "https://3.basecamp.com/.../cards/44444444",
  "assignees": ["Jane Smith"],
  "description": "Add dark mode support to the application.",
  "comments": [
    {
      "id": 1,
      "author": "Jane Smith",
      "content": "Started on the color palette.",
      "created_at": "2025-01-16T10:00:00.000Z"
    },
    {
      "id": 2,
      "author": "John Doe",
      "content": "Looks great!",
      "created_at": "2025-01-17T09:15:00.000Z"
    }
  ]
}
```

### Move a card

```bash
basecamp move <project_id> <board_id> <card_id> --to "Done"
```

Output:
```json
{
  "status": "ok",
  "card_id": "44444444",
  "column": "Done",
  "message": "Card 44444444 moved to 'Done'"
}
```

## Workflow example

To check what's being worked on and move a completed card:

```bash
# Find the project
basecamp projects | jq '.[] | select(.status == "active")'

# Find the board
basecamp boards 12345678

# See cards in progress
basecamp cards 12345678 87654321 --column "Progress"

# View details of a specific card
basecamp card 12345678 44444444 --comments

# Move it to Done
basecamp move 12345678 87654321 44444444 --to "Done"
```

## Configuration

Config files follow XDG Base Directory spec:

| File | Path |
|------|------|
| Config | `~/.config/basecamp/config.json` |
| Token | `~/.local/share/basecamp/token.json` |

To set up from scratch:
```bash
basecamp init   # Configure client_id, client_secret, account_id
basecamp auth   # Authenticate via OAuth
```

## Project-specific config

Create `.basecamp.yml` in your project directory to set a default project:

```yaml
project_id: 12345678
```

The CLI searches current directory and parents, so you can:
- Put it in your repo root for per-project defaults
- Omit project_id from commands when in that directory

```bash
# With .basecamp.yml in current or parent directory:
basecamp boards              # no project_id needed
basecamp cards 87654321      # just board_id
basecamp card 44444444       # just card_id
```

## Error handling

Errors are returned as JSON on stderr:
```json
{"error": "not authenticated, run 'basecamp auth' first"}
```

If you see "not authenticated" or "token expired", the user needs to run:
```bash
basecamp auth
```
