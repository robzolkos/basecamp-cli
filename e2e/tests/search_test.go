package tests

import (
	"testing"

	"github.com/rzolkos/basecamp-cli/e2e/harness"
)

func TestSearch(t *testing.T) {
	h := harness.New(t)

	t.Run("search with query", func(t *testing.T) {
		result := h.Run("search", "test")

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetString("query") != "test" {
			t.Error("expected query in response")
		}
	})

	t.Run("search with type filter", func(t *testing.T) {
		result := h.Run("search", "test", "--type", "Todo")

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		// All results should be todos
		results, ok := result.JSON["results"].([]any)
		if ok {
			for _, r := range results {
				rMap := r.(map[string]any)
				if rMap["type"] != "Todo" {
					t.Errorf("expected all results to be Todo, got %s", rMap["type"])
				}
			}
		}
	})

	t.Run("search with project filter", func(t *testing.T) {
		result := h.Run("search", "test", "--project", h.ProjectID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}
	})

	t.Run("missing query", func(t *testing.T) {
		result := h.Run("search")

		if result.Success() {
			t.Error("expected failure without query")
		}

		if result.ErrorMessage() != "search query required" {
			t.Errorf("expected 'search query required' error, got: %s", result.ErrorMessage())
		}
	})
}
