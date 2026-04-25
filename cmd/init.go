package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the pkmg workspace",
	Long: `Initialize the pkmg config and data directories.

Example:
  pkmg init
`,
	Run: func(cmd *cobra.Command, args []string) {
		configDir := GetConfigDir()
		dataDir := GetDataDir()
		configPath := filepath.Join(configDir, "meta.json")
		metaPath := filepath.Join(dataDir, "meta.json")
		scriptsDir := filepath.Join(dataDir, "scripts")

		createdAnything, err := ensureDir(configDir)
		if err != nil {
			fmt.Printf("Failed to prepare config directory: %v\n", err)
			return
		}

		created, err := ensureDir(dataDir)
		if err != nil {
			fmt.Printf("Failed to prepare data directory: %v\n", err)
			return
		}
		createdAnything = createdAnything || created

		defaultConfig := map[string]interface{}{
			"dataPath":    dataDir,
			"initialized": true,
		}

		if _, err := os.Stat(configPath); err == nil {
		} else if os.IsNotExist(err) {
			configData, err := json.MarshalIndent(defaultConfig, "", "  ")
			if err != nil {
				fmt.Printf("Failed to build config content: %v\n", err)
				return
			}

			if err := os.WriteFile(configPath, configData, 0644); err != nil {
				fmt.Printf("Failed to write config file: %v\n", err)
				return
			}
			createdAnything = true
		} else {
			fmt.Printf("Failed to inspect config file: %v\n", err)
			return
		}

		if _, err := os.Stat(metaPath); os.IsNotExist(err) {
			configData, _ := json.MarshalIndent(defaultConfig, "", "  ")
			if err := os.WriteFile(metaPath, configData, 0644); err != nil {
				fmt.Printf("Failed to write data metadata: %v\n", err)
				return
			}
		}

		created, err = ensureDir(scriptsDir)
		if err != nil {
			fmt.Printf("Failed to prepare scripts directory: %v\n", err)
			return
		}
		createdAnything = createdAnything || created

		if createdAnything {
			fmt.Println("pkmg workspace initialized.")
		} else {
			fmt.Println("pkmg is already initialized.")
		}

		fmt.Printf("Config: %s\n", configPath)
		fmt.Printf("Data: %s\n", dataDir)
		fmt.Printf("Scripts: %s\n", scriptsDir)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func ensureDir(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return false, nil
	} else if !os.IsNotExist(err) {
		return false, err
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return false, err
	}
	return true, nil
}
