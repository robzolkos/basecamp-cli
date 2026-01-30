package commands

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

type TodosCmd struct{}

type Todo struct {
	ID          int        `json:"id"`
	Content     string     `json:"content"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	DueOn       string     `json:"due_on"`
	StartsOn    string     `json:"starts_on"`
	Creator     Creator    `json:"creator"`
	Assignees   []Assignee `json:"assignees"`
	CommentsURL string     `json:"comments_url"`
}

type TodoOutput struct {
	ID        int      `json:"id"`
	Content   string   `json:"content"`
	Completed bool     `json:"completed"`
	DueOn     string   `json:"due_on,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
}

type TodosOutput struct {
	TodolistID int          `json:"todolist_id"`
	Todos      []TodoOutput `json:"todos"`
}

func (c *TodosCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("usage: basecamp todos [project_id] <todolist_id> [--completed]")
	}
	todolistID := remaining[0]

	// Check for --completed flag
	showCompleted := false
	for _, arg := range remaining[1:] {
		if arg == "--completed" {
			showCompleted = true
			break
		}
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	// Build URL with optional completed filter
	url := "/buckets/" + projectID + "/todolists/" + todolistID + "/todos.json"
	if showCompleted {
		url += "?completed=true"
	}

	data, err := cl.Get(url)
	if err != nil {
		return err
	}

	var todos []Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		return err
	}

	output := TodosOutput{
		TodolistID: 0, // Will be parsed from todolistID
		Todos:      make([]TodoOutput, len(todos)),
	}

	// Parse todolistID
	if tlID, err := strconv.Atoi(todolistID); err == nil {
		output.TodolistID = tlID
	}

	for i, todo := range todos {
		var assignees []string
		for _, a := range todo.Assignees {
			assignees = append(assignees, a.Name)
		}

		output.Todos[i] = TodoOutput{
			ID:        todo.ID,
			Content:   stripHTML(todo.Content),
			Completed: todo.Completed,
			DueOn:     todo.DueOn,
			Assignees: assignees,
		}
	}

	return PrintJSON(output)
}

type TodoCmd struct{}

type TodoDetailOutput struct {
	ID          int      `json:"id"`
	Content     string   `json:"content"`
	Description string   `json:"description,omitempty"`
	Completed   bool     `json:"completed"`
	DueOn       string   `json:"due_on,omitempty"`
	StartsOn    string   `json:"starts_on,omitempty"`
	Creator     string   `json:"creator"`
	Assignees   []string `json:"assignees,omitempty"`
}

func (c *TodoCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("usage: basecamp todo [project_id] <todo_id>")
	}
	todoID := remaining[0]

	cl, err := client.New()
	if err != nil {
		return err
	}

	data, err := cl.Get("/buckets/" + projectID + "/todos/" + todoID + ".json")
	if err != nil {
		return err
	}

	var todo Todo
	if err := json.Unmarshal(data, &todo); err != nil {
		return err
	}

	var assignees []string
	for _, a := range todo.Assignees {
		assignees = append(assignees, a.Name)
	}

	output := TodoDetailOutput{
		ID:          todo.ID,
		Content:     stripHTML(todo.Content),
		Description: stripHTML(todo.Description),
		Completed:   todo.Completed,
		DueOn:       todo.DueOn,
		StartsOn:    todo.StartsOn,
		Creator:     todo.Creator.Name,
		Assignees:   assignees,
	}

	return PrintJSON(output)
}

type TodoCreateCmd struct{}

func (c *TodoCreateCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("usage: basecamp todo-create [project_id] <todolist_id> --content <text> [--due <date>] [--description <text>]")
	}
	todolistID := remaining[0]

	var content, description, dueOn string
	for i := 1; i < len(remaining); i++ {
		switch remaining[i] {
		case "--content":
			if i+1 < len(remaining) {
				content = remaining[i+1]
				i++
			}
		case "--description":
			if i+1 < len(remaining) {
				description = remaining[i+1]
				i++
			}
		case "--due":
			if i+1 < len(remaining) {
				dueOn = remaining[i+1]
				i++
			}
		}
	}

	if content == "" {
		return errors.New("--content is required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	payload := map[string]any{
		"content": content,
	}
	if description != "" {
		payload["description"] = description
	}
	if dueOn != "" {
		payload["due_on"] = dueOn
	}

	data, err := cl.Post("/buckets/"+projectID+"/todolists/"+todolistID+"/todos.json", payload)
	if err != nil {
		return err
	}

	var todo Todo
	if err := json.Unmarshal(data, &todo); err != nil {
		return err
	}

	return PrintJSON(map[string]any{
		"status":  "ok",
		"id":      todo.ID,
		"content": stripHTML(todo.Content),
		"message": "Todo created",
	})
}

type TodoCompleteCmd struct{}

func (c *TodoCompleteCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("usage: basecamp todo-complete [project_id] <todo_id>")
	}
	todoID := remaining[0]

	cl, err := client.New()
	if err != nil {
		return err
	}

	_, err = cl.Post("/buckets/"+projectID+"/todos/"+todoID+"/completion.json", nil)
	if err != nil {
		// 204 No Content is expected, which may come back as empty response
		if !strings.Contains(err.Error(), "API error") {
			return err
		}
	}

	return PrintJSON(map[string]any{
		"status":  "ok",
		"todo_id": todoID,
		"message": "Todo completed",
	})
}

type TodoUncompleteCmd struct{}

func (c *TodoUncompleteCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("usage: basecamp todo-uncomplete [project_id] <todo_id>")
	}
	todoID := remaining[0]

	cl, err := client.New()
	if err != nil {
		return err
	}

	_, err = cl.Delete("/buckets/" + projectID + "/todos/" + todoID + "/completion.json")
	if err != nil {
		// 204 No Content is expected
		if !strings.Contains(err.Error(), "API error") {
			return err
		}
	}

	return PrintJSON(map[string]any{
		"status":  "ok",
		"todo_id": todoID,
		"message": "Todo uncompleted",
	})
}

// TodoRepositionCmd repositions a todo within its list
type TodoRepositionCmd struct{}

func (c *TodoRepositionCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("todo_id required")
	}
	todoID := remaining[0]

	// Parse --position flag
	var position string
	for i := 1; i < len(remaining); i++ {
		if remaining[i] == "--position" && i+1 < len(remaining) {
			position = remaining[i+1]
			break
		}
	}

	if position == "" {
		return errors.New("--position required (1-indexed)")
	}

	// Validate position is a number
	pos, err := strconv.Atoi(position)
	if err != nil {
		return errors.New("--position must be a number")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	payload := map[string]any{
		"position": pos,
	}

	_, err = cl.Put("/buckets/"+projectID+"/todos/"+todoID+"/position.json", payload)
	if err != nil {
		return err
	}

	return PrintJSON(map[string]any{
		"status":   "ok",
		"todo_id":  todoID,
		"position": position,
		"message":  "Todo repositioned to " + position,
	})
}
