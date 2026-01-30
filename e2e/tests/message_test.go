package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/rzolkos/basecamp-cli/e2e/harness"
)

func TestMessages(t *testing.T) {
	h := harness.New(t)

	t.Run("list messages", func(t *testing.T) {
		result := h.Run("messages", h.ProjectID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("project_id") == 0 {
			t.Error("expected project_id in response")
		}
		if result.GetInt("message_board_id") == 0 {
			t.Error("expected message_board_id in response")
		}
	})
}

func TestMessageCRUD(t *testing.T) {
	h := harness.New(t)

	var messageID string
	messageSubject := fmt.Sprintf("E2E Test Message %d", time.Now().UnixNano())

	t.Run("create message", func(t *testing.T) {
		result := h.Run("message-create", h.ProjectID, "--subject", messageSubject, "--content", "Test message content from e2e tests")

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status=ok")
		}

		messageID = fmt.Sprintf("%d", result.GetInt("id"))
		if messageID == "0" {
			t.Fatal("expected message id in response")
		}
	})

	t.Run("view message", func(t *testing.T) {
		if messageID == "" {
			t.Skip("no message created")
		}

		result := h.Run("message", h.ProjectID, messageID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("id") == 0 {
			t.Error("expected id in response")
		}
		if result.GetString("subject") == "" {
			t.Error("expected subject in response")
		}
	})

	t.Run("view message with comments flag", func(t *testing.T) {
		if messageID == "" {
			t.Skip("no message created")
		}

		result := h.Run("message", h.ProjectID, messageID, "--comments")

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
}
