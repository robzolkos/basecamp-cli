package tests

import (
	"fmt"
	"testing"

	"github.com/rzolkos/basecamp-cli/e2e/harness"
)

func TestTodolistGroups(t *testing.T) {
	h := harness.New(t)

	t.Run("list todolist groups", func(t *testing.T) {
		result := h.Run("todolist-groups", h.ProjectID, h.TodolistID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("project_id") == 0 {
			t.Error("expected project_id in response")
		}
		if result.GetInt("todolist_id") == 0 {
			t.Error("expected todolist_id in response")
		}

		// Groups array should exist (may be empty)
		if _, ok := result.JSON["groups"]; !ok {
			t.Error("expected groups array in response")
		}
	})

	t.Run("missing todolist_id", func(t *testing.T) {
		result := h.Run("todolist-groups", h.ProjectID)

		if result.Success() {
			t.Error("expected failure without todolist_id")
		}
	})
}

func TestTodolistGroupCRUD(t *testing.T) {
	h := harness.New(t)

	var groupID string

	t.Run("create group", func(t *testing.T) {
		result := h.Run("todolist-group-create", h.ProjectID, h.TodolistID, "--name", "Test Group")

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
		if result.GetString("name") != "Test Group" {
			t.Error("expected name to match")
		}

		groupID = fmt.Sprintf("%d", result.GetInt("id"))
	})

	t.Run("view group", func(t *testing.T) {
		if groupID == "" {
			t.Skip("no group created")
		}

		result := h.Run("todolist-group", h.ProjectID, groupID)

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

	t.Run("missing todolist_id", func(t *testing.T) {
		result := h.Run("todolist-group-create", h.ProjectID)

		if result.Success() {
			t.Error("expected failure without todolist_id")
		}
	})

	t.Run("missing name flag", func(t *testing.T) {
		result := h.Run("todolist-group-create", h.ProjectID, h.TodolistID)

		if result.Success() {
			t.Error("expected failure without --name")
		}
	})

	t.Run("missing group_id", func(t *testing.T) {
		result := h.Run("todolist-group", h.ProjectID)

		if result.Success() {
			t.Error("expected failure without group_id")
		}
	})
}

func TestTodoReposition(t *testing.T) {
	h := harness.New(t)

	// First create a todo to reposition
	createResult := h.Run("todo-create", h.ProjectID, h.TodolistID, "--content", "Reposition test todo")
	if !createResult.Success() {
		t.Fatalf("failed to create todo: %s", createResult.Stderr)
	}

	todoID := fmt.Sprintf("%d", createResult.GetInt("id"))

	t.Run("reposition todo", func(t *testing.T) {
		result := h.Run("todo-reposition", h.ProjectID, todoID, "--position", "1")

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status ok")
		}
		if result.GetString("position") != "1" {
			t.Error("expected position 1")
		}
	})

	t.Run("missing position flag", func(t *testing.T) {
		result := h.Run("todo-reposition", h.ProjectID, todoID)

		if result.Success() {
			t.Error("expected failure without --position")
		}
	})

	t.Run("missing todo_id", func(t *testing.T) {
		result := h.Run("todo-reposition", h.ProjectID)

		if result.Success() {
			t.Error("expected failure without todo_id")
		}
	})
}
