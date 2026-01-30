package tests

import (
	"fmt"
	"testing"

	"github.com/rzolkos/basecamp-cli/e2e/harness"
)

func TestRecordingStatus(t *testing.T) {
	h := harness.New(t)

	// Create a todo to use as a test recording
	createResult := h.Run("todo-create", h.ProjectID, h.TodolistID, "--content", "Recording status test todo")
	if !createResult.Success() {
		t.Fatalf("failed to create todo: %s", createResult.Stderr)
	}

	recordingID := fmt.Sprintf("%d", createResult.GetInt("id"))

	t.Run("archive recording", func(t *testing.T) {
		result := h.Run("archive", h.ProjectID, recordingID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status ok")
		}
		if result.GetString("recording_id") != recordingID {
			t.Error("expected recording_id to match")
		}
	})

	t.Run("unarchive recording", func(t *testing.T) {
		result := h.Run("unarchive", h.ProjectID, recordingID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status ok")
		}
	})

	t.Run("trash recording", func(t *testing.T) {
		result := h.Run("trash", h.ProjectID, recordingID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status ok")
		}
	})

	t.Run("missing recording_id for archive", func(t *testing.T) {
		result := h.Run("archive", h.ProjectID)

		if result.Success() {
			t.Error("expected failure without recording_id")
		}
	})

	t.Run("missing recording_id for unarchive", func(t *testing.T) {
		result := h.Run("unarchive", h.ProjectID)

		if result.Success() {
			t.Error("expected failure without recording_id")
		}
	})

	t.Run("missing recording_id for trash", func(t *testing.T) {
		result := h.Run("trash", h.ProjectID)

		if result.Success() {
			t.Error("expected failure without recording_id")
		}
	})
}
