package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindProjectID(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "sub", "folder")
	os.MkdirAll(subDir, 0755)

	// Create .basecamp.yml in tmpDir
	configContent := "project_id: 12345678\n"
	os.WriteFile(filepath.Join(tmpDir, ".basecamp.yml"), []byte(configContent), 0644)

	// Change to subDir
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(subDir)

	// Should find config in parent
	projectID, err := FindProjectID()
	if err != nil {
		t.Fatalf("FindProjectID() error = %v", err)
	}
	if projectID != "12345678" {
		t.Errorf("FindProjectID() = %v, want %v", projectID, "12345678")
	}
}

func TestFindProjectIDInCurrentDir(t *testing.T) {
	tmpDir := t.TempDir()

	configContent := "project_id: 99999999\n"
	os.WriteFile(filepath.Join(tmpDir, ".basecamp.yml"), []byte(configContent), 0644)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	projectID, err := FindProjectID()
	if err != nil {
		t.Fatalf("FindProjectID() error = %v", err)
	}
	if projectID != "99999999" {
		t.Errorf("FindProjectID() = %v, want %v", projectID, "99999999")
	}
}

func TestFindProjectIDNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	projectID, err := FindProjectID()
	if err != nil {
		t.Fatalf("FindProjectID() error = %v", err)
	}
	if projectID != "" {
		t.Errorf("FindProjectID() = %v, want empty string", projectID)
	}
}

func TestReadProjectIDFormats(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{"simple", "project_id: 123", "123"},
		{"with quotes", `project_id: "456"`, "456"},
		{"with single quotes", "project_id: '789'", "789"},
		{"with spaces", "project_id:   999  ", "999"},
		{"with comment", "# comment\nproject_id: 111", "111"},
		{"empty lines", "\n\nproject_id: 222\n\n", "222"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			path := filepath.Join(tmpDir, ".basecamp.yml")
			os.WriteFile(path, []byte(tt.content), 0644)

			got, err := readProjectID(path)
			if err != nil {
				t.Fatalf("readProjectID() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("readProjectID() = %v, want %v", got, tt.want)
			}
		})
	}
}
