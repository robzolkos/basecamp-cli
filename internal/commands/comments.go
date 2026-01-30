package commands

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

type CommentAddCmd struct{}

type CommentAddOutput struct {
	Status      string `json:"status"`
	ID          int    `json:"id"`
	RecordingID string `json:"recording_id"`
	Message     string `json:"message"`
}

func (c *CommentAddCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// Parse recording_id and --content flag
	var recordingID, content string

	for i := 0; i < len(remaining); i++ {
		switch remaining[i] {
		case "--content":
			if i+1 < len(remaining) {
				content = remaining[i+1]
				i++
			}
		default:
			if recordingID == "" {
				recordingID = remaining[i]
			}
		}
	}

	if recordingID == "" {
		return errors.New("recording_id required")
	}
	if content == "" {
		return errors.New("--content required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	// Create comment
	payload := map[string]string{
		"content": content,
	}

	path := fmt.Sprintf("/buckets/%s/recordings/%s/comments.json", projectID, recordingID)
	responseData, err := cl.Post(path, payload)
	if err != nil {
		return err
	}

	var created Comment
	if err := json.Unmarshal(responseData, &created); err != nil {
		return err
	}

	return PrintJSON(CommentAddOutput{
		Status:      "ok",
		ID:          created.ID,
		RecordingID: recordingID,
		Message:     fmt.Sprintf("Comment added to recording %s", recordingID),
	})
}
