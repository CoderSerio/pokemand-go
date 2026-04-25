package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config [action] [key] [value]",
	Short: "Manage pkmg configuration",
	Long:  "Manage pkmg configuration, including the configured data path",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			listConfigInfo()
			return
		}

		if len(args) == 1 {
			if args[0] == "list" || args[0] == "ls" {
				listConfigInfo()
				return
			}

			fmt.Println("Unknown action.")
			return
		}

		if len(args) == 2 {
			if args[0] == "del" {
				err := SetConfig(args[1], "")
				if err != nil {
					fmt.Printf("Failed to delete config key: %v\n", err)
					return
				}
				fmt.Printf("Deleted config key %q.\n", args[1])
				return
			}

			fmt.Println("Unknown action.")
			return
		}

		if len(args) == 3 {
			if args[0] == "set" {
				key := args[1]
				value := args[2]

				err := SetConfig(key, value)
				if err != nil {
					fmt.Printf("Failed to update config: %v\n", err)
					return
				}
				fmt.Printf("Set config key %q.\n", key)
				return
			}

			fmt.Println("Unknown action.")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func listConfigInfo() {
	str, err := GetConfig()
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		return
	}
	fmt.Println(str)
}

func getConfigFilePath() string {
	configDir := GetConfigDir()
	return filepath.Join(configDir, "meta.json")
}

func GetConfig() (string, error) {
	configPath := getConfigFilePath()
	metaFileContent, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	str := string(metaFileContent)
	return str, nil
}

func SetConfig(key string, value string) error {
	configPath := getConfigFilePath()
	metaFileContent, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			metaFileContent = []byte("{}")
		} else {
			return err
		}
	}

	var config map[string]interface{}
	if err := json.Unmarshal(metaFileContent, &config); err != nil {
		return err
	}

	// TODO: 如果有多级嵌套?
	config[key] = value
	// 后两个参数：行前缀、缩进
	updatedContent, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	configDir := GetConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	return os.WriteFile(configPath, updatedContent, 0644)
}
