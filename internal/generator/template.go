package generator

import (
	"fmt"

	"github.com/suda-3156/leetcode-cli/internal/config"
)

func GenerateFileContent(date, frontendID, title, langName, langSlug, codeSnippet string) string {
	langConfig, ok := config.GetLangConfig(langSlug)
	if !ok {
		// Use // as default comment prefix for unknown languages
		langConfig = config.LangConfig{CommentPrefix: "//"}
	}

	comment := langConfig.CommentPrefix

	return fmt.Sprintf("%s %s\n%s %s. %s\n%s %s\n\n%s\n",
		comment, date,
		comment, frontendID, title,
		comment, langName,
		codeSnippet,
	)
}
