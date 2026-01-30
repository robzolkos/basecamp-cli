package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/rzolkos/basecamp-cli/e2e/harness"
)

func TestDocs(t *testing.T) {
	h := harness.New(t)

	t.Run("list documents", func(t *testing.T) {
		result := h.Run("docs", h.ProjectID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("project_id") == 0 {
			t.Error("expected project_id in response")
		}
		if result.GetInt("vault_id") == 0 {
			t.Error("expected vault_id in response")
		}
	})
}

func TestDocCRUD(t *testing.T) {
	h := harness.New(t)

	var docID string
	docTitle := fmt.Sprintf("E2E Test Document %d", time.Now().UnixNano())

	t.Run("create document", func(t *testing.T) {
		result := h.Run("doc-create", h.ProjectID, "--title", docTitle, "--content", "Test document content from e2e tests")

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status=ok")
		}

		docID = fmt.Sprintf("%d", result.GetInt("id"))
		if docID == "0" {
			t.Fatal("expected document id in response")
		}
	})

	t.Run("view document", func(t *testing.T) {
		if docID == "" {
			t.Skip("no document created")
		}

		result := h.Run("doc", h.ProjectID, docID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("id") == 0 {
			t.Error("expected id in response")
		}
		if result.GetString("title") == "" {
			t.Error("expected title in response")
		}
	})

	t.Run("view document with comments flag", func(t *testing.T) {
		if docID == "" {
			t.Skip("no document created")
		}

		result := h.Run("doc", h.ProjectID, docID, "--comments")

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("id") == 0 {
			t.Error("expected id in response")
		}
	})

	t.Run("missing title flag", func(t *testing.T) {
		result := h.Run("doc-create", h.ProjectID)

		if result.Success() {
			t.Error("expected failure without --title flag")
		}

		if result.ErrorMessage() != "--title required" {
			t.Errorf("expected '--title required' error, got: %s", result.ErrorMessage())
		}
	})
}
