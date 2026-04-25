package integration_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
)

var (
	hanPattern     = regexp.MustCompile(`\p{Han}`)
	cliBinaryPath  string
	cliBinaryErr   error
	cliBinaryBuild sync.Once
)

func runCLI(t *testing.T, configDir string, dataDir string, args ...string) string {
	t.Helper()
	return runCLIWithInput(t, configDir, dataDir, "", args...)
}

func runCLIWithInput(t *testing.T, configDir string, dataDir string, input string, args ...string) string {
	t.Helper()

	cmd := exec.Command("go", append([]string{"run", ".."}, args...)...)
	cmd.Dir = filepath.Join("..", "tests")
	cmd.Env = append(os.Environ(),
		"PKMG_CONFIG_DIR="+configDir,
		"PKMG_DATA_DIR="+dataDir,
	)
	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\n%s", err, string(output))
	}

	return string(output)
}

func buildCLI(t *testing.T) string {
	t.Helper()

	cliBinaryBuild.Do(func() {
		tempDir, err := os.MkdirTemp("", "pkmg-cli-*")
		if err != nil {
			cliBinaryErr = err
			return
		}

		cliBinaryPath = filepath.Join(tempDir, "pkmg-test-bin")
		cmd := exec.Command("go", "build", "-o", cliBinaryPath, "..")
		cmd.Dir = filepath.Join("..", "tests")
		var output bytes.Buffer
		cmd.Stdout = &output
		cmd.Stderr = &output
		if err := cmd.Run(); err != nil {
			cliBinaryErr = err
		}
	})

	if cliBinaryErr != nil {
		t.Fatalf("failed to build test binary: %v", cliBinaryErr)
	}

	return cliBinaryPath
}

func writeManagedScript(t *testing.T, dataDir string, relativePath string, content string) string {
	t.Helper()

	absolutePath := filepath.Join(dataDir, "scripts", filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0755); err != nil {
		t.Fatalf("failed to create script directory: %v", err)
	}
	if err := os.WriteFile(absolutePath, []byte(content), 0755); err != nil {
		t.Fatalf("failed to write managed script: %v", err)
	}
	return absolutePath
}

func requireRuntime(t *testing.T, command string) {
	t.Helper()
	if _, err := exec.LookPath(command); err != nil {
		t.Skipf("skipping because runtime %q is not available", command)
	}
}

func assertNoChineseOutput(t *testing.T, output string) {
	t.Helper()
	if hanPattern.MatchString(output) {
		t.Fatalf("expected English-only output, got:\n%s", output)
	}
}

func assertContains(t *testing.T, output string, needle string) {
	t.Helper()
	if !strings.Contains(output, needle) {
		t.Fatalf("expected output to contain %q, got:\n%s", needle, output)
	}
}

func assertNotContains(t *testing.T, output string, needle string) {
	t.Helper()
	if strings.Contains(output, needle) {
		t.Fatalf("expected output not to contain %q, got:\n%s", needle, output)
	}
}
