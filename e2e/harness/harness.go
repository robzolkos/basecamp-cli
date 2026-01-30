// Package harness provides a test harness for end-to-end testing of the Basecamp CLI
// by executing the CLI binary and capturing stdout, stderr, and exit codes.
package harness

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
)

// Harness provides methods for executing CLI commands and capturing results.
type Harness struct {
	// BinaryPath is the path to the CLI binary
	BinaryPath string

	// ProjectID is the test project ID
	ProjectID string

	// BoardID is the test board ID
	BoardID string

	// CardID is a test card ID for read operations
	CardID string

	// TodolistID is a test todolist ID
	TodolistID string

	// t is the testing context
	t *testing.T
}

// Result contains the output from a CLI command execution.
type Result struct {
	// Stdout is the standard output
	Stdout string

	// Stderr is the standard error output
	Stderr string

	// ExitCode is the process exit code
	ExitCode int

	// JSON is the parsed JSON from stdout (nil if parsing failed)
	JSON map[string]any

	// JSONArray is the parsed JSON array from stdout (nil if not an array)
	JSONArray []map[string]any

	// ParseError is set if JSON parsing failed
	ParseError error
}

// Config holds test harness configuration from environment variables.
type Config struct {
	BinaryPath string
	ProjectID  string
	BoardID    string
	CardID     string
	TodolistID string
}

// Exit codes
const (
	ExitSuccess = 0
	ExitError   = 1
)

// LoadConfig loads test configuration from environment variables.
func LoadConfig() *Config {
	defaultBinary := "./basecamp"
	if cwd, err := os.Getwd(); err == nil {
		// Try to find binary relative to repo root
		for dir := cwd; dir != "/"; dir = filepath.Dir(dir) {
			candidate := filepath.Join(dir, "basecamp")
			if _, err := os.Stat(candidate); err == nil {
				defaultBinary = candidate
				break
			}
		}
	}

	return &Config{
		BinaryPath: getEnvOrDefault("BASECAMP_TEST_BINARY", defaultBinary),
		ProjectID:  os.Getenv("BASECAMP_TEST_PROJECT_ID"),
		BoardID:    os.Getenv("BASECAMP_TEST_BOARD_ID"),
		CardID:     os.Getenv("BASECAMP_TEST_CARD_ID"),
		TodolistID: os.Getenv("BASECAMP_TEST_TODOLIST_ID"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// New creates a new test harness with configuration from environment variables.
func New(t *testing.T) *Harness {
	t.Helper()

	cfg := LoadConfig()

	if cfg.ProjectID == "" {
		t.Skip("BASECAMP_TEST_PROJECT_ID not set, skipping e2e tests")
	}

	return &Harness{
		BinaryPath: cfg.BinaryPath,
		ProjectID:  cfg.ProjectID,
		BoardID:    cfg.BoardID,
		CardID:     cfg.CardID,
		TodolistID: cfg.TodolistID,
		t:          t,
	}
}

// Run executes a CLI command and returns the result.
func (h *Harness) Run(args ...string) *Result {
	h.t.Helper()

	cmd := exec.Command(h.BinaryPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = os.Environ()

	err := cmd.Run()

	result := &Result{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: 0,
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				result.ExitCode = status.ExitStatus()
			}
		} else {
			result.ExitCode = -1
			result.Stderr = err.Error()
		}
	}

	// Try to parse JSON response
	if result.Stdout != "" {
		// Try as object first
		var obj map[string]any
		if err := json.Unmarshal([]byte(result.Stdout), &obj); err == nil {
			result.JSON = obj
		} else {
			// Try as array
			var arr []map[string]any
			if err := json.Unmarshal([]byte(result.Stdout), &arr); err == nil {
				result.JSONArray = arr
			} else {
				result.ParseError = err
			}
		}
	}

	return result
}

// GetString extracts a string value from the JSON response.
func (r *Result) GetString(key string) string {
	if r.JSON == nil {
		return ""
	}
	v, ok := r.JSON[key].(string)
	if !ok {
		return ""
	}
	return v
}

// GetInt extracts an integer value from the JSON response.
func (r *Result) GetInt(key string) int {
	if r.JSON == nil {
		return 0
	}
	// JSON numbers are float64
	v, ok := r.JSON[key].(float64)
	if !ok {
		return 0
	}
	return int(v)
}

// GetNested extracts a nested value from the JSON response.
func (r *Result) GetNested(keys ...string) any {
	if r.JSON == nil {
		return nil
	}
	var current any = r.JSON
	for _, key := range keys {
		m, ok := current.(map[string]any)
		if !ok {
			return nil
		}
		current = m[key]
	}
	return current
}

// Success returns true if the command exited successfully.
func (r *Result) Success() bool {
	return r.ExitCode == ExitSuccess
}

// HasError returns true if stderr contains an error JSON.
func (r *Result) HasError() bool {
	return r.Stderr != "" && r.ExitCode != ExitSuccess
}

// ErrorMessage extracts the error message from stderr.
func (r *Result) ErrorMessage() string {
	if r.Stderr == "" {
		return ""
	}
	var errObj struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal([]byte(r.Stderr), &errObj); err == nil {
		return errObj.Error
	}
	return r.Stderr
}
