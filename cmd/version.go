package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var (
	version string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of leetcode-cli",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("leetcode-cli version %s\n", version)
		return nil
	},
}

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		version = info.Main.Version
	}
	rootCmd.Version = version

	rootCmd.AddCommand(versionCmd)
}
