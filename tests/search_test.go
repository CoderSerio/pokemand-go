package integration_test

import (
	"path/filepath"
	"testing"
)

func TestSearchReturnsMatchesAndJSON(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	dataDir := filepath.Join(tempDir, "data")

	writeManagedScript(t, dataDir, "cleanup.sh", "#!/bin/sh\necho clean\n")
	writeManagedScript(t, dataDir, "team/deploy.sh", "#!/bin/sh\necho deploy\n")

	output := runCLI(t, configDir, dataDir, "search", "deploy")
	assertNoChineseOutput(t, output)
	assertContains(t, output, "team/deploy.sh")

	jsonOutput := runCLI(t, configDir, dataDir, "search", "deploy", "--json")
	assertNoChineseOutput(t, jsonOutput)
	assertContains(t, jsonOutput, `"query": "deploy"`)
	assertContains(t, jsonOutput, `"relativePath": "team/deploy.sh"`)
}

func TestSearchNoResultsMessageIsEnglish(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	dataDir := filepath.Join(tempDir, "data")

	output := runCLI(t, configDir, dataDir, "search", "missing")
	assertNoChineseOutput(t, output)
	assertContains(t, output, `No managed skill scripts matched "missing".`)
}
