package commands

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

// Event represents a Basecamp activity event
type Event struct {
	ID            int     `json:"id"`
	Action        string  `json:"action"`
	CreatedAt     string  `json:"created_at"`
	RecordingType string  `json:"recording_type"`
	Creator       Creator `json:"creator"`
	Recording     struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		Type  string `json:"type"`
	} `json:"recording"`
	Bucket struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"bucket"`
}

// EventsCmd lists all events across all projects
type EventsCmd struct{}

type EventsOutput struct {
	Count  int          `json:"count"`
	Events []EventBrief `json:"events"`
}

type EventBrief struct {
	ID             int    `json:"id"`
	Action         string `json:"action"`
	RecordingType  string `json:"recording_type"`
	RecordingID    int    `json:"recording_id"`
	RecordingTitle string `json:"recording_title"`
	ProjectID      int    `json:"project_id"`
	ProjectName    string `json:"project_name"`
	Creator        string `json:"creator"`
	CreatedAt      string `json:"created_at"`
}

func (c *EventsCmd) Run(args []string) error {
	cl, err := client.New()
	if err != nil {
		return err
	}

	data, err := cl.Get("/events.json")
	if err != nil {
		return err
	}

	var events []Event
	if err := json.Unmarshal(data, &events); err != nil {
		return err
	}

	output := EventsOutput{
		Count:  len(events),
		Events: make([]EventBrief, len(events)),
	}

	for i, e := range events {
		output.Events[i] = EventBrief{
			ID:             e.ID,
			Action:         e.Action,
			RecordingType:  e.RecordingType,
			RecordingID:    e.Recording.ID,
			RecordingTitle: e.Recording.Title,
			ProjectID:      e.Bucket.ID,
			ProjectName:    e.Bucket.Name,
			Creator:        e.Creator.Name,
			CreatedAt:      e.CreatedAt,
		}
	}

	return PrintJSON(output)
}

// EventsProjectCmd lists events for a specific project
type EventsProjectCmd struct{}

type EventsProjectOutput struct {
	ProjectID int          `json:"project_id"`
	Count     int          `json:"count"`
	Events    []EventBrief `json:"events"`
}

func (c *EventsProjectCmd) Run(args []string) error {
	projectID, _, err := getProjectID(args)
	if err != nil {
		return err
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	data, err := cl.Get("/buckets/" + projectID + "/events.json")
	if err != nil {
		return err
	}

	var events []Event
	if err := json.Unmarshal(data, &events); err != nil {
		return err
	}

	var pID int
	fmt.Sscanf(projectID, "%d", &pID)

	output := EventsProjectOutput{
		ProjectID: pID,
		Count:     len(events),
		Events:    make([]EventBrief, len(events)),
	}

	for i, e := range events {
		output.Events[i] = EventBrief{
			ID:             e.ID,
			Action:         e.Action,
			RecordingType:  e.RecordingType,
			RecordingID:    e.Recording.ID,
			RecordingTitle: e.Recording.Title,
			ProjectID:      e.Bucket.ID,
			ProjectName:    e.Bucket.Name,
			Creator:        e.Creator.Name,
			CreatedAt:      e.CreatedAt,
		}
	}

	return PrintJSON(output)
}

// EventsRecordingCmd lists events for a specific recording
type EventsRecordingCmd struct{}

type EventsRecordingOutput struct {
	ProjectID   int          `json:"project_id"`
	RecordingID int          `json:"recording_id"`
	Count       int          `json:"count"`
	Events      []EventBrief `json:"events"`
}

func (c *EventsRecordingCmd) Run(args []string) error {
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

	data, err := cl.Get("/buckets/" + projectID + "/recordings/" + recordingID + "/events.json")
	if err != nil {
		return err
	}

	var events []Event
	if err := json.Unmarshal(data, &events); err != nil {
		return err
	}

	var pID, rID int
	fmt.Sscanf(projectID, "%d", &pID)
	fmt.Sscanf(recordingID, "%d", &rID)

	output := EventsRecordingOutput{
		ProjectID:   pID,
		RecordingID: rID,
		Count:       len(events),
		Events:      make([]EventBrief, len(events)),
	}

	for i, e := range events {
		output.Events[i] = EventBrief{
			ID:             e.ID,
			Action:         e.Action,
			RecordingType:  e.RecordingType,
			RecordingID:    e.Recording.ID,
			RecordingTitle: e.Recording.Title,
			ProjectID:      e.Bucket.ID,
			ProjectName:    e.Bucket.Name,
			Creator:        e.Creator.Name,
			CreatedAt:      e.CreatedAt,
		}
	}

	return PrintJSON(output)
}
