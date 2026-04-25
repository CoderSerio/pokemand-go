package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// 默认的初始配置
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化环境",
	Long: `初始化 pkm 的全局配置和数据目录，会创建必要的文件结构。

使用示例:
  pkm init
`,
	Run: func(cmd *cobra.Command, args []string) {
		configDir := GetConfigDir()
		dataDir := GetDataDir()
		configPath := filepath.Join(configDir, "meta.json")
		metaPath := filepath.Join(dataDir, "meta.json")
		scriptsDir := filepath.Join(dataDir, "scripts")

		_, err := os.Stat(configDir)
		if err == nil {
			fmt.Println("配置目录已存在:", configDir)
		} else if os.IsNotExist(err) {
			fmt.Println("创建配置目录:", configDir)
			if err := os.MkdirAll(configDir, 0755); err != nil {
				fmt.Printf("创建配置目录失败: %v\n", err)
				return
			}
		} else {
			fmt.Printf("检查配置目录失败: %v\n", err)
			return
		}

		_, err = os.Stat(dataDir)
		if err == nil {
			fmt.Println("数据目录已存在:", dataDir)
		} else if os.IsNotExist(err) {
			fmt.Println("创建数据目录:", dataDir)
			if err := os.MkdirAll(dataDir, 0755); err != nil {
				fmt.Printf("创建数据目录失败: %v\n", err)
				return
			}
		} else {
			fmt.Printf("检查数据目录失败: %v\n", err)
			return
		}

		defaultConfig := map[string]interface{}{
			"dataPath":    dataDir,
			"initialized": true,
		}

		_, err = os.Stat(configPath)
		if err == nil {
			fmt.Println("配置文件已存在:", configPath)
		} else if os.IsNotExist(err) {
			configData, err := json.MarshalIndent(defaultConfig, "", "  ")
			if err != nil {
				fmt.Printf("生成配置数据失败: %v\n", err)
				return
			}

			if err := os.WriteFile(configPath, configData, 0644); err != nil {
				fmt.Printf("写入配置文件失败: %v\n", err)
				return
			}

			fmt.Println("配置文件创建成功:", configPath)
		} else {
			fmt.Printf("检查配置文件失败: %v\n", err)
			return
		}

		// 保留 data 目录下的元数据文件，方便用户理解当前数据目录的用途。
		if _, err := os.Stat(metaPath); os.IsNotExist(err) {
			configData, _ := json.MarshalIndent(defaultConfig, "", "  ")
			_ = os.WriteFile(metaPath, configData, 0644)
		}

		_, err = os.Stat(scriptsDir)
		if err == nil {
			fmt.Println("脚本目录已存在:", scriptsDir)
		} else if os.IsNotExist(err) {
			fmt.Println("创建脚本目录:", scriptsDir)
			if err := os.MkdirAll(scriptsDir, 0755); err != nil {
				fmt.Printf("创建脚本目录失败: %v\n", err)
				return
			}
		}

		fmt.Println("\npkm 初始化完成!")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
