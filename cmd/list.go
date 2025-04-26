package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	search string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出文件",
	Long:  "列出所有文件",
	Run: func(cmd *cobra.Command, args []string) {
		dataDir := GetDataDir()
		scriptsDir := filepath.Join(dataDir, "scripts")
		if err := os.MkdirAll(scriptsDir, 0755); err != nil {
			fmt.Printf("创建脚本目录失败: %v\n", err)
			return
		}

		files, err := os.ReadDir(scriptsDir)
		if err != nil {
			fmt.Printf("读取脚本目录失败: %v\n", err)
			return
		}

		hasFiles := false
		for _, file := range files {
			if search != "" && !strings.Contains(file.Name(), search) {
				continue
			}
			hasFiles = true
			fmt.Println(file.Name())
		}

		if !hasFiles {
			if search != "" {
				fmt.Printf("没有找到包含 '%s' 的文件\n", search)
			} else {
				fmt.Printf("%s 目录为空\n", scriptsDir)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&search, "search", "s", "", "搜索文件")
}
