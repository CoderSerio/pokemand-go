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

func TestRunSupportsJavaScriptAndPythonScripts(t *testing.T) {
	requireRuntime(t, "node")
	requireRuntime(t, "python3")

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	dataDir := filepath.Join(tempDir, "data")

	writeManagedScript(t, dataDir, "hello.js", "#!/usr/bin/env node\nconsole.log(process.argv.slice(2).join('-'))\n")
	jsOutput := runCLI(t, configDir, dataDir, "run", "hello.js", "node", "skill")
	assertNoChineseOutput(t, jsOutput)
	assertContains(t, jsOutput, "node-skill")

	writeManagedScript(t, dataDir, "hello.py", "#!/usr/bin/env python3\nimport sys\nprint('-'.join(sys.argv[1:]))\n")
	pyOutput := runCLI(t, configDir, dataDir, "run", "hello.py", "python", "skill")
	assertNoChineseOutput(t, pyOutput)
	assertContains(t, pyOutput, "python-skill")
}
