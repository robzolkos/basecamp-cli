package tests

import (
	"testing"

	"github.com/rzolkos/basecamp-cli/e2e/harness"
)

func TestEvents(t *testing.T) {
	h := harness.New(t)

	t.Run("list all events", func(t *testing.T) {
		result := h.Run("events")

		// This endpoint may not be available in all accounts
		if !result.Success() {
			errMsg := result.ErrorMessage()
			if errMsg != "" {
				t.Skip("events endpoint not available")
			}
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		// count should exist
		if _, ok := result.JSON["count"]; !ok {
			t.Error("expected count in response")
		}

		// events array should exist
		if _, ok := result.JSON["events"]; !ok {
			t.Error("expected events array in response")
		}
	})
}

func TestEventsProject(t *testing.T) {
	h := harness.New(t)

	t.Run("list project events", func(t *testing.T) {
		result := h.Run("events-project", h.ProjectID)

		// This endpoint may not be available in all accounts
		if !result.Success() {
			errMsg := result.ErrorMessage()
			if errMsg != "" {
				t.Skip("project events endpoint not available")
			}
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("project_id") == 0 {
			t.Error("expected project_id in response")
		}

		// events array should exist
		if _, ok := result.JSON["events"]; !ok {
			t.Error("expected events array in response")
		}
	})
}

func TestEventsRecording(t *testing.T) {
	h := harness.New(t)

	t.Run("list recording events", func(t *testing.T) {
		// Use the test card as the recording
		result := h.Run("events-recording", h.ProjectID, h.CardID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("project_id") == 0 {
			t.Error("expected project_id in response")
		}
		if result.GetInt("recording_id") == 0 {
			t.Error("expected recording_id in response")
		}

		// events array should exist
		if _, ok := result.JSON["events"]; !ok {
			t.Error("expected events array in response")
		}
	})

	t.Run("missing recording_id", func(t *testing.T) {
		result := h.Run("events-recording", h.ProjectID)

		if result.Success() {
			t.Error("expected failure without recording_id")
		}
	})
}
