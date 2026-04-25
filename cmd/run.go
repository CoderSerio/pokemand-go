package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [filePath]",
	Short: "Run a managed skill script",
	Long:  "Run a managed skill script with optional arguments",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: 未来支持一下并行/并发执行多个命令

		script, err := findManagedScript(args[0], 1)
		if err != nil {
			fmt.Printf("Managed skill script not found: %s\n", args[0])
			return
		}

		runCmd, err := buildRunCommand(script, args[1:])
		if err != nil {
			fmt.Printf("Script execution failed: %v\n", err)
			return
		}

		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr

		if err := runCmd.Run(); err != nil {
			fmt.Printf("Script execution failed: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func buildRunCommand(script ScriptInfo, args []string) (*exec.Cmd, error) {
	if command, commandArgs, ok := commandFromShebang(script.Shebang); ok {
		return exec.Command(command, append(commandArgs, append([]string{script.AbsolutePath}, args...)...)...), nil
	}

	command, commandArgs, err := commandFromExtension(script.AbsolutePath)
	if err != nil {
		return nil, err
	}

	return exec.Command(command, append(commandArgs, append([]string{script.AbsolutePath}, args...)...)...), nil
}

func commandFromShebang(shebang string) (string, []string, bool) {
	trimmed := strings.TrimSpace(strings.TrimPrefix(shebang, "#!"))
	if trimmed == "" {
		return "", nil, false
	}

	fields := strings.Fields(trimmed)
	if len(fields) == 0 {
		return "", nil, false
	}

	command := fields[0]
	commandArgs := fields[1:]
	if filepath.Base(command) == "env" && len(commandArgs) > 0 {
		command = commandArgs[0]
		commandArgs = commandArgs[1:]
	}

	if _, err := exec.LookPath(command); err != nil {
		return "", nil, false
	}

	return command, commandArgs, true
}

func commandFromExtension(path string) (string, []string, error) {
	extension := strings.ToLower(filepath.Ext(path))

	type runtimeSpec struct {
		command string
		args    []string
	}

	runtimes := map[string]runtimeSpec{
		"":      {command: "sh"},
		".sh":   {command: "sh"},
		".bash": {command: "bash"},
		".zsh":  {command: "zsh"},
		".js":   {command: "node"},
		".cjs":  {command: "node"},
		".mjs":  {command: "node"},
		".py":   {command: "python3"},
		".rb":   {command: "ruby"},
		".pl":   {command: "perl"},
		".php":  {command: "php"},
		".lua":  {command: "lua"},
		".ps1":  {command: "pwsh", args: []string{"-File"}},
	}

	runtime, ok := runtimes[extension]
	if !ok {
		return "", nil, fmt.Errorf("unsupported script type %q without a shebang", extension)
	}
	if _, err := exec.LookPath(runtime.command); err != nil {
		return "", nil, fmt.Errorf("required runtime %q is not available", runtime.command)
	}

	return runtime.command, runtime.args, nil
}
