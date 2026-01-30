package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

// Upload represents a Basecamp upload (file in a vault)
type Upload struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	ContentType   string  `json:"content_type"`
	ByteSize      int64   `json:"byte_size"`
	Width         int     `json:"width"`
	Height        int     `json:"height"`
	DownloadURL   string  `json:"download_url"`
	AppURL        string  `json:"app_url"`
	CommentsCount int     `json:"comments_count"`
	CommentsURL   string  `json:"comments_url"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	Creator       Creator `json:"creator"`
}

// UploadCmd uploads a file and returns the attachable_sgid
type UploadCmd struct{}

type UploadOutput struct {
	Status         string `json:"status"`
	AttachableSGID string `json:"attachable_sgid"`
	Message        string `json:"message"`
}

func (c *UploadCmd) Run(args []string) error {
	if len(args) < 1 {
		return errors.New("file path required")
	}
	filePath := args[0]

	// Read the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// Determine content type
	ext := filepath.Ext(filePath)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Read file contents
	fileData, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	// Upload the file
	fileName := filepath.Base(filePath)
	data, err := cl.UploadFile("/attachments.json?name="+fileName, fileData, contentType, fileInfo.Size())
	if err != nil {
		return err
	}

	var result struct {
		AttachableSGID string `json:"attachable_sgid"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	return PrintJSON(UploadOutput{
		Status:         "ok",
		AttachableSGID: result.AttachableSGID,
		Message:        fmt.Sprintf("File '%s' uploaded", fileName),
	})
}

// UploadsCmd lists uploads in a vault
type UploadsCmd struct{}

type UploadsOutput struct {
	ProjectID int           `json:"project_id"`
	VaultID   int           `json:"vault_id"`
	Uploads   []UploadBrief `json:"uploads"`
}

type UploadBrief struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ContentType string `json:"content_type"`
	ByteSize    int64  `json:"byte_size"`
	Creator     string `json:"creator"`
	CreatedAt   string `json:"created_at"`
}

func (c *UploadsCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("vault_id required")
	}
	vaultID := remaining[0]

	cl, err := client.New()
	if err != nil {
		return err
	}

	data, err := cl.Get("/buckets/" + projectID + "/vaults/" + vaultID + "/uploads.json")
	if err != nil {
		return err
	}

	var uploads []Upload
	if err := json.Unmarshal(data, &uploads); err != nil {
		return err
	}

	var pID, vID int
	fmt.Sscanf(projectID, "%d", &pID)
	fmt.Sscanf(vaultID, "%d", &vID)

	output := UploadsOutput{
		ProjectID: pID,
		VaultID:   vID,
		Uploads:   make([]UploadBrief, len(uploads)),
	}

	for i, u := range uploads {
		output.Uploads[i] = UploadBrief{
			ID:          u.ID,
			Title:       u.Title,
			ContentType: u.ContentType,
			ByteSize:    u.ByteSize,
			Creator:     u.Creator.Name,
			CreatedAt:   u.CreatedAt,
		}
	}

	return PrintJSON(output)
}

// UploadViewCmd views a single upload
type UploadViewCmd struct{}

type UploadDetailOutput struct {
	ID            int             `json:"id"`
	Title         string          `json:"title"`
	Description   string          `json:"description,omitempty"`
	ContentType   string          `json:"content_type"`
	ByteSize      int64           `json:"byte_size"`
	Width         int             `json:"width,omitempty"`
	Height        int             `json:"height,omitempty"`
	DownloadURL   string          `json:"download_url"`
	URL           string          `json:"url"`
	Creator       string          `json:"creator"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
	CommentsCount int             `json:"comments_count"`
	Comments      []CommentOutput `json:"comments,omitempty"`
}

func (c *UploadViewCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// Parse upload_id and flags
	var uploadID string
	showComments := false

	for i := 0; i < len(remaining); i++ {
		if remaining[i] == "--comments" {
			showComments = true
		} else if uploadID == "" {
			uploadID = remaining[i]
		}
	}

	if uploadID == "" {
		return errors.New("upload_id required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	data, err := cl.Get("/buckets/" + projectID + "/uploads/" + uploadID + ".json")
	if err != nil {
		return err
	}

	var upload Upload
	if err := json.Unmarshal(data, &upload); err != nil {
		return err
	}

	output := UploadDetailOutput{
		ID:            upload.ID,
		Title:         upload.Title,
		Description:   stripHTML(upload.Description),
		ContentType:   upload.ContentType,
		ByteSize:      upload.ByteSize,
		Width:         upload.Width,
		Height:        upload.Height,
		DownloadURL:   upload.DownloadURL,
		URL:           upload.AppURL,
		Creator:       upload.Creator.Name,
		CreatedAt:     upload.CreatedAt,
		UpdatedAt:     upload.UpdatedAt,
		CommentsCount: upload.CommentsCount,
	}

	if showComments && upload.CommentsURL != "" {
		comments, err := fetchComments(cl, upload.CommentsURL)
		if err != nil {
			return err
		}
		output.Comments = comments
	}

	return PrintJSON(output)
}
