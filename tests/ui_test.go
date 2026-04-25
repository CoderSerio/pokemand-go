package integration_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestUISmokeStartAndShutdown(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	dataDir := filepath.Join(tempDir, "data")

	binaryPath := buildCLI(t)
	cmd := exec.Command(binaryPath, "ui", "--no-open", "--port", "0")
	cmd.Env = append(os.Environ(),
		"PKMG_CONFIG_DIR="+configDir,
		"PKMG_DATA_DIR="+dataDir,
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start ui command: %v", err)
	}

	deadline := time.Now().Add(5 * time.Second)
	for {
		if strings.Contains(output.String(), "pkmg ui is running at http://") {
			break
		}
		if time.Now().After(deadline) {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			t.Fatalf("timed out waiting for ui startup output:\n%s", output.String())
		}
		time.Sleep(50 * time.Millisecond)
	}

	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGINT); err != nil {
		t.Fatalf("failed to interrupt ui command: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		t.Fatalf("ui command exited with error: %v\n%s", err, output.String())
	}

	text := output.String()
	assertNoChineseOutput(t, text)
	assertContains(t, text, "pkmg ui is running at http://")
	assertContains(t, text, "Received signal interrupt, shutting down UI...")
}
