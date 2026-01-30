package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

// fetchMessageBoard gets the message board for a project
func fetchMessageBoard(cl *client.Client, projectID string) (ProjectDetail, MessageBoard, error) {
	project, err := fetchProject(cl, projectID)
	if err != nil {
		return ProjectDetail{}, MessageBoard{}, err
	}

	messageBoardURL, err := getDockURL(project, "message_board")
	if err != nil {
		return ProjectDetail{}, MessageBoard{}, err
	}

	boardData, err := cl.Get(messageBoardURL)
	if err != nil {
		return ProjectDetail{}, MessageBoard{}, err
	}

	var board MessageBoard
	if err := json.Unmarshal(boardData, &board); err != nil {
		return ProjectDetail{}, MessageBoard{}, err
	}

	return project, board, nil
}

type MessagesCmd struct{}

type MessageBoard struct {
	ID          int    `json:"id"`
	MessagesURL string `json:"messages_url"`
}

type Message struct {
	ID            int     `json:"id"`
	Subject       string  `json:"subject"`
	Content       string  `json:"content"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	CommentsCount int     `json:"comments_count"`
	CommentsURL   string  `json:"comments_url"`
	URL           string  `json:"app_url"`
	Creator       Creator `json:"creator"`
}

type MessageListOutput struct {
	ProjectID      int                  `json:"project_id"`
	MessageBoardID int                  `json:"message_board_id"`
	Messages       []MessageOutputBrief `json:"messages"`
}

type MessageOutputBrief struct {
	ID            int    `json:"id"`
	Subject       string `json:"subject"`
	Creator       string `json:"creator"`
	CreatedAt     string `json:"created_at"`
	CommentsCount int    `json:"comments_count"`
}

func (c *MessagesCmd) Run(args []string) error {
	projectID, _, err := getProjectID(args)
	if err != nil {
		return err
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	project, board, err := fetchMessageBoard(cl, projectID)
	if err != nil {
		return err
	}

	// Get messages
	messagesData, err := cl.Get(board.MessagesURL)
	if err != nil {
		return err
	}

	var messages []Message
	if err := json.Unmarshal(messagesData, &messages); err != nil {
		return err
	}

	output := MessageListOutput{
		ProjectID:      project.ID,
		MessageBoardID: board.ID,
		Messages:       make([]MessageOutputBrief, len(messages)),
	}

	for i, m := range messages {
		output.Messages[i] = MessageOutputBrief{
			ID:            m.ID,
			Subject:       m.Subject,
			Creator:       m.Creator.Name,
			CreatedAt:     m.CreatedAt,
			CommentsCount: m.CommentsCount,
		}
	}

	return PrintJSON(output)
}

type MessageCmd struct{}

type MessageDetailOutput struct {
	ID            int             `json:"id"`
	Subject       string          `json:"subject"`
	Content       string          `json:"content"`
	Creator       string          `json:"creator"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
	CommentsCount int             `json:"comments_count"`
	URL           string          `json:"url"`
	Comments      []CommentOutput `json:"comments,omitempty"`
}

func (c *MessageCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// Parse message_id and flags
	var messageID string
	showComments := false

	for i := 0; i < len(remaining); i++ {
		if remaining[i] == "--comments" {
			showComments = true
		} else if messageID == "" {
			messageID = remaining[i]
		}
	}

	if messageID == "" {
		return errors.New("message_id required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	// Get message directly
	messageData, err := cl.Get("/buckets/" + projectID + "/messages/" + messageID + ".json")
	if err != nil {
		return err
	}

	var message Message
	if err := json.Unmarshal(messageData, &message); err != nil {
		return err
	}

	output := MessageDetailOutput{
		ID:            message.ID,
		Subject:       message.Subject,
		Content:       stripHTML(message.Content),
		Creator:       message.Creator.Name,
		CreatedAt:     message.CreatedAt,
		UpdatedAt:     message.UpdatedAt,
		CommentsCount: message.CommentsCount,
		URL:           message.URL,
	}

	if showComments && message.CommentsURL != "" {
		comments, err := fetchComments(cl, message.CommentsURL)
		if err != nil {
			return err
		}
		output.Comments = comments
	}

	return PrintJSON(output)
}

type MessageCreateCmd struct{}

type MessageCreateOutput struct {
	Status  string `json:"status"`
	ID      int    `json:"id"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (c *MessageCreateCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// Parse flags
	var subject, content string

	for i := 0; i < len(remaining); i++ {
		switch remaining[i] {
		case "--subject":
			if i+1 < len(remaining) {
				subject = remaining[i+1]
				i++
			}
		case "--content":
			if i+1 < len(remaining) {
				content = remaining[i+1]
				i++
			}
		}
	}

	if subject == "" {
		return errors.New("--subject required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	_, board, err := fetchMessageBoard(cl, projectID)
	if err != nil {
		return err
	}

	// Create message
	payload := map[string]any{
		"subject": subject,
		"status":  "active",
	}
	if content != "" {
		payload["content"] = content
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// POST to messages URL
	messagesURL := board.MessagesURL
	// Convert from full URL to path
	if idx := strings.Index(messagesURL, "/buckets/"); idx != -1 {
		messagesURL = messagesURL[idx:]
	}

	responseData, err := cl.Post(messagesURL, payloadBytes)
	if err != nil {
		return err
	}

	var created Message
	if err := json.Unmarshal(responseData, &created); err != nil {
		return err
	}

	return PrintJSON(MessageCreateOutput{
		Status:  "ok",
		ID:      created.ID,
		Subject: created.Subject,
		Message: fmt.Sprintf("Message '%s' created", created.Subject),
	})
}
