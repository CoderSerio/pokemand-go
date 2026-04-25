package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	searchAsJSON bool
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "搜索脚本",
	Long:  "按名称或相对路径搜索已管理脚本",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		scripts, err := listManagedScripts(0)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		results := filterScriptsBySearch(scripts, args[0])
		if searchAsJSON {
			payload := struct {
				Query   string       `json:"query"`
				Count   int          `json:"count"`
				Scripts []ScriptInfo `json:"scripts"`
			}{
				Query:   args[0],
				Count:   len(results),
				Scripts: results,
			}
			if err := printJSON(payload); err != nil {
				fmt.Printf("输出 JSON 失败: %v\n", err)
			}
			return
		}

		if len(results) == 0 {
			fmt.Printf("没有找到包含 '%s' 的文件\n", args[0])
			return
		}

		for _, script := range results {
			fmt.Println(script.RelativePath)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolVar(&searchAsJSON, "json", false, "以 JSON 输出搜索结果")
}
