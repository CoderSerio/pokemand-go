/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
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
var configDir string
var configDirOnce sync.Once

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

func GetConfigDir() string {
	configDirOnce.Do(func() {
		if envPath := os.Getenv("PKMG_CONFIG_DIR"); envPath != "" {
			configDir = envPath
			return
		}

		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			fmt.Println("获取用户配置目录失败:", err)
			return
		}
		configDir = filepath.Join(userConfigDir, "pkmg")
	})

	return configDir
}

func GetDataDir() string {
	dataDirOnce.Do(func() {
		if envPath := os.Getenv("PKMG_DATA_DIR"); envPath != "" {
			dataDir = envPath
			return
		}

		if configuredPath := readConfiguredDataDir(); configuredPath != "" {
			dataDir = configuredPath
			return
		}

		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			fmt.Println("获取用户数据目录失败:", err)
			return
		}
		dataDir = filepath.Join(userConfigDir, "pkmg")
	})

	return dataDir
}

func resetPathCaches() {
	dataDir = ""
	configDir = ""
	dataDirOnce = sync.Once{}
	configDirOnce = sync.Once{}
}

func readConfiguredDataDir() string {
	configPath := filepath.Join(GetConfigDir(), "meta.json")
	content, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}

	var config map[string]interface{}
	if err := json.Unmarshal(content, &config); err != nil {
		return ""
	}

	value, ok := config["dataPath"].(string)
	if !ok || value == "" {
		return ""
	}

	return value
}
