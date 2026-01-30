package tests

import (
	"fmt"
	"testing"

	"github.com/rzolkos/basecamp-cli/e2e/harness"
)

func TestPeople(t *testing.T) {
	h := harness.New(t)

	t.Run("list all people", func(t *testing.T) {
		result := h.Run("people")

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		count := result.GetInt("count")
		if count == 0 {
			t.Error("expected at least one person")
		}

		people, ok := result.JSON["people"].([]any)
		if !ok {
			t.Fatal("expected people array")
		}

		if len(people) != count {
			t.Errorf("expected %d people, got %d", count, len(people))
		}

		// Check first person has required fields
		if len(people) > 0 {
			first := people[0].(map[string]any)
			if _, ok := first["id"]; !ok {
				t.Error("expected id in person")
			}
			if _, ok := first["name"]; !ok {
				t.Error("expected name in person")
			}
		}
	})
}

func TestPeoplePingable(t *testing.T) {
	h := harness.New(t)

	t.Run("list pingable people", func(t *testing.T) {
		result := h.Run("people-pingable")

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		// Should have count and people array
		if _, ok := result.JSON["count"]; !ok {
			t.Error("expected count in response")
		}
		if _, ok := result.JSON["people"]; !ok {
			t.Error("expected people array in response")
		}
	})
}

func TestPeopleProject(t *testing.T) {
	h := harness.New(t)

	t.Run("list people on project", func(t *testing.T) {
		result := h.Run("people-project", h.ProjectID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("project_id") == 0 {
			t.Error("expected project_id in response")
		}

		if _, ok := result.JSON["people"]; !ok {
			t.Error("expected people array in response")
		}
	})
}

func TestMyProfile(t *testing.T) {
	h := harness.New(t)

	t.Run("view my profile", func(t *testing.T) {
		result := h.Run("my-profile")

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
		if result.GetString("email") == "" {
			t.Error("expected email in response")
		}
	})
}

func TestPerson(t *testing.T) {
	h := harness.New(t)

	// First get a person ID from the people list
	peopleResult := h.Run("people")
	if !peopleResult.Success() {
		t.Fatalf("failed to get people: %s", peopleResult.Stderr)
	}

	people, ok := peopleResult.JSON["people"].([]any)
	if !ok || len(people) == 0 {
		t.Skip("no people found")
	}

	firstPerson := people[0].(map[string]any)
	personID := fmt.Sprintf("%.0f", firstPerson["id"].(float64))

	t.Run("view person details", func(t *testing.T) {
		result := h.Run("person", personID)

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

	t.Run("missing person_id", func(t *testing.T) {
		result := h.Run("person")

		if result.Success() {
			t.Error("expected failure without person_id")
		}
	})
}
