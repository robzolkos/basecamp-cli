package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/rzolkos/basecamp-cli/e2e/harness"
)

func TestColumns(t *testing.T) {
	h := harness.New(t)

	if h.BoardID == "" {
		t.Skip("BASECAMP_TEST_BOARD_ID not set")
	}

	t.Run("list columns", func(t *testing.T) {
		result := h.Run("columns", h.ProjectID, h.BoardID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("board_id") == 0 {
			t.Error("expected board_id in response")
		}

		columns, ok := result.JSON["columns"].([]any)
		if !ok || len(columns) == 0 {
			t.Error("expected at least one column")
		}
	})
}

func TestCardCreate(t *testing.T) {
	h := harness.New(t)

	if h.BoardID == "" {
		t.Skip("BASECAMP_TEST_BOARD_ID not set")
	}

	// First get a column ID
	columnsResult := h.Run("columns", h.ProjectID, h.BoardID)
	if !columnsResult.Success() {
		t.Fatalf("failed to get columns: %s", columnsResult.Stderr)
	}

	columns, ok := columnsResult.JSON["columns"].([]any)
	if !ok || len(columns) == 0 {
		t.Fatal("no columns found")
	}

	firstColumn := columns[0].(map[string]any)
	columnID := fmt.Sprintf("%.0f", firstColumn["id"].(float64))

	var cardID string
	cardTitle := fmt.Sprintf("E2E Test Card %d", time.Now().UnixNano())

	t.Run("create card", func(t *testing.T) {
		result := h.Run("card-create", h.ProjectID, h.BoardID, "--column", columnID, "--title", cardTitle)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status=ok")
		}

		cardID = fmt.Sprintf("%d", result.GetInt("id"))
		if cardID == "0" {
			t.Fatal("expected card id in response")
		}
	})

	t.Run("update card", func(t *testing.T) {
		if cardID == "" {
			t.Skip("no card created")
		}

		newTitle := fmt.Sprintf("Updated Card %d", time.Now().UnixNano())
		result := h.Run("card-update", h.ProjectID, cardID, "--title", newTitle)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status=ok")
		}
	})

	t.Run("missing required flags", func(t *testing.T) {
		result := h.Run("card-create", h.ProjectID, h.BoardID)

		if result.Success() {
			t.Error("expected failure without required flags")
		}
	})
}
