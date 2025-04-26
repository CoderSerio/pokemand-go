/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"
)

var (
	Version = "dev" // 这个变量会在编译时通过 -ldflags 注入
)

var dataDir string
var dataDirOnce sync.Once

var rootCmd = &cobra.Command{
	Use:   "pkmg",
	Short: "一个用于封装和管理自定义命令的 CLI 工具",
	Long: `pkmg 是一个简易的命令行工具，用于封装和管理自定义命令。

使用示例:
  pkmg init          初始化环境
  pkmg open file     打开文件
`,
	Version: Version,
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "启用详细输出模式")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func GetDataDir() string {
	dataDirOnce.Do(func() {
		// 使用当前工作目录
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("获取当前工作目录失败:", err)
			return
		}
		dataDir = filepath.Join(cwd, "data")
	})

	return dataDir
}
