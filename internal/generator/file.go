package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/suda-3156/leetcode-cli/internal/api"
	"github.com/suda-3156/leetcode-cli/internal/config"
)

func GenerateFile(outputPath string, question *api.QuestionDetail, snippet *api.CodeSnippet) error {
	// Create directory
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate file content
	content := GenerateFileContent(
		config.GetCurrentDate(),
		question.QuestionFrontendID,
		question.Title,
		snippet.Lang,
		snippet.LangSlug,
		snippet.Code,
	)

	// Write to file
	if err := os.WriteFile(outputPath, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// GetOutputPath returns the output path based on the provided path flag and defaults.
func GetOutputPath(pathFlag, frontendID, titleSlug, langSlug string) string {
	langConfig, ok := config.GetLangConfig(langSlug)
	if !ok {
		langConfig = config.LangConfig{Extension: ".txt"}
	}

	if pathFlag == "" || pathFlag == "default" {
		return config.GetDefaultOutputPath(frontendID, titleSlug, langConfig.Extension)
	}

	return pathFlag
}
