package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

// fetchSchedule gets the schedule for a project
func fetchSchedule(cl *client.Client, projectID string) (ProjectDetail, Schedule, error) {
	project, err := fetchProject(cl, projectID)
	if err != nil {
		return ProjectDetail{}, Schedule{}, err
	}

	scheduleURL, err := getDockURL(project, "schedule")
	if err != nil {
		return ProjectDetail{}, Schedule{}, err
	}

	scheduleData, err := cl.Get(scheduleURL)
	if err != nil {
		return ProjectDetail{}, Schedule{}, err
	}

	var schedule Schedule
	if err := json.Unmarshal(scheduleData, &schedule); err != nil {
		return ProjectDetail{}, Schedule{}, err
	}

	return project, schedule, nil
}

type ScheduleCmd struct{}

type Schedule struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	EntriesURL string `json:"entries_url"`
}

type ScheduleEntry struct {
	ID            int       `json:"id"`
	Summary       string    `json:"summary"`
	Description   string    `json:"description"`
	StartsAt      string    `json:"starts_at"`
	EndsAt        string    `json:"ends_at"`
	AllDay        bool      `json:"all_day"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
	CommentsCount int       `json:"comments_count"`
	CommentsURL   string    `json:"comments_url"`
	URL           string    `json:"app_url"`
	Creator       Creator   `json:"creator"`
	Participants  []Creator `json:"participants"`
}

type ScheduleListOutput struct {
	ProjectID  int                  `json:"project_id"`
	ScheduleID int                  `json:"schedule_id"`
	Entries    []ScheduleEntryBrief `json:"entries"`
}

type ScheduleEntryBrief struct {
	ID       int    `json:"id"`
	Summary  string `json:"summary"`
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
	AllDay   bool   `json:"all_day"`
}

func (c *ScheduleCmd) Run(args []string) error {
	projectID, _, err := getProjectID(args)
	if err != nil {
		return err
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	project, schedule, err := fetchSchedule(cl, projectID)
	if err != nil {
		return err
	}

	// Get entries
	entriesData, err := cl.Get(schedule.EntriesURL)
	if err != nil {
		return err
	}

	var entries []ScheduleEntry
	if err := json.Unmarshal(entriesData, &entries); err != nil {
		return err
	}

	output := ScheduleListOutput{
		ProjectID:  project.ID,
		ScheduleID: schedule.ID,
		Entries:    make([]ScheduleEntryBrief, len(entries)),
	}

	for i, e := range entries {
		output.Entries[i] = ScheduleEntryBrief{
			ID:       e.ID,
			Summary:  e.Summary,
			StartsAt: e.StartsAt,
			EndsAt:   e.EndsAt,
			AllDay:   e.AllDay,
		}
	}

	return PrintJSON(output)
}

type EventCmd struct{}

type EventDetailOutput struct {
	ID            int             `json:"id"`
	Summary       string          `json:"summary"`
	Description   string          `json:"description,omitempty"`
	StartsAt      string          `json:"starts_at"`
	EndsAt        string          `json:"ends_at"`
	AllDay        bool            `json:"all_day"`
	Creator       string          `json:"creator"`
	Participants  []string        `json:"participants,omitempty"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
	CommentsCount int             `json:"comments_count"`
	URL           string          `json:"url"`
	Comments      []CommentOutput `json:"comments,omitempty"`
}

func (c *EventCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// Parse entry_id and flags
	var entryID string
	showComments := false

	for i := 0; i < len(remaining); i++ {
		if remaining[i] == "--comments" {
			showComments = true
		} else if entryID == "" {
			entryID = remaining[i]
		}
	}

	if entryID == "" {
		return errors.New("entry_id required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	// Get entry directly
	entryData, err := cl.Get("/buckets/" + projectID + "/schedule_entries/" + entryID + ".json")
	if err != nil {
		return err
	}

	var entry ScheduleEntry
	if err := json.Unmarshal(entryData, &entry); err != nil {
		return err
	}

	var participants []string
	for _, p := range entry.Participants {
		participants = append(participants, p.Name)
	}

	output := EventDetailOutput{
		ID:            entry.ID,
		Summary:       entry.Summary,
		Description:   stripHTML(entry.Description),
		StartsAt:      entry.StartsAt,
		EndsAt:        entry.EndsAt,
		AllDay:        entry.AllDay,
		Creator:       entry.Creator.Name,
		Participants:  participants,
		CreatedAt:     entry.CreatedAt,
		UpdatedAt:     entry.UpdatedAt,
		CommentsCount: entry.CommentsCount,
		URL:           entry.URL,
	}

	if showComments && entry.CommentsURL != "" {
		comments, err := fetchComments(cl, entry.CommentsURL)
		if err != nil {
			return err
		}
		output.Comments = comments
	}

	return PrintJSON(output)
}

type EventCreateCmd struct{}

type EventCreateOutput struct {
	Status  string `json:"status"`
	ID      int    `json:"id"`
	Summary string `json:"summary"`
	Message string `json:"message"`
}

func (c *EventCreateCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// Parse flags
	var summary, description, startsAt, endsAt string
	allDay := false

	for i := 0; i < len(remaining); i++ {
		switch remaining[i] {
		case "--summary":
			if i+1 < len(remaining) {
				summary = remaining[i+1]
				i++
			}
		case "--description":
			if i+1 < len(remaining) {
				description = remaining[i+1]
				i++
			}
		case "--starts-at":
			if i+1 < len(remaining) {
				startsAt = remaining[i+1]
				i++
			}
		case "--ends-at":
			if i+1 < len(remaining) {
				endsAt = remaining[i+1]
				i++
			}
		case "--all-day":
			allDay = true
		}
	}

	if summary == "" {
		return errors.New("--summary required")
	}
	if startsAt == "" {
		return errors.New("--starts-at required")
	}
	if endsAt == "" {
		return errors.New("--ends-at required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	_, schedule, err := fetchSchedule(cl, projectID)
	if err != nil {
		return err
	}

	// Create entry
	payload := map[string]any{
		"summary":   summary,
		"starts_at": startsAt,
		"ends_at":   endsAt,
		"all_day":   allDay,
	}
	if description != "" {
		payload["description"] = description
	}

	// POST to entries URL
	entriesURL := schedule.EntriesURL
	// Convert from full URL to path
	if idx := strings.Index(entriesURL, "/buckets/"); idx != -1 {
		entriesURL = entriesURL[idx:]
	}

	responseData, err := cl.Post(entriesURL, payload)
	if err != nil {
		return err
	}

	var created ScheduleEntry
	if err := json.Unmarshal(responseData, &created); err != nil {
		return err
	}

	return PrintJSON(EventCreateOutput{
		Status:  "ok",
		ID:      created.ID,
		Summary: created.Summary,
		Message: fmt.Sprintf("Event '%s' created", created.Summary),
	})
}
