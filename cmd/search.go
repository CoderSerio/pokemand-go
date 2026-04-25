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
	Short: "Search managed skill scripts",
	Long:  "Search managed skill scripts by name or relative path",
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
				fmt.Printf("Failed to print JSON output: %v\n", err)
			}
			return
		}

		if len(results) == 0 {
			fmt.Printf("No managed skill scripts matched %q.\n", args[0])
			return
		}

		for _, script := range results {
			fmt.Println(script.RelativePath)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolVar(&searchAsJSON, "json", false, "Print search results as JSON")
}
