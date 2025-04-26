package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [filePath]",
	Short: "运行命令",
	Long:  "运行命令",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dataDir := GetDataDir()
		scriptsDir := filepath.Join(dataDir, "scripts")

		// TODO: 未来支持一下并行/并发执行多个命令

		targetFilePath := filepath.Join(scriptsDir, args[0])
		if _, err := os.Stat(targetFilePath); os.IsNotExist(err) {
			fmt.Printf("文件不存在: %s\n", targetFilePath)
			return
		}

		// fmt.Printf("开始执行脚本: %s\n\n", targetFilePath)
		// startTime := time.Now()

		// TODO: 这里执行可能会有一定的兼容问题，可能需要根据不同的操作系统执行不同的命令
		shellCmd := exec.Command("sh", targetFilePath)
		shellCmd.Stdout = os.Stdout
		shellCmd.Stderr = os.Stderr

		if err := shellCmd.Run(); err != nil {
			fmt.Printf("\n执行失败: %v\n", err)
			return
		}

		// duration := time.Since(startTime)
		// fmt.Printf("\n执行完成 ✅ (耗时: %.2f秒)\n", duration.Seconds())
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
