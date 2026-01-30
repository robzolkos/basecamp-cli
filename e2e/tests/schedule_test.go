package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/rzolkos/basecamp-cli/e2e/harness"
)

func TestSchedule(t *testing.T) {
	h := harness.New(t)

	t.Run("list schedule entries", func(t *testing.T) {
		result := h.Run("schedule", h.ProjectID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("project_id") == 0 {
			t.Error("expected project_id in response")
		}
		if result.GetInt("schedule_id") == 0 {
			t.Error("expected schedule_id in response")
		}
	})
}

func TestEventCRUD(t *testing.T) {
	h := harness.New(t)

	var eventID string
	eventSummary := fmt.Sprintf("E2E Test Event %d", time.Now().UnixNano())
	startsAt := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	endsAt := time.Now().Add(25 * time.Hour).Format(time.RFC3339)

	t.Run("create event", func(t *testing.T) {
		result := h.Run("event-create", h.ProjectID,
			"--summary", eventSummary,
			"--starts-at", startsAt,
			"--ends-at", endsAt)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status=ok")
		}

		eventID = fmt.Sprintf("%d", result.GetInt("id"))
		if eventID == "0" {
			t.Fatal("expected event id in response")
		}
	})

	t.Run("view event", func(t *testing.T) {
		if eventID == "" {
			t.Skip("no event created")
		}

		result := h.Run("event", h.ProjectID, eventID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("id") == 0 {
			t.Error("expected id in response")
		}
		if result.GetString("summary") == "" {
			t.Error("expected summary in response")
		}
	})

	t.Run("view event with comments flag", func(t *testing.T) {
		if eventID == "" {
			t.Skip("no event created")
		}

		result := h.Run("event", h.ProjectID, eventID, "--comments")

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

	t.Run("missing required flags", func(t *testing.T) {
		result := h.Run("event-create", h.ProjectID)

		if result.Success() {
			t.Error("expected failure without required flags")
		}

		if result.ErrorMessage() != "--summary required" {
			t.Errorf("expected '--summary required' error, got: %s", result.ErrorMessage())
		}
	})
}
