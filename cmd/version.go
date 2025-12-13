package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var version = "devel"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of lcli(leetcode-cli)",
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Printf("lcli version %s\n", version)
		return nil
	},
}

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	mainVersion := info.Main.Version
	if mainVersion != "" && mainVersion != "(devel)" {
		version = mainVersion
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		version = info.Main.Version
	}
	rootCmd.Version = version

	rootCmd.AddCommand(versionCmd)
}
