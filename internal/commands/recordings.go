package commands

import (
	"errors"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

// ArchiveCmd archives a recording
type ArchiveCmd struct{}

type RecordingStatusOutput struct {
	Status      string `json:"status"`
	RecordingID string `json:"recording_id"`
	Message     string `json:"message"`
}

func (c *ArchiveCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("recording_id required")
	}
	recordingID := remaining[0]

	cl, err := client.New()
	if err != nil {
		return err
	}

	_, err = cl.Put("/buckets/"+projectID+"/recordings/"+recordingID+"/status/archived.json", nil)
	if err != nil {
		return err
	}

	return PrintJSON(RecordingStatusOutput{
		Status:      "ok",
		RecordingID: recordingID,
		Message:     "Recording archived",
	})
}

// UnarchiveCmd unarchives a recording (sets to active)
type UnarchiveCmd struct{}

func (c *UnarchiveCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("recording_id required")
	}
	recordingID := remaining[0]

	cl, err := client.New()
	if err != nil {
		return err
	}

	_, err = cl.Put("/buckets/"+projectID+"/recordings/"+recordingID+"/status/active.json", nil)
	if err != nil {
		return err
	}

	return PrintJSON(RecordingStatusOutput{
		Status:      "ok",
		RecordingID: recordingID,
		Message:     "Recording unarchived",
	})
}

// TrashCmd trashes a recording
type TrashCmd struct{}

func (c *TrashCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("recording_id required")
	}
	recordingID := remaining[0]

	cl, err := client.New()
	if err != nil {
		return err
	}

	_, err = cl.Put("/buckets/"+projectID+"/recordings/"+recordingID+"/status/trashed.json", nil)
	if err != nil {
		return err
	}

	return PrintJSON(RecordingStatusOutput{
		Status:      "ok",
		RecordingID: recordingID,
		Message:     "Recording trashed",
	})
}
