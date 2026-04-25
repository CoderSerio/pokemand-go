package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveManagedSkillPathCreatesScriptsDir(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("PKMG_CONFIG_DIR", filepath.Join(tempDir, "config"))
	t.Setenv("PKMG_DATA_DIR", filepath.Join(tempDir, "data"))
	resetPathCaches()
	defer resetPathCaches()

	relativePath, absolutePath, err := resolveManagedSkillPath("team/deploy")
	if err != nil {
		t.Fatalf("resolveManagedSkillPath failed: %v", err)
	}

	if relativePath != "team/deploy.sh" {
		t.Fatalf("unexpected relative path: %s", relativePath)
	}

	expectedAbsolutePath := filepath.Join(tempDir, "data", "scripts", "team", "deploy.sh")
	if absolutePath != expectedAbsolutePath {
		t.Fatalf("unexpected absolute path: %s", absolutePath)
	}

	if _, err := os.Stat(filepath.Join(tempDir, "data", "scripts")); err != nil {
		t.Fatalf("expected scripts dir to exist: %v", err)
	}
}

func TestCreateLocalSkillWithoutInitCreatesNestedFile(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("PKMG_CONFIG_DIR", filepath.Join(tempDir, "config"))
	t.Setenv("PKMG_DATA_DIR", filepath.Join(tempDir, "data"))
	resetPathCaches()
	defer resetPathCaches()

	created, err := createLocalSkill("nested/cleanup", "")
	if err != nil {
		t.Fatalf("createLocalSkill failed: %v", err)
	}

	if created.RelativePath != "nested/cleanup.sh" {
		t.Fatalf("unexpected relative path: %s", created.RelativePath)
	}

	content, err := os.ReadFile(created.AbsolutePath)
	if err != nil {
		t.Fatalf("expected created file to exist: %v", err)
	}

	if string(content) != "#!/usr/bin/env sh\n# cleanup\n\n" {
		t.Fatalf("unexpected default content: %q", string(content))
	}
}
