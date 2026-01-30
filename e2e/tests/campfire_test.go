package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/rzolkos/basecamp-cli/e2e/harness"
)

func TestCampfire(t *testing.T) {
	h := harness.New(t)

	t.Run("list campfire messages", func(t *testing.T) {
		result := h.Run("campfire", h.ProjectID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("project_id") == 0 {
			t.Error("expected project_id in response")
		}
		if result.GetInt("campfire_id") == 0 {
			t.Error("expected campfire_id in response")
		}
	})
}

func TestCampfirePost(t *testing.T) {
	h := harness.New(t)

	t.Run("post to campfire", func(t *testing.T) {
		content := fmt.Sprintf("E2E Test Message %d", time.Now().UnixNano())
		result := h.Run("campfire-post", h.ProjectID, "--content", content)

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
			t.Error("expected id in response")
		}
	})

	t.Run("missing content flag", func(t *testing.T) {
		result := h.Run("campfire-post", h.ProjectID)

		if result.Success() {
			t.Error("expected failure without --content flag")
		}

		if result.ErrorMessage() != "--content required" {
			t.Errorf("expected '--content required' error, got: %s", result.ErrorMessage())
		}
	})
}
