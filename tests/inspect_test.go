package integration_test

import (
	"path/filepath"
	"testing"
)

func TestInspectShowsScriptDetailsAndJSON(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	dataDir := filepath.Join(tempDir, "data")

	writeManagedScript(t, dataDir, "cleanup.sh", "#!/bin/sh\n# Cleanup\necho clean\n")

	output := runCLI(t, configDir, dataDir, "inspect", "cleanup.sh")
	assertNoChineseOutput(t, output)
	assertContains(t, output, "name: cleanup.sh")
	assertContains(t, output, "path: cleanup.sh")
	assertContains(t, output, "preview:")

	jsonOutput := runCLI(t, configDir, dataDir, "inspect", "cleanup.sh", "--json")
	assertNoChineseOutput(t, jsonOutput)
	assertContains(t, jsonOutput, `"name": "cleanup.sh"`)
	assertContains(t, jsonOutput, `"relativePath": "cleanup.sh"`)
}

func TestInspectMissingScriptMessageIsEnglish(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	dataDir := filepath.Join(tempDir, "data")

	output := runCLI(t, configDir, dataDir, "inspect", "missing.sh")
	assertNoChineseOutput(t, output)
	assertContains(t, output, "Failed to inspect managed skill script:")
}
