package commands

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

// Step represents a card step/checklist item
type Step struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Completed bool       `json:"completed"`
	DueOn     string     `json:"due_on"`
	Position  int        `json:"position"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
	Creator   Creator    `json:"creator"`
	Assignees []Assignee `json:"assignees"`
}

// StepOutput is the brief output for steps in card listings
type StepOutput struct {
	ID        int      `json:"id"`
	Title     string   `json:"title"`
	Completed bool     `json:"completed"`
	DueOn     string   `json:"due_on,omitempty"`
	Position  int      `json:"position"`
	Assignees []string `json:"assignees,omitempty"`
}

// StepCreateCmd creates a step in a card
type StepCreateCmd struct{}

type StepCreateOutput struct {
	Status  string `json:"status"`
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

func (c *StepCreateCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// Parse card_id and flags
	var cardID, title, dueOn, assignees string

	for i := 0; i < len(remaining); i++ {
		switch remaining[i] {
		case "--title":
			if i+1 < len(remaining) {
				title = remaining[i+1]
				i++
			}
		case "--due":
			if i+1 < len(remaining) {
				dueOn = remaining[i+1]
				i++
			}
		case "--assignees":
			if i+1 < len(remaining) {
				assignees = remaining[i+1]
				i++
			}
		default:
			if cardID == "" {
				cardID = remaining[i]
			}
		}
	}

	if cardID == "" {
		return errors.New("card_id required")
	}
	if title == "" {
		return errors.New("--title required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	payload := map[string]any{
		"title": title,
	}
	if dueOn != "" {
		payload["due_on"] = dueOn
	}
	if assignees != "" {
		payload["assignees"] = assignees
	}

	path := fmt.Sprintf("/buckets/%s/card_tables/cards/%s/steps.json", projectID, cardID)
	responseData, err := cl.Post(path, payload)
	if err != nil {
		return err
	}

	var created Step
	if err := json.Unmarshal(responseData, &created); err != nil {
		return err
	}

	return PrintJSON(StepCreateOutput{
		Status:  "ok",
		ID:      created.ID,
		Title:   created.Title,
		Message: fmt.Sprintf("Step '%s' created", created.Title),
	})
}

// StepUpdateCmd updates a step
type StepUpdateCmd struct{}

type StepUpdateOutput struct {
	Status  string `json:"status"`
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

func (c *StepUpdateCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// Parse step_id and flags
	var stepID, title, dueOn, assignees string

	for i := 0; i < len(remaining); i++ {
		switch remaining[i] {
		case "--title":
			if i+1 < len(remaining) {
				title = remaining[i+1]
				i++
			}
		case "--due":
			if i+1 < len(remaining) {
				dueOn = remaining[i+1]
				i++
			}
		case "--assignees":
			if i+1 < len(remaining) {
				assignees = remaining[i+1]
				i++
			}
		default:
			if stepID == "" {
				stepID = remaining[i]
			}
		}
	}

	if stepID == "" {
		return errors.New("step_id required")
	}
	if title == "" && dueOn == "" && assignees == "" {
		return errors.New("at least one of --title, --due, or --assignees required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	payload := map[string]any{}
	if title != "" {
		payload["title"] = title
	}
	if dueOn != "" {
		payload["due_on"] = dueOn
	}
	if assignees != "" {
		payload["assignees"] = assignees
	}

	path := fmt.Sprintf("/buckets/%s/card_tables/steps/%s.json", projectID, stepID)
	responseData, err := cl.Put(path, payload)
	if err != nil {
		return err
	}

	var updated Step
	if err := json.Unmarshal(responseData, &updated); err != nil {
		return err
	}

	return PrintJSON(StepUpdateOutput{
		Status:  "ok",
		ID:      updated.ID,
		Title:   updated.Title,
		Message: fmt.Sprintf("Step '%s' updated", updated.Title),
	})
}

// StepCompleteCmd marks a step as complete
type StepCompleteCmd struct{}

type StepCompleteOutput struct {
	Status  string `json:"status"`
	StepID  string `json:"step_id"`
	Message string `json:"message"`
}

func (c *StepCompleteCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("usage: basecamp step-complete [project_id] <step_id>")
	}
	stepID := remaining[0]

	cl, err := client.New()
	if err != nil {
		return err
	}

	payload := map[string]string{
		"completion": "on",
	}

	path := fmt.Sprintf("/buckets/%s/card_tables/steps/%s/completions.json", projectID, stepID)
	_, err = cl.Put(path, payload)
	if err != nil {
		return err
	}

	return PrintJSON(StepCompleteOutput{
		Status:  "ok",
		StepID:  stepID,
		Message: "Step completed",
	})
}

// StepUncompleteCmd marks a step as uncomplete
type StepUncompleteCmd struct{}

func (c *StepUncompleteCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("usage: basecamp step-uncomplete [project_id] <step_id>")
	}
	stepID := remaining[0]

	cl, err := client.New()
	if err != nil {
		return err
	}

	payload := map[string]string{
		"completion": "off",
	}

	path := fmt.Sprintf("/buckets/%s/card_tables/steps/%s/completions.json", projectID, stepID)
	_, err = cl.Put(path, payload)
	if err != nil {
		return err
	}

	return PrintJSON(StepCompleteOutput{
		Status:  "ok",
		StepID:  stepID,
		Message: "Step uncompleted",
	})
}

// StepRepositionCmd changes a step's position
type StepRepositionCmd struct{}

type StepRepositionOutput struct {
	Status   string `json:"status"`
	StepID   string `json:"step_id"`
	Position string `json:"position"`
	Message  string `json:"message"`
}

func (c *StepRepositionCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// Parse card_id, step_id, and --position flag
	var cardID, stepID, position string

	for i := 0; i < len(remaining); i++ {
		switch remaining[i] {
		case "--position":
			if i+1 < len(remaining) {
				position = remaining[i+1]
				i++
			}
		default:
			if cardID == "" {
				cardID = remaining[i]
			} else if stepID == "" {
				stepID = remaining[i]
			}
		}
	}

	if cardID == "" {
		return errors.New("card_id required")
	}
	if stepID == "" {
		return errors.New("step_id required")
	}
	if position == "" {
		return errors.New("--position required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	payload := map[string]any{
		"source_id": stepID,
		"position":  position,
	}

	path := fmt.Sprintf("/buckets/%s/card_tables/cards/%s/positions.json", projectID, cardID)
	_, err = cl.Post(path, payload)
	if err != nil {
		return err
	}

	return PrintJSON(StepRepositionOutput{
		Status:   "ok",
		StepID:   stepID,
		Position: position,
		Message:  fmt.Sprintf("Step repositioned to %s", position),
	})
}
