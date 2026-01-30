package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/rzolkos/basecamp-cli/e2e/harness"
)

func TestStepCRUD(t *testing.T) {
	h := harness.New(t)

	if h.CardID == "" {
		t.Skip("BASECAMP_TEST_CARD_ID not set")
	}

	var stepID string
	stepTitle := fmt.Sprintf("E2E Test Step %d", time.Now().UnixNano())

	t.Run("create step", func(t *testing.T) {
		result := h.Run("step-create", h.ProjectID, h.CardID, "--title", stepTitle)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status=ok")
		}

		stepID = fmt.Sprintf("%d", result.GetInt("id"))
		if stepID == "0" {
			t.Fatal("expected step id in response")
		}
	})

	t.Run("view card shows step", func(t *testing.T) {
		if stepID == "" {
			t.Skip("no step created")
		}

		result := h.Run("card", h.ProjectID, h.CardID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		steps, ok := result.JSON["steps"].([]any)
		if !ok || len(steps) == 0 {
			t.Error("expected steps in card output")
		}
	})

	t.Run("update step", func(t *testing.T) {
		if stepID == "" {
			t.Skip("no step created")
		}

		newTitle := fmt.Sprintf("Updated Step %d", time.Now().UnixNano())
		result := h.Run("step-update", h.ProjectID, stepID, "--title", newTitle)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status=ok")
		}
	})

	t.Run("complete step", func(t *testing.T) {
		if stepID == "" {
			t.Skip("no step created")
		}

		result := h.Run("step-complete", h.ProjectID, stepID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status=ok")
		}
	})

	t.Run("uncomplete step", func(t *testing.T) {
		if stepID == "" {
			t.Skip("no step created")
		}

		result := h.Run("step-uncomplete", h.ProjectID, stepID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status=ok")
		}
	})

	t.Run("missing title flag", func(t *testing.T) {
		result := h.Run("step-create", h.ProjectID, h.CardID)

		if result.Success() {
			t.Error("expected failure without --title flag")
		}

		if result.ErrorMessage() != "--title required" {
			t.Errorf("expected '--title required' error, got: %s", result.ErrorMessage())
		}
	})
}

func TestStepReposition(t *testing.T) {
	h := harness.New(t)

	if h.CardID == "" {
		t.Skip("BASECAMP_TEST_CARD_ID not set")
	}

	// Create two steps to test repositioning
	step1Title := fmt.Sprintf("E2E Reposition Step 1 %d", time.Now().UnixNano())
	step2Title := fmt.Sprintf("E2E Reposition Step 2 %d", time.Now().UnixNano())

	result1 := h.Run("step-create", h.ProjectID, h.CardID, "--title", step1Title)
	if !result1.Success() {
		t.Fatalf("failed to create step 1: %s", result1.Stderr)
	}
	step1ID := fmt.Sprintf("%d", result1.GetInt("id"))

	result2 := h.Run("step-create", h.ProjectID, h.CardID, "--title", step2Title)
	if !result2.Success() {
		t.Fatalf("failed to create step 2: %s", result2.Stderr)
	}
	step2ID := fmt.Sprintf("%d", result2.GetInt("id"))

	t.Run("reposition step", func(t *testing.T) {
		// Move step2 to position 0 (first)
		result := h.Run("step-reposition", h.ProjectID, h.CardID, step2ID, "--position", "0")

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status=ok")
		}
	})

	t.Run("missing position flag", func(t *testing.T) {
		result := h.Run("step-reposition", h.ProjectID, h.CardID, step1ID)

		if result.Success() {
			t.Error("expected failure without --position flag")
		}

		if result.ErrorMessage() != "--position required" {
			t.Errorf("expected '--position required' error, got: %s", result.ErrorMessage())
		}
	})
}
