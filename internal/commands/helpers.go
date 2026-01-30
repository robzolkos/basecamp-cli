package commands

import (
	"encoding/json"
	"fmt"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

// getDockURL finds a dock item URL by name from a project
func getDockURL(project ProjectDetail, dockName string) (string, error) {
	for _, dock := range project.Dock {
		if dock.Name == dockName {
			return dock.URL, nil
		}
	}
	return "", fmt.Errorf("no %s found in this project", dockName)
}

// fetchProject gets a project by ID and returns the parsed ProjectDetail
func fetchProject(cl *client.Client, projectID string) (ProjectDetail, error) {
	data, err := cl.Get("/projects/" + projectID + ".json")
	if err != nil {
		return ProjectDetail{}, err
	}

	var project ProjectDetail
	if err := json.Unmarshal(data, &project); err != nil {
		return ProjectDetail{}, err
	}
	return project, nil
}

// fetchComments fetches and parses comments from a comments URL
func fetchComments(cl *client.Client, commentsURL string) ([]CommentOutput, error) {
	commentsData, err := cl.GetAll(commentsURL)
	if err != nil {
		return nil, err
	}

	comments := make([]CommentOutput, len(commentsData))
	for i, commentJSON := range commentsData {
		var comment Comment
		if err := json.Unmarshal(commentJSON, &comment); err != nil {
			return nil, err
		}
		author := "Unknown"
		if comment.Creator.Name != "" {
			author = comment.Creator.Name
		}
		comments[i] = CommentOutput{
			ID:        comment.ID,
			Author:    author,
			Content:   stripHTML(comment.Content),
			CreatedAt: comment.CreatedAt,
		}
	}
	return comments, nil
}
