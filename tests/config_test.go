package integration_test

import (
	"path/filepath"
	"testing"
)

func TestConfigListSetDeleteAndInvalidAction(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	dataDir := filepath.Join(tempDir, "data")

	initOutput := runCLI(t, configDir, dataDir, "init")
	assertNoChineseOutput(t, initOutput)

	listOutput := runCLI(t, configDir, dataDir, "config", "list")
	assertNoChineseOutput(t, listOutput)
	assertContains(t, listOutput, "\"initialized\": true")

	setOutput := runCLI(t, configDir, dataDir, "config", "set", "theme", "dark")
	assertNoChineseOutput(t, setOutput)
	assertContains(t, setOutput, `Set config key "theme".`)

	listAfterSet := runCLI(t, configDir, dataDir, "config", "list")
	assertNoChineseOutput(t, listAfterSet)
	assertContains(t, listAfterSet, `"theme": "dark"`)

	deleteOutput := runCLI(t, configDir, dataDir, "config", "del", "theme")
	assertNoChineseOutput(t, deleteOutput)
	assertContains(t, deleteOutput, `Deleted config key "theme".`)

	listAfterDelete := runCLI(t, configDir, dataDir, "config", "list")
	assertNoChineseOutput(t, listAfterDelete)
	assertContains(t, listAfterDelete, `"theme": ""`)

	invalidOutput := runCLI(t, configDir, dataDir, "config", "wat")
	assertNoChineseOutput(t, invalidOutput)
	assertContains(t, invalidOutput, "Unknown action.")
}
