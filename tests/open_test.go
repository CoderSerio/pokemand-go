package integration_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenCreatesMissingManagedScript(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	dataDir := filepath.Join(tempDir, "data")

	output := runCLIWithInput(t, configDir, dataDir, "y\n", "open", "team/deploy", "--editor", "true")
	assertNoChineseOutput(t, output)
	assertContains(t, output, "Managed skill script does not exist. Create")

	targetPath := filepath.Join(dataDir, "scripts", "team", "deploy.sh")
	content, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("expected created script to exist: %v", err)
	}
	assertContains(t, string(content), "#!/usr/bin/env sh")
}

func TestOpenCreatesJavaScriptTemplateWhenExtensionIsProvided(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	dataDir := filepath.Join(tempDir, "data")

	output := runCLIWithInput(t, configDir, dataDir, "y\n", "open", "tooling/task.js", "--editor", "true")
	assertNoChineseOutput(t, output)
	assertContains(t, output, "Managed skill script does not exist. Create")

	targetPath := filepath.Join(dataDir, "scripts", "tooling", "task.js")
	content, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("expected created script to exist: %v", err)
	}
	assertContains(t, string(content), "#!/usr/bin/env node")
	assertContains(t, string(content), "// task")
}

func TestOpenCanBeCanceled(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	dataDir := filepath.Join(tempDir, "data")

	output := runCLIWithInput(t, configDir, dataDir, "n\n", "open", "cleanup.sh", "--editor", "true")
	assertNoChineseOutput(t, output)
	assertContains(t, output, "Operation canceled.")

	if _, err := os.Stat(filepath.Join(dataDir, "scripts", "cleanup.sh")); !os.IsNotExist(err) {
		t.Fatalf("expected canceled open not to create a script, got err=%v", err)
	}
}
