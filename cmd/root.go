package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/suda-3156/leetcode-cli/internal/tui"
)

var (
	slugFlag      string
	langFlag      string
	pathFlag      string
	overwriteFlag string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lcli [keyword]",
	Short: "Search LeetCode problems and generate code snippet files",
	Long: `
LeetCode CLI is a tool to search LeetCode problems by keyword
and generate code snippet files for solving them.

Examples:
  lcli "reverse polish"
  lcli --slug evaluate-reverse-polish-notation
  lcli "two sum" --lang golang --path path/to/file.go`,
	Args: func(_ *cobra.Command, args []string) error {
		// If --slug is specified, no keyword argument is needed
		if slugFlag != "" {
			return nil
		}
		if len(args) < 1 {
			return fmt.Errorf(
				"keyword argument is required when --slug is not specified, " +
					"because this cli is not intended to search all problems, " +
					"just for generating code snippet files",
			)
		}
		return nil
	},
	RunE: func(_ *cobra.Command, args []string) error {
		keyword := ""
		if len(args) > 0 {
			keyword = args[0]
		}

		_, err := tui.Run(keyword, slugFlag, langFlag, pathFlag)
		if err != nil {
			return err
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&slugFlag, "slug", "", "Skip search and use titleSlug directly")
	rootCmd.Flags().StringVarP(&langFlag, "lang", "l", "", "Specify language (langSlug: golang, python3, python, etc.)")
	rootCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Output file path")
	rootCmd.Flags().StringVarP(&overwriteFlag, "overwrite", "o", "", "Overwrite behavior (prompt, always, backup, never)")
}
