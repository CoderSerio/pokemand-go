package integration_test

import (
	"path/filepath"
	"testing"
)

func TestListHandlesEmptyWorkspace(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	dataDir := filepath.Join(tempDir, "data")

	output := runCLI(t, configDir, dataDir, "list")
	assertNoChineseOutput(t, output)
	assertContains(t, output, "No managed skill scripts found in")
}

func TestListReturnsScriptsAndJSON(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	dataDir := filepath.Join(tempDir, "data")

	writeManagedScript(t, dataDir, "team/deploy.sh", "#!/bin/sh\necho deploy\n")

	output := runCLI(t, configDir, dataDir, "list")
	assertNoChineseOutput(t, output)
	assertContains(t, output, "team/deploy.sh")

	jsonOutput := runCLI(t, configDir, dataDir, "list", "--json")
	assertNoChineseOutput(t, jsonOutput)
	assertContains(t, jsonOutput, `"count": 1`)
	assertContains(t, jsonOutput, `"relativePath": "team/deploy.sh"`)
}
