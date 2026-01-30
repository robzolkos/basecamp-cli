package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/rzolkos/basecamp-cli/e2e/harness"
)

func TestTodolists(t *testing.T) {
	h := harness.New(t)

	t.Run("list todolists", func(t *testing.T) {
		result := h.Run("todolists", h.ProjectID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("project_id") == 0 {
			t.Error("expected project_id in response")
		}
		if result.GetInt("todoset_id") == 0 {
			t.Error("expected todoset_id in response")
		}
	})
}

func TestTodos(t *testing.T) {
	h := harness.New(t)

	if h.TodolistID == "" {
		t.Skip("BASECAMP_TEST_TODOLIST_ID not set")
	}

	t.Run("list todos", func(t *testing.T) {
		result := h.Run("todos", h.ProjectID, h.TodolistID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("todolist_id") == 0 {
			t.Error("expected todolist_id in response")
		}
	})
}

func TestTodoCRUD(t *testing.T) {
	h := harness.New(t)

	if h.TodolistID == "" {
		t.Skip("BASECAMP_TEST_TODOLIST_ID not set")
	}

	var todoID string
	todoContent := fmt.Sprintf("E2E Test Todo %d", time.Now().UnixNano())

	t.Run("create todo", func(t *testing.T) {
		result := h.Run("todo-create", h.ProjectID, h.TodolistID, "--content", todoContent)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status=ok")
		}

		todoID = fmt.Sprintf("%d", result.GetInt("id"))
		if todoID == "0" {
			t.Fatal("expected todo id in response")
		}
	})

	t.Run("view todo", func(t *testing.T) {
		if todoID == "" {
			t.Skip("no todo created")
		}

		result := h.Run("todo", h.ProjectID, todoID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("id") == 0 {
			t.Error("expected id in response")
		}
		if result.GetString("content") == "" {
			t.Error("expected content in response")
		}
	})

	t.Run("complete todo", func(t *testing.T) {
		if todoID == "" {
			t.Skip("no todo created")
		}

		result := h.Run("todo-complete", h.ProjectID, todoID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status=ok")
		}
	})

	t.Run("verify todo is completed", func(t *testing.T) {
		if todoID == "" {
			t.Skip("no todo created")
		}

		result := h.Run("todo", h.ProjectID, todoID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		completed, ok := result.JSON["completed"].(bool)
		if !ok || !completed {
			t.Error("expected todo to be completed")
		}
	})

	t.Run("uncomplete todo", func(t *testing.T) {
		if todoID == "" {
			t.Skip("no todo created")
		}

		result := h.Run("todo-uncomplete", h.ProjectID, todoID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status=ok")
		}
	})

	t.Run("verify todo is uncompleted", func(t *testing.T) {
		if todoID == "" {
			t.Skip("no todo created")
		}

		result := h.Run("todo", h.ProjectID, todoID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		completed, ok := result.JSON["completed"].(bool)
		if !ok || completed {
			t.Error("expected todo to be uncompleted")
		}
	})

	// Note: Basecamp API doesn't have a delete endpoint for todos,
	// they can only be trashed via the recordings API which is out of scope.
	// The test todo will remain in the todolist.
}
