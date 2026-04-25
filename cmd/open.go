package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var (
	editor string
)

var openCmd = &cobra.Command{
	Use:   "open [path]",
	Short: "Open a skill script for editing",
	Long:  "Open an existing skill script or create it in the managed scripts directory before editing",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		filePath, err := findFile(path)
		if err != nil {
			relativePath, internalPath, err := resolveManagedSkillPath(path)
			if err != nil {
				fmt.Printf("Failed to resolve managed skill path: %v\n", err)
				return
			}

			if askForConfirmation(fmt.Sprintf("Managed skill script does not exist. Create %s?", internalPath)) {
				created, err := createLocalSkill(relativePath, "")
				if err != nil {
					fmt.Printf("Failed to create managed skill script: %v\n", err)
					return
				}
				filePath = created.AbsolutePath
			} else {
				fmt.Println("Operation canceled.")
				return
			}
		}

		if err := openFile(filePath, editor); err != nil {
			fmt.Printf("Failed to open managed skill script: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
	openCmd.Flags().StringVarP(&editor, "editor", "e", "vim", "Editor command to use")
}

func findFile(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		internalFilePath := filepath.Join(GetDataDir(), "scripts", path)
		if _, err := os.Stat(internalFilePath); os.IsNotExist(err) {
			return "", err
		}
		path = internalFilePath
	}

	return path, nil
}

func resolveManagedSkillPath(path string) (string, string, error) {
	relativePath, err := normalizeSkillRelativePath(path)
	if err != nil {
		return "", "", err
	}

	scriptsDir, err := ensureScriptsDir()
	if err != nil {
		return "", "", fmt.Errorf("failed to prepare scripts directory: %w", err)
	}

	absolutePath := filepath.Join(scriptsDir, filepath.FromSlash(relativePath))
	return relativePath, absolutePath, nil
}

func getSystemOpenCommand() (string, []string) {
	switch runtime.GOOS {
	case "darwin":
		return "open", []string{}
	case "windows":
		return "cmd", []string{"/c", "start"}
	default: // Linux 等
		return "xdg-open", []string{}
	}
}

func openFile(filePath string, editor string) error {
	if editor != "" {
		cmd := exec.Command(editor, filePath)
		// 对于交互式编辑器，需要将标准输入输出连接到当前进程
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err == nil {
			return nil
		}
		fmt.Printf("Failed to open with editor %q: %v\n", editor, err)
	}

	openCmd, args := getSystemOpenCommand()
	args = append(args, filePath)
	cmd := exec.Command(openCmd, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to open file: %v (%s)", err, string(output))
	}
	return nil
}

func askForConfirmation(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [Y/n]: ", message)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "" || response == "y" || response == "yes"
}
