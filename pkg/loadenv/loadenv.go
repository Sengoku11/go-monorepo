// Package loadenv helps to load .env file for local development
package loadenv

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var errGoWorkNotFound = errors.New(`go.work file not found`)

// Local attempts to load an env file located at the go.work folder.
// It returns error if the ENVIRONMENT is local and an env file couldn't be loaded.
func Local() error {
	if environment, exists := os.LookupEnv("ENVIRONMENT"); exists {
		if environment != "local" {
			return nil
		}
	}

	workspaceFolder, err := findGoWorkRoot()
	if err != nil {
		return fmt.Errorf("cannot find workspace folder: %w", err)
	}

	path := filepath.Join(workspaceFolder, ".env")

	err = godotenv.Load(path)
	if err != nil {
		return fmt.Errorf("cannot load .env file: %w", err)
	}

	return nil
}

// findGoWorkRoot starts from the current directory and walks up the folder tree
// until it finds a directory containing a "go.work" file.
// It returns the relative path from the original starting directory to the directory containing go.work.
func findGoWorkRoot() (string, error) {
	startDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	current := startDir

	for {
		// Build the candidate path by joining the current directory with "go.work".
		candidate := filepath.Join(current, "go.work")
		if _, err := os.Stat(candidate); err == nil {
			relPath, err := filepath.Rel(startDir, current)
			if err != nil {
				return "", fmt.Errorf("failed to compute relative path: %w", err)
			}

			return relPath, nil
		}

		// Move one directory up.
		parent := filepath.Dir(current)
		// If the parent is the same as current, we've reached the root.
		if parent == current {
			break
		}

		current = parent
	}

	return "", errGoWorkNotFound
}
