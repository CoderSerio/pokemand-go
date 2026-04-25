package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	search     string
	listAsJSON bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List managed skill scripts",
	Long:  "List all managed skill scripts",
	Run: func(cmd *cobra.Command, args []string) {
		scripts, err := listManagedScripts(0)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		results := filterScriptsBySearch(scripts, search)

		if listAsJSON {
			payload := struct {
				Count   int          `json:"count"`
				Search  string       `json:"search,omitempty"`
				Scripts []ScriptInfo `json:"scripts"`
			}{
				Count:   len(results),
				Search:  search,
				Scripts: results,
			}
			if err := printJSON(payload); err != nil {
				fmt.Printf("Failed to print JSON output: %v\n", err)
			}
			return
		}

		if len(results) == 0 {
			if search != "" {
				fmt.Printf("No managed skill scripts matched %q.\n", search)
			} else {
				fmt.Printf("No managed skill scripts found in %s.\n", getScriptsDir())
			}
			return
		}

		for _, script := range results {
			fmt.Println(script.RelativePath)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&search, "search", "s", "", "Filter scripts by name or path")
	listCmd.Flags().BoolVar(&listAsJSON, "json", false, "Print the script list as JSON")
}
