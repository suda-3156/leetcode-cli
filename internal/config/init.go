package config

import (
	"embed"
	"fmt"

	"github.com/suda-3156/leetcode-cli/internal/file"
	prompt "github.com/suda-3156/leetcode-cli/internal/tui/subcommands/init"
)

//go:embed template/*.yaml
var configTemplateFS embed.FS

// GetDefaultConfigTemplate returns the default config template content.
func GetDefaultConfigTemplate() (string, error) {
	content, err := configTemplateFS.ReadFile("template/default.yaml")
	if err != nil {
		return "", fmt.Errorf("failed to read default config template: %w", err)
	}
	return string(content), nil
}

// InitConfig creates a new config file at the specified path.
// If force is true, it will overwrite existing files without prompting.
// If force is false and the file exists, it will prompt the user.
func InitConfig(path string, force bool) error {
	// Check if file exists
	exists := file.FileExists(path)

	if exists && !force {
		// Prompt user for confirmation
		shouldOverwrite, err := prompt.Run(path)
		if err != nil {
			return err
		}
		if !shouldOverwrite {
			return fmt.Errorf("config file already exists: %s", path)
		}
	}

	// Get template content
	template, err := GetDefaultConfigTemplate()
	if err != nil {
		return err
	}

	// Write to file
	err = file.Save(path, template)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	return nil
}
