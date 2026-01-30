package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

type MoveCmd struct{}

type ColumnInfo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type CardTableForMove struct {
	ID    int          `json:"id"`
	Lists []ColumnInfo `json:"lists"`
}

func (c *MoveCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 2 {
		return errors.New("usage: basecamp move [project_id] <board_id> <card_id> --to <column>")
	}
	boardID := remaining[0]
	cardID := remaining[1]

	var targetColumn string
	for i := 2; i < len(remaining); i++ {
		if remaining[i] == "--to" && i+1 < len(remaining) {
			targetColumn = remaining[i+1]
			break
		}
	}

	if targetColumn == "" {
		return errors.New("--to <column> flag is required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	// Get the card table to find the target column
	data, err := cl.Get("/buckets/" + projectID + "/card_tables/" + boardID + ".json")
	if err != nil {
		return err
	}

	var cardTable CardTableForMove
	if err := json.Unmarshal(data, &cardTable); err != nil {
		return err
	}

	// Find target column
	var column *ColumnInfo
	var columnNames []string
	for _, col := range cardTable.Lists {
		columnNames = append(columnNames, col.Title)
		if strings.EqualFold(col.Title, targetColumn) {
			column = &col
			break
		}
	}

	if column == nil {
		return fmt.Errorf("column '%s' not found. Available columns: %s", targetColumn, strings.Join(columnNames, ", "))
	}

	// Move the card
	_, err = cl.Post("/buckets/"+projectID+"/card_tables/cards/"+cardID+"/moves.json", map[string]int{
		"column_id": column.ID,
	})
	if err != nil {
		return err
	}

	return PrintJSON(map[string]any{
		"status":  "ok",
		"card_id": cardID,
		"column":  column.Title,
		"message": fmt.Sprintf("Card %s moved to '%s'", cardID, column.Title),
	})
}
