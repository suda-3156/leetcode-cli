package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/suda-3156/leetcode-cli/internal/config"
)

var (
	initPathFlag  string
	initForceFlag bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new config file",
	Long: `Initialize a new LeetCode CLI config file with default settings.

The config file will be created at the specified path (default: ./.leetcode-cli.yaml).
If the file already exists, you will be prompted to confirm overwriting unless --force is specified.

Examples:
  lcli init
  lcli init --path ~/.leetcode-cli.yaml
  lcli init --force
  lcli init --path ./config.yaml --force`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Use default path if not specified
		if initPathFlag == "" {
			initPathFlag = "./.leetcode-cli.yaml"
		}

		err := config.InitConfig(initPathFlag, initForceFlag)
		if err != nil {
			return err
		}

		fmt.Printf("Config file created successfully: %s\n", initPathFlag)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&initPathFlag, "path", "", "Path where the config file will be created (default: ./.leetcode-cli.yaml)")
	initCmd.Flags().BoolVar(&initForceFlag, "force", false, "Force overwrite if config file already exists")
}
