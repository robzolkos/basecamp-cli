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
	"init":                  func() Command { return &InitCmd{} },
	"auth":                  func() Command { return &AuthCmd{} },
	"projects":              func() Command { return &ProjectsCmd{} },
	"boards":                func() Command { return &BoardsCmd{} },
	"cards":                 func() Command { return &CardsCmd{} },
	"card":                  func() Command { return &CardCmd{} },
	"move":                  func() Command { return &MoveCmd{} },
	"todolists":             func() Command { return &TodolistsCmd{} },
	"todos":                 func() Command { return &TodosCmd{} },
	"todo":                  func() Command { return &TodoCmd{} },
	"todo-create":           func() Command { return &TodoCreateCmd{} },
	"todo-complete":         func() Command { return &TodoCompleteCmd{} },
	"todo-uncomplete":       func() Command { return &TodoUncompleteCmd{} },
	"messages":              func() Command { return &MessagesCmd{} },
	"message":               func() Command { return &MessageCmd{} },
	"message-create":        func() Command { return &MessageCreateCmd{} },
	"comment-add":           func() Command { return &CommentAddCmd{} },
	"docs":                  func() Command { return &DocsCmd{} },
	"doc":                   func() Command { return &DocCmd{} },
	"doc-create":            func() Command { return &DocCreateCmd{} },
	"schedule":              func() Command { return &ScheduleCmd{} },
	"event":                 func() Command { return &EventCmd{} },
	"event-create":          func() Command { return &EventCreateCmd{} },
	"campfire":              func() Command { return &CampfireCmd{} },
	"campfire-post":         func() Command { return &CampfirePostCmd{} },
	"columns":               func() Command { return &ColumnsCmd{} },
	"card-create":           func() Command { return &CardCreateCmd{} },
	"card-update":           func() Command { return &CardUpdateCmd{} },
	"search":                func() Command { return &SearchCmd{} },
	"step-create":           func() Command { return &StepCreateCmd{} },
	"step-update":           func() Command { return &StepUpdateCmd{} },
	"step-complete":         func() Command { return &StepCompleteCmd{} },
	"step-uncomplete":       func() Command { return &StepUncompleteCmd{} },
	"step-reposition":       func() Command { return &StepRepositionCmd{} },
	"people":                func() Command { return &PeopleCmd{} },
	"person":                func() Command { return &PersonCmd{} },
	"people-pingable":       func() Command { return &PeoplePingableCmd{} },
	"people-project":        func() Command { return &PeopleProjectCmd{} },
	"my-profile":            func() Command { return &MyProfileCmd{} },
	"project-access":        func() Command { return &ProjectAccessCmd{} },
	"questionnaire":         func() Command { return &QuestionnaireCmd{} },
	"questions":             func() Command { return &QuestionsCmd{} },
	"question":              func() Command { return &QuestionCmd{} },
	"question-answers":      func() Command { return &QuestionAnswersCmd{} },
	"question-answer":       func() Command { return &QuestionAnswerCmd{} },
	"todolist-groups":       func() Command { return &TodolistGroupsCmd{} },
	"todolist-group":        func() Command { return &TodolistGroupCmd{} },
	"todolist-group-create": func() Command { return &TodolistGroupCreateCmd{} },
	"todo-reposition":       func() Command { return &TodoRepositionCmd{} },
	"upload":                func() Command { return &UploadCmd{} },
	"uploads":               func() Command { return &UploadsCmd{} },
	"upload-view":           func() Command { return &UploadViewCmd{} },
	"archive":               func() Command { return &ArchiveCmd{} },
	"unarchive":             func() Command { return &UnarchiveCmd{} },
	"trash":                 func() Command { return &TrashCmd{} },
	"message-types":         func() Command { return &MessageTypesCmd{} },
	"message-type":          func() Command { return &MessageTypeCmd{} },
	"message-type-create":   func() Command { return &MessageTypeCreateCmd{} },
	"message-type-update":   func() Command { return &MessageTypeUpdateCmd{} },
	"message-type-delete":   func() Command { return &MessageTypeDeleteCmd{} },
	"events":                func() Command { return &EventsCmd{} },
	"events-project":        func() Command { return &EventsProjectCmd{} },
	"events-recording":      func() Command { return &EventsRecordingCmd{} },
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

Card Tables:
  boards [project_id]               List card tables in a project
  columns [project_id] <board_id>   List columns in a board
  cards [project_id] <board_id>     List cards (--column <name> to filter)
  card [project_id] <card_id>       View card details (--comments for comments)
  card-create [project_id] <board>  Create card (--column, --title required)
  card-update [project_id] <card>   Update card (--title, --content, --due)
  move [project_id] <board> <card>  Move card (--to <column> required)

Card Steps:
  step-create [project_id] <card>   Create step (--title required)
  step-update [project_id] <step>   Update step (--title, --due, --assignees)
  step-complete [project_id] <step> Mark step as complete
  step-uncomplete [project_id] <step> Mark step as incomplete
  step-reposition [project_id] <card> <step> Reposition step (--position required)

Todos:
  todolists [project_id]            List todo lists
  todos [project_id] <todolist_id>  List todos (--completed for completed)
  todo [project_id] <todo_id>       View todo details
  todo-create [project_id] <list>   Create todo (--content required)
  todo-complete [project_id] <id>   Mark todo as complete
  todo-uncomplete [project_id] <id> Mark todo as incomplete
  todo-reposition [project_id] <id> Reposition todo (--position required)

Todo Groups:
  todolist-groups [project_id] <list> List groups in a todolist
  todolist-group [project_id] <id>    View group details
  todolist-group-create [project_id] <list> Create group (--name required)

Messages:
  messages [project_id]             List messages
  message [project_id] <message_id> View message (--comments for comments)
  message-create [project_id]       Create message (--subject required)

Comments:
  comment-add [project_id] <id>     Add comment to recording (--content required)

Documents:
  docs [project_id]                 List documents
  doc [project_id] <doc_id>         View document (--comments for comments)
  doc-create [project_id]           Create document (--title required)

Schedule:
  schedule [project_id]             List schedule entries
  event [project_id] <entry_id>     View event (--comments for comments)
  event-create [project_id]         Create event (--summary, --starts-at, --ends-at)

Campfire:
  campfire [project_id]             List campfire messages
  campfire-post [project_id]        Post to campfire (--content required)

Search:
  search <query>                    Search across all projects
                                    (--type <type>, --project <id> optional)

People:
  people                            List all people
  person <person_id>                View person details
  people-pingable                   List pingable people
  people-project [project_id]       List people on a project
  my-profile                        View your own profile
  project-access [project_id]       Update project access (--grant, --revoke)

Automatic Check-ins:
  questionnaire [project_id]        Get questionnaire info
  questions [project_id]            List check-in questions
  question [project_id] <id>        View question (--comments for comments)
  question-answers [project_id] <q> List answers to a question
  question-answer [project_id] <id> View answer (--comments for comments)

Uploads:
  upload <file>                     Upload a file (returns attachable_sgid)
  uploads [project_id] <vault_id>   List uploads in a vault
  upload-view [project_id] <id>     View upload (--comments for comments)

Recordings:
  archive [project_id] <id>         Archive a recording
  unarchive [project_id] <id>       Unarchive a recording
  trash [project_id] <id>           Trash a recording

Message Types:
  message-types [project_id]        List message types
  message-type [project_id] <id>    View message type
  message-type-create [project_id]  Create type (--name, --icon required)
  message-type-update [project_id] <id> Update type (--name, --icon)
  message-type-delete [project_id] <id> Delete type

Activity Events:
  events                            List all events (across all projects)
  events-project [project_id]       List events for a project
  events-recording [project_id] <id> List events for a recording

  version                           Show version

Project ID can be omitted if .basecamp.yml exists in current or parent directory:
  project_id: 12345678

Examples:
  basecamp projects
  basecamp boards
  basecamp todolists
  basecamp todos 12345
  basecamp todo-create 12345 --content "Fix bug"
  basecamp todo-complete 67890
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
