package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/rzolkos/basecamp-cli/internal/config"
)

type Command interface {
	Run(args []string) error
}

var commands = map[string]func() Command{
	"init":     func() Command { return &InitCmd{} },
	"auth":     func() Command { return &AuthCmd{} },
	"projects": func() Command { return &ProjectsCmd{} },
	"boards":   func() Command { return &BoardsCmd{} },
	"cards":    func() Command { return &CardsCmd{} },
	"card":     func() Command { return &CardCmd{} },
	"move":     func() Command { return &MoveCmd{} },
}

func Execute(args []string, version string) {
	if len(args) < 1 {
		printHelp(version)
		os.Exit(1)
	}

	cmd := args[0]

	if cmd == "version" || cmd == "--version" || cmd == "-v" {
		fmt.Println(version)
		return
	}

	if cmd == "help" || cmd == "--help" || cmd == "-h" {
		printHelp(version)
		return
	}

	factory, ok := commands[cmd]
	if !ok {
		PrintError(fmt.Errorf("unknown command: %s", cmd))
		os.Exit(1)
	}

	if err := factory().Run(args[1:]); err != nil {
		PrintError(err)
		os.Exit(1)
	}
}

func printHelp(version string) {
	fmt.Printf(`basecamp - Basecamp CLI %s

Usage: basecamp <command> [arguments] [flags]

Commands:
  init                              Configure credentials
  auth                              Authenticate with OAuth
  projects                          List all projects
  boards [project_id]               List card tables in a project
  cards [project_id] <board_id>     List cards (--column <name> to filter)
  card [project_id] <card_id>       View card details (--comments for comments)
  move [project_id] <board> <card>  Move card (--to <column> required)
  version                           Show version

Project ID can be omitted if .basecamp.yml exists in current or parent directory:
  project_id: 12345678

Examples:
  basecamp projects
  basecamp boards 12345678
  basecamp boards                   # uses .basecamp.yml
  basecamp cards 12345678 87654321 --column "In Progress"
  basecamp card 12345678 44444444 --comments
  basecamp move 12345678 87654321 44444444 --to "Done"
`, version)
}

func PrintError(err error) {
	errJSON, _ := json.Marshal(map[string]string{"error": err.Error()})
	fmt.Fprintln(os.Stderr, string(errJSON))
}

func PrintJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// getProjectID returns project ID from args[0] or .basecamp.yml, plus remaining args.
// If project_id comes from config, args are returned unchanged.
// If project_id comes from args[0], remaining args are returned.
func getProjectID(args []string) (projectID string, remaining []string, err error) {
	// First try to get from config
	configProjectID, err := config.FindProjectID()
	if err != nil {
		return "", nil, err
	}

	if configProjectID != "" {
		// Use config, all args are remaining
		return configProjectID, args, nil
	}

	// Need project_id from args
	if len(args) < 1 {
		return "", nil, errors.New("project_id required: provide as argument or create .basecamp.yml with project_id")
	}
	return args[0], args[1:], nil
}
