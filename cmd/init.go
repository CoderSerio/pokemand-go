package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// 默认的初始配置
var defaultConfig = map[string]interface{}{
	"dataPath":    "data",
	"initialized": true,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化环境",
	Long: `初始化 pkm 的全局配置和数据目录，会创建必要的文件结构。

使用示例:
  pkm init
`,
	Run: func(cmd *cobra.Command, args []string) {
		dataDir := GetDataDir()
		metaPath := filepath.Join(dataDir, "meta.json")
		scriptsDir := filepath.Join(dataDir, "scripts")

		_, err := os.Stat(dataDir)
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

		_, err = os.Stat(metaPath)
		if err == nil {
			fmt.Println("配置文件已存在:", metaPath)
		} else if os.IsNotExist(err) {
			configData, err := json.MarshalIndent(defaultConfig, "", "  ")
			if err != nil {
				fmt.Printf("生成配置数据失败: %v\n", err)
				return
			}

			if err := os.WriteFile(metaPath, configData, 0644); err != nil {
				fmt.Printf("写入配置文件失败: %v\n", err)
				return
			}

			fmt.Println("配置文件创建成功:", metaPath)
		} else {
			fmt.Printf("检查配置文件失败: %v\n", err)
			return
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
