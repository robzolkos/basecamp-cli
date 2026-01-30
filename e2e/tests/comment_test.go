package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/rzolkos/basecamp-cli/e2e/harness"
)

func TestCommentAdd(t *testing.T) {
	h := harness.New(t)

	if h.TodolistID == "" {
		t.Skip("BASECAMP_TEST_TODOLIST_ID not set")
	}

	// First create a todo to comment on
	todoContent := fmt.Sprintf("E2E Comment Test Todo %d", time.Now().UnixNano())
	createResult := h.Run("todo-create", h.ProjectID, h.TodolistID, "--content", todoContent)

	if !createResult.Success() {
		t.Fatalf("failed to create test todo: %s", createResult.Stderr)
	}

	todoID := fmt.Sprintf("%d", createResult.GetInt("id"))
	if todoID == "0" {
		t.Fatal("failed to get todo ID")
	}

	t.Run("add comment to todo", func(t *testing.T) {
		commentContent := fmt.Sprintf("E2E Test Comment %d", time.Now().UnixNano())
		result := h.Run("comment-add", h.ProjectID, todoID, "--content", commentContent)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status=ok")
		}

		if result.GetInt("id") == 0 {
			t.Error("expected comment id in response")
		}

		if result.GetString("recording_id") != todoID {
			t.Errorf("expected recording_id=%s, got %s", todoID, result.GetString("recording_id"))
		}
	})

	t.Run("missing content flag", func(t *testing.T) {
		result := h.Run("comment-add", h.ProjectID, todoID)

		if result.Success() {
			t.Error("expected failure without --content flag")
		}

		if result.ErrorMessage() != "--content required" {
			t.Errorf("expected '--content required' error, got: %s", result.ErrorMessage())
		}
	})
}
