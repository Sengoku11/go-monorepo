package loadenv_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Sengoku11/go-monorepo/pkg/loadenv"
)

const (
	envKey = "TEST_ENV"
	envVal = "123"
)

func TestLocal(t *testing.T) {
	tempDir := t.TempDir()
	goWorkFile := filepath.Join(tempDir, "go.work")
	envFile := filepath.Join(tempDir, ".env")
	envContent := fmt.Sprintf("%s=%s", envKey, envVal)

	t.Chdir(tempDir)
	t.Setenv("ENVIRONMENT", "local")

	if err := os.WriteFile(goWorkFile, []byte("test"), 0o600); err != nil {
		t.Fatalf("failed to create go.work file: %v", err)
	}

	if err := os.WriteFile(envFile, []byte(envContent), 0o600); err != nil {
		t.Fatalf("failed to create .env file: %v", err)
	}

	if err := loadenv.Local(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if env := os.Getenv(envKey); env != envVal {
		t.Errorf("env data is incorrect: %s", envVal)
	}

	if err := os.Unsetenv(envKey); err != nil {
		t.Errorf("cannot unset envKey: %v", err)
	}
}

func TestLocal_NonLocalEnvironment(t *testing.T) {
	t.Setenv("ENVIRONMENT", "production")

	if err := loadenv.Local(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if env := os.Getenv(envKey); env == envVal {
		t.Errorf("env data shouln't exist, ensure you clean envs")
	}
}

//nolint:paralleltest
func TestLocal_GoWorkNotFound(t *testing.T) {
	tempDir := t.TempDir()

	t.Chdir(tempDir)

	if err := loadenv.Local(); err == nil {
		t.Fatalf("expected error, got nil")
	} else if !strings.Contains(err.Error(), "cannot find workspace folder:") {
		t.Errorf("expected error, got %v", err)
	}

	if env := os.Getenv(envKey); env == envVal {
		t.Errorf("env data shouln't exist, ensure you clean envs")
	}
}

//nolint:paralleltest
func TestLocal_EnvFileNotFound(t *testing.T) {
	tempDir := t.TempDir()
	goWorkFile := filepath.Join(tempDir, "go.work")

	t.Chdir(tempDir)

	if err := os.WriteFile(goWorkFile, []byte("test"), 0o600); err != nil {
		t.Fatalf("failed to create go.work file: %v", err)
	}

	if err := loadenv.Local(); err == nil {
		t.Fatalf("expected error, got nil")
	} else if !strings.Contains(err.Error(), "cannot load .env file") {
		t.Errorf("expected error, got %v", err)
	}

	if env := os.Getenv(envKey); env == envVal {
		t.Errorf("env data shouln't exist, ensure you clean envs")
	}
}
