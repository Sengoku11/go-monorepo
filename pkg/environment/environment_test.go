package environment_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Sengoku11/go-monorepo/pkg/environment"
)

const (
	envKey = "TEST_ENV"
	envVal = "123"
)

func TestCurrent(t *testing.T) {
	tests := []struct {
		name      string
		env       string
		expectEnv environment.Environment
		expectErr error
	}{
		{
			name:      "uppercase",
			env:       "LOCAL",
			expectEnv: environment.Local,
			expectErr: nil,
		},
		{
			name:      "camel case",
			env:       "PrOd",
			expectEnv: environment.Prod,
			expectErr: nil,
		},
		{
			name:      "invalid value",
			env:       "random",
			expectEnv: "",
			expectErr: environment.ErrUnknownEnvironment,
		},
		{
			name:      "caps",
			env:       "",
			expectEnv: "",
			expectErr: environment.ErrEnvironmentNotSet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env != "" {
				t.Setenv("ENVIRONMENT", tt.env)
			}

			env, err := environment.Current()
			if !errors.Is(err, tt.expectErr) {
				t.Errorf("expected %v error, but got %v", tt.expectErr, err)
			}

			if env != tt.expectEnv {
				t.Errorf("expected %s env, but got %s", tt.expectEnv, env)
			}
		})
	}
}

func TestLoadLocalDotEnv(t *testing.T) {
	tempDir := t.TempDir()
	goWorkFile := filepath.Join(tempDir, "go.work")
	envFile := filepath.Join(tempDir, ".env")
	envContent := fmt.Sprintf("%s=%s", envKey, envVal)

	t.Chdir(tempDir)
	t.Setenv("ENVIRONMENT", string(environment.Local))

	if err := os.WriteFile(goWorkFile, []byte("test"), 0o600); err != nil {
		t.Fatalf("failed to create go.work file: %v", err)
	}

	if err := os.WriteFile(envFile, []byte(envContent), 0o600); err != nil {
		t.Fatalf("failed to create .env file: %v", err)
	}

	if err := environment.LoadLocalDotEnv(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if env := os.Getenv(envKey); env != envVal {
		t.Errorf("env data is incorrect: %s", envVal)
	}

	if err := os.Unsetenv(envKey); err != nil {
		t.Errorf("cannot unset envKey: %v", err)
	}
}

func TestLoadLocalDotEnv_NonLocalEnvironment(t *testing.T) {
	t.Setenv("ENVIRONMENT", string(environment.Prod))

	if err := environment.LoadLocalDotEnv(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if env := os.Getenv(envKey); env == envVal {
		t.Errorf("env data shouln't exist, ensure you clean envs")
	}
}

//nolint:paralleltest
func TestLoadLocalDotEnv_GoWorkNotFound(t *testing.T) {
	tempDir := t.TempDir()

	t.Chdir(tempDir)

	if err := environment.LoadLocalDotEnv(); err == nil {
		t.Fatalf("expected error, got nil")
	} else if !strings.Contains(err.Error(), "cannot find workspace folder:") {
		t.Errorf("expected error, got %v", err)
	}

	if env := os.Getenv(envKey); env == envVal {
		t.Errorf("env data shouln't exist, ensure you clean envs")
	}
}

//nolint:paralleltest
func TestLoadLocalDotEnv_EnvFileNotFound(t *testing.T) {
	tempDir := t.TempDir()
	goWorkFile := filepath.Join(tempDir, "go.work")

	t.Chdir(tempDir)

	if err := os.WriteFile(goWorkFile, []byte("test"), 0o600); err != nil {
		t.Fatalf("failed to create go.work file: %v", err)
	}

	if err := environment.LoadLocalDotEnv(); err == nil {
		t.Fatalf("expected error, got nil")
	} else if !strings.Contains(err.Error(), "cannot load .env file") {
		t.Errorf("expected error, got %v", err)
	}

	if env := os.Getenv(envKey); env == envVal {
		t.Errorf("env data shouln't exist, ensure you clean envs")
	}
}

//nolint:paralleltest
func TestLoadLocalDotEnv_EnvironmentVarIsMissing(t *testing.T) {
	tempDir := t.TempDir()
	goWorkFile := filepath.Join(tempDir, "go.work")
	envFile := filepath.Join(tempDir, ".env")
	envContent := fmt.Sprintf("%s=%s", envKey, envVal)

	t.Chdir(tempDir)

	if err := os.WriteFile(goWorkFile, []byte("test"), 0o600); err != nil {
		t.Fatalf("failed to create go.work file: %v", err)
	}

	if err := os.WriteFile(envFile, []byte(envContent), 0o600); err != nil {
		t.Fatalf("failed to create .env file: %v", err)
	}

	if err := environment.LoadLocalDotEnv(); !errors.Is(err, environment.ErrEnvironmentNotSet) {
		t.Fatalf("unexpected error: %v", err)
	}
}
