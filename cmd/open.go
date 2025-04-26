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
	Short: "打开文件",
	Long:  "打开文件进行编辑，可以指定编辑工具",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		filePath, err := findFile(path)
		if err != nil {
			// 文件不存在，让用户确认是否创建文件
			internalPath := filepath.Join(GetDataDir(), "scripts", path)
			if askForConfirmation(fmt.Sprintf("文件不存在，是否创建 %s？", internalPath)) {
				created, err := createNewFile(internalPath)
				if err != nil {
					fmt.Printf("创建文件失败: %v\n", err)
					return
				}
				if created {
					// fmt.Printf("创建成功: %s\n", internalPath)
					filePath = internalPath
				}
			} else {
				fmt.Println("操作已取消")
				return
			}
		}

		// 打开文件进行编辑
		if err := openFile(filePath, editor); err != nil {
			fmt.Printf("打开文件失败: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
	openCmd.Flags().StringVarP(&editor, "editor", "e", "vim", "指定编辑工具")
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

func createNewFile(path string) (bool, error) {
	// 创建文件并写入初始内容
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0755)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// 写入 shebang 和初始注释
	_, err = file.WriteString("#!/bin/sh\n")
	if err != nil {
		os.Remove(path) // 如果写入失败，删除文件
		return false, fmt.Errorf("写入文件失败: %v", err)
	}

	return true, nil
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
		fmt.Printf("无法使用 %s 打开文件: %v\n", editor, err)
	}

	openCmd, args := getSystemOpenCommand()
	args = append(args, filePath)
	cmd := exec.Command(openCmd, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("无法打开文件: %v (%s)", err, string(output))
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
