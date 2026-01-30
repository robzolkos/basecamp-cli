package commands

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

// fetchTodoSet gets the todoset for a project
func fetchTodoSet(cl *client.Client, projectID string) (ProjectDetail, TodoSet, error) {
	project, err := fetchProject(cl, projectID)
	if err != nil {
		return ProjectDetail{}, TodoSet{}, err
	}

	todosetURL, err := getDockURL(project, "todoset")
	if err != nil {
		return ProjectDetail{}, TodoSet{}, err
	}

	data, err := cl.Get(todosetURL)
	if err != nil {
		return ProjectDetail{}, TodoSet{}, err
	}

	var todoset TodoSet
	if err := json.Unmarshal(data, &todoset); err != nil {
		return ProjectDetail{}, TodoSet{}, err
	}

	return project, todoset, nil
}

type TodolistsCmd struct{}

type TodoSet struct {
	ID           int    `json:"id"`
	TodolistsURL string `json:"todolists_url"`
	GroupsURL    string `json:"groups_url"`
}

type Todolist struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	TodosURL       string `json:"todos_url"`
	CompletedRatio string `json:"completed_ratio"`
	Completed      bool   `json:"completed"`
}

type TodolistOutput struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	CompletedRatio string `json:"completed_ratio"`
	Completed      bool   `json:"completed"`
}

type TodolistsOutput struct {
	ProjectID int              `json:"project_id"`
	TodosetID int              `json:"todoset_id"`
	Todolists []TodolistOutput `json:"todolists"`
}

func (c *TodolistsCmd) Run(args []string) error {
	projectID, _, err := getProjectID(args)
	if err != nil {
		return err
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	project, todoset, err := fetchTodoSet(cl, projectID)
	if err != nil {
		return err
	}

	// Get todolists
	todolistsData, err := cl.Get(todoset.TodolistsURL)
	if err != nil {
		return err
	}

	var todolists []Todolist
	if err := json.Unmarshal(todolistsData, &todolists); err != nil {
		return err
	}

	output := TodolistsOutput{
		ProjectID: project.ID,
		TodosetID: todoset.ID,
		Todolists: make([]TodolistOutput, len(todolists)),
	}

	for i, tl := range todolists {
		output.Todolists[i] = TodolistOutput{
			ID:             tl.ID,
			Title:          tl.Title,
			CompletedRatio: tl.CompletedRatio,
			Completed:      tl.Completed,
		}
	}

	return PrintJSON(output)
}

// TodolistGroup represents a group of todolists
type TodolistGroup struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Color     string `json:"color"`
}

// TodolistGroupsCmd lists todolist groups within a todolist
type TodolistGroupsCmd struct{}

type TodolistGroupsOutput struct {
	ProjectID  int                  `json:"project_id"`
	TodolistID int                  `json:"todolist_id"`
	Groups     []TodolistGroupBrief `json:"groups"`
}

type TodolistGroupBrief struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Position int    `json:"position"`
	Color    string `json:"color,omitempty"`
}

func (c *TodolistGroupsCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("todolist_id required")
	}
	todolistID := remaining[0]

	cl, err := client.New()
	if err != nil {
		return err
	}

	// Groups are listed under a specific todolist
	groupsURL := fmt.Sprintf("/buckets/%s/todolists/%s/groups.json", projectID, todolistID)

	data, err := cl.Get(groupsURL)
	if err != nil {
		return err
	}

	var groups []TodolistGroup
	if err := json.Unmarshal(data, &groups); err != nil {
		return err
	}

	var tlID int
	fmt.Sscanf(todolistID, "%d", &tlID)

	output := TodolistGroupsOutput{
		ProjectID:  0,
		TodolistID: tlID,
		Groups:     make([]TodolistGroupBrief, len(groups)),
	}
	fmt.Sscanf(projectID, "%d", &output.ProjectID)

	for i, g := range groups {
		output.Groups[i] = TodolistGroupBrief{
			ID:       g.ID,
			Name:     g.Name,
			Position: g.Position,
			Color:    g.Color,
		}
	}

	return PrintJSON(output)
}

// TodolistGroupCmd views a single group (groups are todolists with group_position_url)
type TodolistGroupCmd struct{}

type TodolistGroupDetailOutput struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Position  int    `json:"position"`
	Color     string `json:"color,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (c *TodolistGroupCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("group_id required")
	}
	groupID := remaining[0]

	cl, err := client.New()
	if err != nil {
		return err
	}

	// Groups are accessed via the todolists endpoint (they are todolists with group_position_url)
	data, err := cl.Get("/buckets/" + projectID + "/todolists/" + groupID + ".json")
	if err != nil {
		return err
	}

	var group TodolistGroup
	if err := json.Unmarshal(data, &group); err != nil {
		return err
	}

	return PrintJSON(TodolistGroupDetailOutput{
		ID:        group.ID,
		Name:      group.Name,
		Position:  group.Position,
		Color:     group.Color,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
	})
}

// TodolistGroupCreateCmd creates a new group within a todolist
type TodolistGroupCreateCmd struct{}

type TodolistGroupCreateOutput struct {
	Status  string `json:"status"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (c *TodolistGroupCreateCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// First arg should be todolist_id
	if len(remaining) < 1 {
		return errors.New("todolist_id required")
	}
	todolistID := remaining[0]

	// Parse flags
	var name, color string

	for i := 1; i < len(remaining); i++ {
		switch remaining[i] {
		case "--name":
			if i+1 < len(remaining) {
				name = remaining[i+1]
				i++
			}
		case "--color":
			if i+1 < len(remaining) {
				color = remaining[i+1]
				i++
			}
		}
	}

	if name == "" {
		return errors.New("--name required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	payload := map[string]any{
		"name": name,
	}
	if color != "" {
		payload["color"] = color
	}

	// Groups are created under a specific todolist
	data, err := cl.Post("/buckets/"+projectID+"/todolists/"+todolistID+"/groups.json", payload)
	if err != nil {
		return err
	}

	var group TodolistGroup
	if err := json.Unmarshal(data, &group); err != nil {
		return err
	}

	return PrintJSON(TodolistGroupCreateOutput{
		Status:  "ok",
		ID:      group.ID,
		Name:    group.Name,
		Message: fmt.Sprintf("Group '%s' created", group.Name),
	})
}
