package integration_test

import (
	"path/filepath"
	"testing"
)

func TestRunPassesThroughArguments(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	dataDir := filepath.Join(tempDir, "data")

	writeManagedScript(t, dataDir, "echo-args.sh", "#!/bin/sh\necho \"$1-$2\"\n")

	output := runCLI(t, configDir, dataDir, "run", "echo-args.sh", "hello", "world")
	assertNoChineseOutput(t, output)
	assertContains(t, output, "hello-world")
}

func TestRunMissingScriptAndFailureMessagesAreEnglish(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	dataDir := filepath.Join(tempDir, "data")

	missingOutput := runCLI(t, configDir, dataDir, "run", "missing.sh")
	assertNoChineseOutput(t, missingOutput)
	assertContains(t, missingOutput, "Managed skill script not found: missing.sh")

	writeManagedScript(t, dataDir, "fail.sh", "#!/bin/sh\nexit 2\n")

	failureOutput := runCLI(t, configDir, dataDir, "run", "fail.sh")
	assertNoChineseOutput(t, failureOutput)
	assertContains(t, failureOutput, "Script execution failed:")
}
