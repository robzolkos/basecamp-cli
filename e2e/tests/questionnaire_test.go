package tests

import (
	"fmt"
	"testing"

	"github.com/rzolkos/basecamp-cli/e2e/harness"
)

func TestQuestionnaire(t *testing.T) {
	h := harness.New(t)

	t.Run("get questionnaire", func(t *testing.T) {
		result := h.Run("questionnaire", h.ProjectID)

		if !result.Success() {
			// Questionnaire may not exist in project - that's ok
			errMsg := result.ErrorMessage()
			if errMsg == "no questionnaire found in this project" {
				t.Skip("no questionnaire in test project")
			}
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("project_id") == 0 {
			t.Error("expected project_id in response")
		}
		if result.GetInt("questionnaire_id") == 0 {
			t.Error("expected questionnaire_id in response")
		}
	})
}

func TestQuestions(t *testing.T) {
	h := harness.New(t)

	t.Run("list questions", func(t *testing.T) {
		result := h.Run("questions", h.ProjectID)

		if !result.Success() {
			// Questionnaire may not exist in project - that's ok
			errMsg := result.ErrorMessage()
			if errMsg == "no questionnaire found in this project" {
				t.Skip("no questionnaire in test project")
			}
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("project_id") == 0 {
			t.Error("expected project_id in response")
		}
		if result.GetInt("questionnaire_id") == 0 {
			t.Error("expected questionnaire_id in response")
		}

		// Questions array should exist
		if _, ok := result.JSON["questions"]; !ok {
			t.Error("expected questions array in response")
		}
	})
}

func TestQuestion(t *testing.T) {
	h := harness.New(t)

	// First get questions list
	questionsResult := h.Run("questions", h.ProjectID)
	if !questionsResult.Success() {
		errMsg := questionsResult.ErrorMessage()
		if errMsg == "no questionnaire found in this project" {
			t.Skip("no questionnaire in test project")
		}
		t.Fatalf("failed to get questions: %s", questionsResult.Stderr)
	}

	questions, ok := questionsResult.JSON["questions"].([]any)
	if !ok || len(questions) == 0 {
		t.Skip("no questions found")
	}

	firstQuestion := questions[0].(map[string]any)
	questionID := fmt.Sprintf("%.0f", firstQuestion["id"].(float64))

	t.Run("view question", func(t *testing.T) {
		result := h.Run("question", h.ProjectID, questionID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("id") == 0 {
			t.Error("expected id in response")
		}
		if result.GetString("title") == "" {
			t.Error("expected title in response")
		}
	})

	t.Run("view question with comments", func(t *testing.T) {
		result := h.Run("question", h.ProjectID, questionID, "--comments")

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		// Should have comments field (even if empty)
		if result.JSON != nil {
			// Comments may or may not be present based on whether there are any
			if result.GetInt("id") == 0 {
				t.Error("expected id in response")
			}
		}
	})

	t.Run("missing question_id", func(t *testing.T) {
		result := h.Run("question", h.ProjectID)

		if result.Success() {
			t.Error("expected failure without question_id")
		}
	})
}

func TestQuestionAnswers(t *testing.T) {
	h := harness.New(t)

	// First get questions list
	questionsResult := h.Run("questions", h.ProjectID)
	if !questionsResult.Success() {
		errMsg := questionsResult.ErrorMessage()
		if errMsg == "no questionnaire found in this project" {
			t.Skip("no questionnaire in test project")
		}
		t.Fatalf("failed to get questions: %s", questionsResult.Stderr)
	}

	questions, ok := questionsResult.JSON["questions"].([]any)
	if !ok || len(questions) == 0 {
		t.Skip("no questions found")
	}

	firstQuestion := questions[0].(map[string]any)
	questionID := fmt.Sprintf("%.0f", firstQuestion["id"].(float64))

	t.Run("list answers", func(t *testing.T) {
		result := h.Run("question-answers", h.ProjectID, questionID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("question_id") == 0 {
			t.Error("expected question_id in response")
		}

		// Answers array should exist
		if _, ok := result.JSON["answers"]; !ok {
			t.Error("expected answers array in response")
		}
	})

	t.Run("missing question_id", func(t *testing.T) {
		result := h.Run("question-answers", h.ProjectID)

		if result.Success() {
			t.Error("expected failure without question_id")
		}
	})
}

func TestQuestionAnswer(t *testing.T) {
	h := harness.New(t)

	// First get questions list
	questionsResult := h.Run("questions", h.ProjectID)
	if !questionsResult.Success() {
		errMsg := questionsResult.ErrorMessage()
		if errMsg == "no questionnaire found in this project" {
			t.Skip("no questionnaire in test project")
		}
		t.Fatalf("failed to get questions: %s", questionsResult.Stderr)
	}

	questions, ok := questionsResult.JSON["questions"].([]any)
	if !ok || len(questions) == 0 {
		t.Skip("no questions found")
	}

	firstQuestion := questions[0].(map[string]any)
	questionID := fmt.Sprintf("%.0f", firstQuestion["id"].(float64))

	// Get answers
	answersResult := h.Run("question-answers", h.ProjectID, questionID)
	if !answersResult.Success() {
		t.Fatalf("failed to get answers: %s", answersResult.Stderr)
	}

	answers, ok := answersResult.JSON["answers"].([]any)
	if !ok || len(answers) == 0 {
		t.Skip("no answers found")
	}

	firstAnswer := answers[0].(map[string]any)
	answerID := fmt.Sprintf("%.0f", firstAnswer["id"].(float64))

	t.Run("view answer", func(t *testing.T) {
		result := h.Run("question-answer", h.ProjectID, answerID)

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

	t.Run("view answer with comments", func(t *testing.T) {
		result := h.Run("question-answer", h.ProjectID, answerID, "--comments")

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON != nil {
			if result.GetInt("id") == 0 {
				t.Error("expected id in response")
			}
		}
	})

	t.Run("missing answer_id", func(t *testing.T) {
		result := h.Run("question-answer", h.ProjectID)

		if result.Success() {
			t.Error("expected failure without answer_id")
		}
	})
}
