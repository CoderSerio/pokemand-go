package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	inspectAsJSON    bool
	inspectLineCount int
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [path]",
	Short: "Inspect a managed skill script",
	Long:  "Show metadata, paths, permissions, and content preview for a managed skill script",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		script, err := findManagedScript(args[0], inspectLineCount)
		if err != nil {
			fmt.Printf("Failed to inspect managed skill script: %v\n", err)
			return
		}

		if inspectAsJSON {
			if err := printJSON(script); err != nil {
				fmt.Printf("Failed to print JSON output: %v\n", err)
			}
			return
		}

		fmt.Printf("name: %s\n", script.Name)
		fmt.Printf("path: %s\n", script.RelativePath)
		fmt.Printf("absolutePath: %s\n", script.AbsolutePath)
		fmt.Printf("sizeBytes: %d\n", script.SizeBytes)
		fmt.Printf("mode: %s\n", script.Mode)
		fmt.Printf("executable: %t\n", script.Executable)
		if script.Shebang != "" {
			fmt.Printf("shebang: %s\n", script.Shebang)
		}
		fmt.Printf("modifiedAt: %s\n", script.ModifiedAt.Format("2006-01-02T15:04:05Z07:00"))
		fmt.Println("preview:")
		for _, line := range script.Preview {
			fmt.Printf("  %s\n", line)
		}
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().BoolVar(&inspectAsJSON, "json", false, "Print script details as JSON")
	inspectCmd.Flags().IntVar(&inspectLineCount, "preview-lines", defaultPreviewLines, "Number of preview lines to show")
}
