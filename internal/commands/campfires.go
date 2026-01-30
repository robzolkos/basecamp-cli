package commands

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

// fetchCampfire gets the campfire for a project
func fetchCampfire(cl *client.Client, projectID string) (ProjectDetail, Campfire, error) {
	project, err := fetchProject(cl, projectID)
	if err != nil {
		return ProjectDetail{}, Campfire{}, err
	}

	campfireURL, err := getDockURL(project, "chat")
	if err != nil {
		return ProjectDetail{}, Campfire{}, err
	}

	campfireData, err := cl.Get(campfireURL)
	if err != nil {
		return ProjectDetail{}, Campfire{}, err
	}

	var campfire Campfire
	if err := json.Unmarshal(campfireData, &campfire); err != nil {
		return ProjectDetail{}, Campfire{}, err
	}

	return project, campfire, nil
}

type CampfireCmd struct{}

type Campfire struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	LinesURL string `json:"lines_url"`
}

type CampfireLine struct {
	ID        int     `json:"id"`
	Content   string  `json:"content"`
	CreatedAt string  `json:"created_at"`
	Creator   Creator `json:"creator"`
}

type CampfireListOutput struct {
	ProjectID  int                 `json:"project_id"`
	CampfireID int                 `json:"campfire_id"`
	Lines      []CampfireLineBrief `json:"lines"`
}

type CampfireLineBrief struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	Creator   string `json:"creator"`
	CreatedAt string `json:"created_at"`
}

func (c *CampfireCmd) Run(args []string) error {
	projectID, _, err := getProjectID(args)
	if err != nil {
		return err
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	project, campfire, err := fetchCampfire(cl, projectID)
	if err != nil {
		return err
	}

	// Get lines
	linesData, err := cl.Get(campfire.LinesURL)
	if err != nil {
		return err
	}

	var lines []CampfireLine
	if err := json.Unmarshal(linesData, &lines); err != nil {
		return err
	}

	output := CampfireListOutput{
		ProjectID:  project.ID,
		CampfireID: campfire.ID,
		Lines:      make([]CampfireLineBrief, len(lines)),
	}

	for i, line := range lines {
		output.Lines[i] = CampfireLineBrief{
			ID:        line.ID,
			Content:   line.Content,
			Creator:   line.Creator.Name,
			CreatedAt: line.CreatedAt,
		}
	}

	return PrintJSON(output)
}

type CampfirePostCmd struct{}

type CampfirePostOutput struct {
	Status  string `json:"status"`
	ID      int    `json:"id"`
	Content string `json:"content"`
	Message string `json:"message"`
}

func (c *CampfirePostCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// Parse --content flag
	var content string

	for i := 0; i < len(remaining); i++ {
		if remaining[i] == "--content" {
			if i+1 < len(remaining) {
				content = remaining[i+1]
				i++
			}
		}
	}

	if content == "" {
		return errors.New("--content required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	_, campfire, err := fetchCampfire(cl, projectID)
	if err != nil {
		return err
	}

	// Post line
	payload := map[string]string{
		"content": content,
	}

	// POST to lines URL
	linesURL := campfire.LinesURL
	// Convert from full URL to path
	if idx := strings.Index(linesURL, "/buckets/"); idx != -1 {
		linesURL = linesURL[idx:]
	}

	responseData, err := cl.Post(linesURL, payload)
	if err != nil {
		return err
	}

	var created CampfireLine
	if err := json.Unmarshal(responseData, &created); err != nil {
		return err
	}

	return PrintJSON(CampfirePostOutput{
		Status:  "ok",
		ID:      created.ID,
		Content: created.Content,
		Message: "Message posted to campfire",
	})
}
