package tests

import (
	"fmt"
	"testing"

	"github.com/rzolkos/basecamp-cli/e2e/harness"
)

func TestMessageTypes(t *testing.T) {
	h := harness.New(t)

	t.Run("list message types", func(t *testing.T) {
		result := h.Run("message-types", h.ProjectID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("project_id") == 0 {
			t.Error("expected project_id in response")
		}

		// message_types array should exist
		if _, ok := result.JSON["message_types"]; !ok {
			t.Error("expected message_types array in response")
		}
	})
}

func TestMessageTypeCRUD(t *testing.T) {
	h := harness.New(t)

	var typeID string

	t.Run("create message type", func(t *testing.T) {
		result := h.Run("message-type-create", h.ProjectID, "--name", "Test Type", "--icon", "🧪")

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status ok")
		}
		if result.GetInt("id") == 0 {
			t.Error("expected id in response")
		}
		if result.GetString("name") != "Test Type" {
			t.Error("expected name to match")
		}

		typeID = fmt.Sprintf("%d", result.GetInt("id"))
	})

	t.Run("view message type", func(t *testing.T) {
		if typeID == "" {
			t.Skip("no type created")
		}

		result := h.Run("message-type", h.ProjectID, typeID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("id") == 0 {
			t.Error("expected id in response")
		}
		if result.GetString("name") == "" {
			t.Error("expected name in response")
		}
	})

	t.Run("update message type", func(t *testing.T) {
		if typeID == "" {
			t.Skip("no type created")
		}

		result := h.Run("message-type-update", h.ProjectID, typeID, "--name", "Updated Type", "--icon", "✅")

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status ok")
		}
		if result.GetString("name") != "Updated Type" {
			t.Error("expected name to be updated")
		}
	})

	t.Run("delete message type", func(t *testing.T) {
		if typeID == "" {
			t.Skip("no type created")
		}

		result := h.Run("message-type-delete", h.ProjectID, typeID)

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

	t.Run("missing name and icon for create", func(t *testing.T) {
		result := h.Run("message-type-create", h.ProjectID)

		if result.Success() {
			t.Error("expected failure without --name and --icon")
		}
	})

	t.Run("missing type_id for view", func(t *testing.T) {
		result := h.Run("message-type", h.ProjectID)

		if result.Success() {
			t.Error("expected failure without type_id")
		}
	})
}
