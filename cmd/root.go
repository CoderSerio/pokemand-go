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
	Short: "A local-first CLI for managing reusable skill scripts",
	Long: `pkmg is a local-first CLI for managing reusable skill scripts.

Examples:
  pkmg init
  pkmg open cleanup.sh
`,
	Version: Version,
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
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
			fmt.Println("Failed to resolve user config directory:", err)
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
			fmt.Println("Failed to resolve user data directory:", err)
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
