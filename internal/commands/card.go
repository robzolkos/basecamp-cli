package commands

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

type CardCmd struct{}

type Assignee struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CardDetail struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	Description   string     `json:"description"`
	CreatedAt     string     `json:"created_at"`
	UpdatedAt     string     `json:"updated_at"`
	AppURL        string     `json:"app_url"`
	CommentsCount int        `json:"comments_count"`
	CommentsURL   string     `json:"comments_url"`
	Creator       Creator    `json:"creator"`
	Assignees     []Assignee `json:"assignees"`
	Steps         []Step     `json:"steps"`
}

type Comment struct {
	ID        int     `json:"id"`
	Content   string  `json:"content"`
	CreatedAt string  `json:"created_at"`
	Creator   Creator `json:"creator"`
}

type CommentOutput struct {
	ID        int    `json:"id"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type CardDetailOutput struct {
	ID          int             `json:"id"`
	Title       string          `json:"title"`
	Creator     string          `json:"creator"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
	URL         string          `json:"url"`
	Assignees   []string        `json:"assignees,omitempty"`
	Description string          `json:"description"`
	Steps       []StepOutput    `json:"steps,omitempty"`
	Comments    []CommentOutput `json:"comments,omitempty"`
}

func (c *CardCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("usage: basecamp card [project_id] <card_id> [--comments]")
	}
	cardID := remaining[0]

	showComments := false
	for _, arg := range remaining[1:] {
		if arg == "--comments" {
			showComments = true
			break
		}
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	// Get the card
	data, err := cl.Get("/buckets/" + projectID + "/card_tables/cards/" + cardID + ".json")
	if err != nil {
		return err
	}

	var card CardDetail
	if err := json.Unmarshal(data, &card); err != nil {
		return err
	}

	output := CardDetailOutput{
		ID:          card.ID,
		Title:       card.Title,
		Creator:     card.Creator.Name,
		CreatedAt:   card.CreatedAt,
		UpdatedAt:   card.UpdatedAt,
		URL:         card.AppURL,
		Description: stripHTML(coalesce(card.Content, card.Description, "No description")),
	}

	if len(card.Assignees) > 0 {
		output.Assignees = make([]string, len(card.Assignees))
		for i, a := range card.Assignees {
			output.Assignees[i] = a.Name
		}
	}

	if len(card.Steps) > 0 {
		output.Steps = make([]StepOutput, len(card.Steps))
		for i, s := range card.Steps {
			var assignees []string
			for _, a := range s.Assignees {
				assignees = append(assignees, a.Name)
			}
			output.Steps[i] = StepOutput{
				ID:        s.ID,
				Title:     s.Title,
				Completed: s.Completed,
				DueOn:     s.DueOn,
				Position:  s.Position,
				Assignees: assignees,
			}
		}
	}

	if showComments && card.CommentsCount > 0 {
		comments, err := fetchComments(cl, card.CommentsURL)
		if err != nil {
			return err
		}
		output.Comments = comments
	}

	return PrintJSON(output)
}

var htmlTagRegex = regexp.MustCompile(`<[^>]*>`)
var whitespaceRegex = regexp.MustCompile(`\s+`)

func stripHTML(html string) string {
	if html == "" {
		return ""
	}
	text := htmlTagRegex.ReplaceAllString(html, " ")
	text = whitespaceRegex.ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}

func coalesce(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
