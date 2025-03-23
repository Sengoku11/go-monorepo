package bootstrap_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Sengoku11/go-monorepo/pkg/bootstrap"
	"github.com/Sengoku11/go-monorepo/pkg/environment"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	ctx, cancel, log := bootstrap.Default()
	defer cancel(nil)

	if ctx == nil {
		t.Errorf("expected context to be non nil")
	}

	if cancel == nil {
		t.Errorf("expected cancel func to be non nil")
	}

	if log == nil {
		t.Errorf("expected logger to be non nil")
	}

	if ctx != nil && ctx.Err() != nil {
		t.Errorf("expected non-canceled context, but this was canceled: %v", ctx.Err())
	}
}

func TestDefault_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("did not panic")
		}
	}()

	tempDir := t.TempDir()
	goWorkFile := filepath.Join(tempDir, "go.work")

	t.Chdir(tempDir)
	t.Setenv("ENVIRONMENT", string(environment.Local))

	if err := os.WriteFile(goWorkFile, []byte("test"), 0o600); err != nil {
		t.Fatalf("failed to create go.work file: %v", err)
	}

	bootstrap.Default()
}
