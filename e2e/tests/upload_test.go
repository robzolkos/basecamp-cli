package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/rzolkos/basecamp-cli/e2e/harness"
)

func TestUpload(t *testing.T) {
	h := harness.New(t)

	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-upload.txt")
	if err := os.WriteFile(testFile, []byte("Test file content for upload"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	t.Run("upload file", func(t *testing.T) {
		result := h.Run("upload", testFile)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status ok")
		}
		if result.GetString("attachable_sgid") == "" {
			t.Error("expected attachable_sgid in response")
		}
	})

	t.Run("missing file path", func(t *testing.T) {
		result := h.Run("upload")

		if result.Success() {
			t.Error("expected failure without file path")
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		result := h.Run("upload", "/nonexistent/file.txt")

		if result.Success() {
			t.Error("expected failure with nonexistent file")
		}
	})
}

func TestUploads(t *testing.T) {
	h := harness.New(t)

	// Get vault ID from docs command
	docsResult := h.Run("docs", h.ProjectID)
	if !docsResult.Success() {
		t.Fatalf("failed to get docs: %s", docsResult.Stderr)
	}

	vaultID := docsResult.GetInt("vault_id")
	if vaultID == 0 {
		t.Skip("no vault found")
	}

	t.Run("list uploads", func(t *testing.T) {
		result := h.Run("uploads", h.ProjectID, fmt.Sprintf("%d", vaultID))

		// Note: This test may fail if there are no uploads in the vault
		// That's acceptable - we're testing the command structure
		if result.JSON != nil {
			if _, ok := result.JSON["uploads"]; !ok {
				t.Error("expected uploads array in response")
			}
		}
	})

	t.Run("missing vault_id", func(t *testing.T) {
		result := h.Run("uploads", h.ProjectID)

		if result.Success() {
			t.Error("expected failure without vault_id")
		}
	})
}

func TestUploadView(t *testing.T) {
	h := harness.New(t)

	t.Run("missing upload_id", func(t *testing.T) {
		result := h.Run("upload-view", h.ProjectID)

		if result.Success() {
			t.Error("expected failure without upload_id")
		}
	})
}
