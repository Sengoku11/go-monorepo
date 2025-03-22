// Package environment defines Environment constants and provides helpers
// for loading a .env file during local development.
package environment

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// Environment is a typed string representing the environment name.
type Environment string

// IsValid checks if there is no typo and correct Environment defined.
func (e Environment) IsValid() bool {
	switch e {
	case Local, UAT, Prod:
		return true
	}

	return false
}

// Use these constants to standardize environment naming across the monorepo.
const (
	Local Environment = "local"
	UAT   Environment = "uat"
	Prod  Environment = "prod"
)

var (
	// ErrGoWorkNotFound indicates that a go.work file was not found in the folder tree.
	ErrGoWorkNotFound = errors.New("go.work file not found")

	// ErrEnvironmentNotSet is returned when the ENVIRONMENT variable is not set.
	ErrEnvironmentNotSet = errors.New("env ENVIRONMENT is not set")

	// ErrUnknownEnvironment is returned when the ENVIRONMENT variable's value is not recognized.
	ErrUnknownEnvironment = errors.New("unknown environment")
)

// Current returns Environment specified in ENV.
func Current() (Environment, error) {
	env, exists := os.LookupEnv("ENVIRONMENT")
	if !exists {
		return "", ErrEnvironmentNotSet
	}

	normalized := Environment(strings.ToLower(env))
	if !normalized.IsValid() {
		return "", ErrUnknownEnvironment
	}

	return normalized, nil
}

// LoadLocalDotEnv attempts to load an env file located at the go.work folder.
// It returns error only if the Environment is Local and an env file couldn't be loaded.
func LoadLocalDotEnv() error {
	if env, err := Current(); err == nil && env != Local {
		// No need to load a dot env file.
		return nil
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

	// Control check, after loading .env file the ENVIRONMENT variable must be defined.
	// TODO: parse .env file instead.
	if _, err := Current(); err != nil {
		return err
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

	return "", ErrGoWorkNotFound
}
