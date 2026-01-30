package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

const ProjectConfigFile = ".basecamp.yml"

// FindProjectID looks for .basecamp.yml in current directory and parents,
// returning the project_id if found.
func FindProjectID() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		configPath := filepath.Join(dir, ProjectConfigFile)
		if projectID, err := readProjectID(configPath); err == nil {
			return projectID, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root
			break
		}
		dir = parent
	}

	return "", nil
}

func readProjectID(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Look for project_id: value
		if strings.HasPrefix(line, "project_id:") {
			value := strings.TrimPrefix(line, "project_id:")
			value = strings.TrimSpace(value)
			// Remove quotes if present
			value = strings.Trim(value, `"'`)
			if value != "" {
				return value, nil
			}
		}
	}

	return "", os.ErrNotExist
}
