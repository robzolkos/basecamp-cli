package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

// fetchVault gets the vault for a project
func fetchVault(cl *client.Client, projectID string) (ProjectDetail, Vault, error) {
	project, err := fetchProject(cl, projectID)
	if err != nil {
		return ProjectDetail{}, Vault{}, err
	}

	vaultURL, err := getDockURL(project, "vault")
	if err != nil {
		return ProjectDetail{}, Vault{}, err
	}

	vaultData, err := cl.Get(vaultURL)
	if err != nil {
		return ProjectDetail{}, Vault{}, err
	}

	var vault Vault
	if err := json.Unmarshal(vaultData, &vault); err != nil {
		return ProjectDetail{}, Vault{}, err
	}

	return project, vault, nil
}

type DocsCmd struct{}

type Vault struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	DocumentsURL string `json:"documents_url"`
}

type Document struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Content       string  `json:"content"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	CommentsCount int     `json:"comments_count"`
	CommentsURL   string  `json:"comments_url"`
	URL           string  `json:"app_url"`
	Creator       Creator `json:"creator"`
}

type DocsListOutput struct {
	ProjectID int              `json:"project_id"`
	VaultID   int              `json:"vault_id"`
	Documents []DocOutputBrief `json:"documents"`
}

type DocOutputBrief struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Creator   string `json:"creator"`
	UpdatedAt string `json:"updated_at"`
}

func (c *DocsCmd) Run(args []string) error {
	projectID, _, err := getProjectID(args)
	if err != nil {
		return err
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	project, vault, err := fetchVault(cl, projectID)
	if err != nil {
		return err
	}

	// Get documents
	docsData, err := cl.Get(vault.DocumentsURL)
	if err != nil {
		return err
	}

	var documents []Document
	if err := json.Unmarshal(docsData, &documents); err != nil {
		return err
	}

	output := DocsListOutput{
		ProjectID: project.ID,
		VaultID:   vault.ID,
		Documents: make([]DocOutputBrief, len(documents)),
	}

	for i, d := range documents {
		output.Documents[i] = DocOutputBrief{
			ID:        d.ID,
			Title:     d.Title,
			Creator:   d.Creator.Name,
			UpdatedAt: d.UpdatedAt,
		}
	}

	return PrintJSON(output)
}

type DocCmd struct{}

type DocDetailOutput struct {
	ID            int             `json:"id"`
	Title         string          `json:"title"`
	Content       string          `json:"content"`
	Creator       string          `json:"creator"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
	CommentsCount int             `json:"comments_count"`
	URL           string          `json:"url"`
	Comments      []CommentOutput `json:"comments,omitempty"`
}

func (c *DocCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// Parse doc_id and flags
	var docID string
	showComments := false

	for i := 0; i < len(remaining); i++ {
		if remaining[i] == "--comments" {
			showComments = true
		} else if docID == "" {
			docID = remaining[i]
		}
	}

	if docID == "" {
		return errors.New("document_id required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	// Get document directly
	docData, err := cl.Get("/buckets/" + projectID + "/documents/" + docID + ".json")
	if err != nil {
		return err
	}

	var doc Document
	if err := json.Unmarshal(docData, &doc); err != nil {
		return err
	}

	output := DocDetailOutput{
		ID:            doc.ID,
		Title:         doc.Title,
		Content:       stripHTML(doc.Content),
		Creator:       doc.Creator.Name,
		CreatedAt:     doc.CreatedAt,
		UpdatedAt:     doc.UpdatedAt,
		CommentsCount: doc.CommentsCount,
		URL:           doc.URL,
	}

	if showComments && doc.CommentsURL != "" {
		comments, err := fetchComments(cl, doc.CommentsURL)
		if err != nil {
			return err
		}
		output.Comments = comments
	}

	return PrintJSON(output)
}

type DocCreateCmd struct{}

type DocCreateOutput struct {
	Status  string `json:"status"`
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

func (c *DocCreateCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// Parse flags
	var title, content string

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
		}
	}

	if title == "" {
		return errors.New("--title required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	_, vault, err := fetchVault(cl, projectID)
	if err != nil {
		return err
	}

	// Create document
	payload := map[string]any{
		"title":  title,
		"status": "active",
	}
	if content != "" {
		payload["content"] = content
	}

	// POST to documents URL
	docsURL := vault.DocumentsURL
	// Convert from full URL to path
	if idx := strings.Index(docsURL, "/buckets/"); idx != -1 {
		docsURL = docsURL[idx:]
	}

	responseData, err := cl.Post(docsURL, payload)
	if err != nil {
		return err
	}

	var created Document
	if err := json.Unmarshal(responseData, &created); err != nil {
		return err
	}

	return PrintJSON(DocCreateOutput{
		Status:  "ok",
		ID:      created.ID,
		Title:   created.Title,
		Message: fmt.Sprintf("Document '%s' created", created.Title),
	})
}
