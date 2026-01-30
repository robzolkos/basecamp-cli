package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

var errBoardIDRequired = errors.New("usage: basecamp cards [project_id] <board_id> [--column <name>]")

type CardsCmd struct{}

type ColumnDetail struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Color      string `json:"color"`
	CardsCount int    `json:"cards_count"`
	CardsURL   string `json:"cards_url"`
}

type CardTableDetail struct {
	ID    int            `json:"id"`
	Title string         `json:"title"`
	Lists []ColumnDetail `json:"lists"`
}

type Creator struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CardSummary struct {
	ID      int     `json:"id"`
	Title   string  `json:"title"`
	Creator Creator `json:"creator"`
}

type CardOutput struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Creator string `json:"creator"`
}

type ColumnCards struct {
	Column string       `json:"column"`
	Cards  []CardOutput `json:"cards"`
}

type CardsOutput struct {
	BoardID    int           `json:"board_id"`
	BoardTitle string        `json:"board_title"`
	Columns    []ColumnCards `json:"columns"`
}

func (c *CardsCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errBoardIDRequired
	}
	boardID := remaining[0]

	var columnFilter string
	for i := 1; i < len(remaining); i++ {
		if remaining[i] == "--column" && i+1 < len(remaining) {
			columnFilter = remaining[i+1]
			break
		}
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	// Get the card table
	data, err := cl.Get("/buckets/" + projectID + "/card_tables/" + boardID + ".json")
	if err != nil {
		return err
	}

	var cardTable CardTableDetail
	if err := json.Unmarshal(data, &cardTable); err != nil {
		return err
	}

	output := CardsOutput{
		BoardID:    cardTable.ID,
		BoardTitle: cardTable.Title,
		Columns:    []ColumnCards{},
	}

	for _, list := range cardTable.Lists {
		// Filter by column if specified
		if columnFilter != "" && !strings.Contains(strings.ToLower(list.Title), strings.ToLower(columnFilter)) {
			continue
		}

		if list.CardsCount == 0 {
			continue
		}

		// Fetch cards from this column
		cardsData, err := cl.Get(list.CardsURL)
		if err != nil {
			return err
		}

		var cards []CardSummary
		if err := json.Unmarshal(cardsData, &cards); err != nil {
			return err
		}

		columnCards := ColumnCards{
			Column: list.Title,
			Cards:  make([]CardOutput, len(cards)),
		}

		for i, card := range cards {
			creator := "Unknown"
			if card.Creator.Name != "" {
				creator = card.Creator.Name
			}
			columnCards.Cards[i] = CardOutput{
				ID:      card.ID,
				Title:   card.Title,
				Creator: creator,
			}
		}

		output.Columns = append(output.Columns, columnCards)
	}

	return PrintJSON(output)
}

// ColumnsCmd lists columns in a card table
type ColumnsCmd struct{}

type ColumnOutputBrief struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Color      string `json:"color"`
	CardsCount int    `json:"cards_count"`
}

type ColumnsOutput struct {
	BoardID    int                 `json:"board_id"`
	BoardTitle string              `json:"board_title"`
	Columns    []ColumnOutputBrief `json:"columns"`
}

func (c *ColumnsCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("usage: basecamp columns [project_id] <board_id>")
	}
	boardID := remaining[0]

	cl, err := client.New()
	if err != nil {
		return err
	}

	// Get the card table
	data, err := cl.Get("/buckets/" + projectID + "/card_tables/" + boardID + ".json")
	if err != nil {
		return err
	}

	var cardTable CardTableDetail
	if err := json.Unmarshal(data, &cardTable); err != nil {
		return err
	}

	output := ColumnsOutput{
		BoardID:    cardTable.ID,
		BoardTitle: cardTable.Title,
		Columns:    make([]ColumnOutputBrief, len(cardTable.Lists)),
	}

	for i, col := range cardTable.Lists {
		output.Columns[i] = ColumnOutputBrief{
			ID:         col.ID,
			Title:      col.Title,
			Color:      col.Color,
			CardsCount: col.CardsCount,
		}
	}

	return PrintJSON(output)
}

// CardCreateCmd creates a new card in a column
type CardCreateCmd struct{}

type CardCreateOutput struct {
	Status  string `json:"status"`
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

func (c *CardCreateCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// Parse board_id, column_id, and flags
	var boardID, columnID, title, content, dueOn string

	for i := 0; i < len(remaining); i++ {
		switch remaining[i] {
		case "--column":
			if i+1 < len(remaining) {
				columnID = remaining[i+1]
				i++
			}
		case "--title":
			if i+1 < len(remaining) {
				title = remaining[i+1]
				i++
			}
		case "--content":
			if i+1 < len(remaining) {
				content = remaining[i+1]
				i++
			}
		case "--due":
			if i+1 < len(remaining) {
				dueOn = remaining[i+1]
				i++
			}
		default:
			if boardID == "" {
				boardID = remaining[i]
			}
		}
	}

	if boardID == "" {
		return errors.New("board_id required")
	}
	if columnID == "" {
		return errors.New("--column required (column ID)")
	}
	if title == "" {
		return errors.New("--title required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	// Create card
	payload := map[string]any{
		"title": title,
	}
	if content != "" {
		payload["content"] = content
	}
	if dueOn != "" {
		payload["due_on"] = dueOn
	}

	path := fmt.Sprintf("/buckets/%s/card_tables/lists/%s/cards.json", projectID, columnID)
	responseData, err := cl.Post(path, payload)
	if err != nil {
		return err
	}

	var created CardSummary
	if err := json.Unmarshal(responseData, &created); err != nil {
		return err
	}

	return PrintJSON(CardCreateOutput{
		Status:  "ok",
		ID:      created.ID,
		Title:   created.Title,
		Message: fmt.Sprintf("Card '%s' created", created.Title),
	})
}

// CardUpdateCmd updates an existing card
type CardUpdateCmd struct{}

type CardUpdateOutput struct {
	Status  string `json:"status"`
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

func (c *CardUpdateCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// Parse card_id and flags
	var cardID, title, content, dueOn string

	for i := 0; i < len(remaining); i++ {
		switch remaining[i] {
		case "--title":
			if i+1 < len(remaining) {
				title = remaining[i+1]
				i++
			}
		case "--content":
			if i+1 < len(remaining) {
				content = remaining[i+1]
				i++
			}
		case "--due":
			if i+1 < len(remaining) {
				dueOn = remaining[i+1]
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
	if title == "" && content == "" && dueOn == "" {
		return errors.New("at least one of --title, --content, or --due required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	// Update card
	payload := map[string]any{}
	if title != "" {
		payload["title"] = title
	}
	if content != "" {
		payload["content"] = content
	}
	if dueOn != "" {
		payload["due_on"] = dueOn
	}

	path := fmt.Sprintf("/buckets/%s/card_tables/cards/%s.json", projectID, cardID)
	responseData, err := cl.Put(path, payload)
	if err != nil {
		return err
	}

	var updated CardSummary
	if err := json.Unmarshal(responseData, &updated); err != nil {
		return err
	}

	return PrintJSON(CardUpdateOutput{
		Status:  "ok",
		ID:      updated.ID,
		Title:   updated.Title,
		Message: fmt.Sprintf("Card '%s' updated", updated.Title),
	})
}
