package commands

import (
	"encoding/json"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

type BoardsCmd struct{}

type DockItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ProjectDetail struct {
	ID   int        `json:"id"`
	Name string     `json:"name"`
	Dock []DockItem `json:"dock"`
}

type ColumnSummary struct {
	Title      string `json:"title"`
	CardsCount int    `json:"cards_count"`
}

type CardTable struct {
	ID    int             `json:"id"`
	Title string          `json:"title"`
	Lists []ColumnSummary `json:"lists"`
}

type BoardOutput struct {
	ProjectID   int             `json:"project_id"`
	ProjectName string          `json:"project_name"`
	BoardID     int             `json:"board_id"`
	BoardTitle  string          `json:"board_title"`
	Columns     []ColumnSummary `json:"columns"`
}

func (c *BoardsCmd) Run(args []string) error {
	projectID, _, err := getProjectID(args)
	if err != nil {
		return err
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	project, err := fetchProject(cl, projectID)
	if err != nil {
		return err
	}

	cardTableURL, err := getDockURL(project, "kanban_board")
	if err != nil {
		return PrintJSON(map[string]any{
			"project_id":   project.ID,
			"project_name": project.Name,
			"error":        "No card table found in this project",
		})
	}

	// Get the card table
	cardTableData, err := cl.Get(cardTableURL)
	if err != nil {
		return err
	}

	var cardTable CardTable
	if err := json.Unmarshal(cardTableData, &cardTable); err != nil {
		return err
	}

	return PrintJSON(BoardOutput{
		ProjectID:   project.ID,
		ProjectName: project.Name,
		BoardID:     cardTable.ID,
		BoardTitle:  cardTable.Title,
		Columns:     cardTable.Lists,
	})
}
