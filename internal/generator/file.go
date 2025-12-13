package generator

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/suda-3156/leetcode-cli/internal/config"
)

// GetDefaultTemplate returns the default file path template
func GetDefaultTemplate() string {
	return config.FILEPATH_TMPL
}

// GetOutputPath returns the output file path based on the configuration and frontend ID.
func GetOutputPath(cfg *config.Config, titleSlug, frontendID, langSlug string) (string, error) {
	tmpl, err := template.New("outputPathTemplate").Parse(cfg.OutPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse output path template: %w", err)
	}

	langConfig := config.GetLangConfig(langSlug)

	var buf bytes.Buffer
	data := map[string]string{
		"Date":       cfg.GetCurrentDate(),
		"FrontendID": frontendID,
		"TitleSlug":  titleSlug,
		"Extension":  langConfig.Extension,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute output path template: %w", err)
	}
	return buf.String(), nil
}
