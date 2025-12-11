package generator

import (
	"bytes"
	"text/template"

	"github.com/suda-3156/leetcode-cli/internal/config"
)

// GetDefaultTemplate returns the default file path template
func GetDefaultTemplate() string {
	return config.FILEPATH_TMPL
}

// GetOutputPath returns the output file path based on the configuration and frontend ID.
func GetOutputPath(cfg *config.Config, titleSlug, frontendID, langSlug string) string {
	tmpl, err := template.New("outputPathTemplate").Parse(cfg.OutPath)
	if err != nil {
		panic("failed to parse output path template: " + err.Error())
	}

	langConfig := config.GetLangConfig(langSlug)

	var buf bytes.Buffer
	data := map[string]string{
		"Date":       config.GetCurrentDate(cfg),
		"FrontendID": frontendID,
		"TitleSlug":  titleSlug,
		"Extension":  langConfig.Extension,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		panic("failed to execute output path template: " + err.Error())
	}
	return buf.String()
}
