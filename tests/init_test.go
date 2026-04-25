package integration_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitIsIdempotentAndEnglish(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	dataDir := filepath.Join(tempDir, "data")

	firstOutput := runCLI(t, configDir, dataDir, "init")
	assertNoChineseOutput(t, firstOutput)
	assertContains(t, firstOutput, "pkmg workspace initialized.")
	assertContains(t, firstOutput, "Config: "+filepath.Join(configDir, "meta.json"))
	assertContains(t, firstOutput, "Data: "+dataDir)
	assertContains(t, firstOutput, "Scripts: "+filepath.Join(dataDir, "scripts"))

	secondOutput := runCLI(t, configDir, dataDir, "init")
	assertNoChineseOutput(t, secondOutput)
	assertContains(t, secondOutput, "pkmg is already initialized.")
	assertContains(t, secondOutput, "Config: "+filepath.Join(configDir, "meta.json"))
	assertContains(t, secondOutput, "Data: "+dataDir)
	assertContains(t, secondOutput, "Scripts: "+filepath.Join(dataDir, "scripts"))

	if _, err := os.Stat(filepath.Join(dataDir, "scripts")); err != nil {
		t.Fatalf("expected scripts directory to exist: %v", err)
	}
}
