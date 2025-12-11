package generator

import (
	"bytes"
	"fmt"
	"html"
	"regexp"
	"strings"
	"text/template"

	"github.com/suda-3156/leetcode-cli/internal/config"
)

func GenerateFileContent(date, frontendID, title, langName, langSlug, codeSnippet string) string {
	// Decode HTML entities (e.g., &gt; -> >, &lt; -> <)
	codeSnippet = html.UnescapeString(codeSnippet)

	langConfig := config.GetLangConfig(langSlug)

	replaceData := ReplaceData{
		Date:        date,
		FrontendID:  frontendID,
		Title:       title,
		LangName:    langName,
		CodeSnippet: codeSnippet,
		// Extract function name from code snippet
		FuncName:      GetFunctionName(langConfig, codeSnippet),
		CommentPrefix: langConfig.CommentPrefix,
	}

	content, err := Replace(langConfig, replaceData)
	if err != nil {
		panic(fmt.Sprintf("failed to generate file content: %v", err))
	}

	return content
}

type ReplaceData struct {
	Date          string
	FrontendID    string
	Title         string
	LangName      string
	CodeSnippet   string
	FuncName      string
	CommentPrefix string
}

func Replace(langConfig *config.LangConfig, data ReplaceData) (string, error) {
	tmplContent, err := GetTemplate(langConfig.TemplateFile)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("langTemplate").Parse(tmplContent)
	if err != nil {
		return "", err
	}

	dataMap := map[string]string{
		"Date":          data.Date,
		"FrontendID":    data.FrontendID,
		"Title":         data.Title,
		"LangName":      data.LangName,
		"CodeSnippet":   data.CodeSnippet,
		"FuncName":      data.FuncName,
		"CommentPrefix": data.CommentPrefix,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, dataMap); err != nil {
		return "", err
	}

	return buf.String(), nil
}

const DEFAULT_FUNC_NAME = "FUNC_NAME"

func GetFunctionName(langConfig *config.LangConfig, code string) string {
	pattern := langConfig.FuncDefRegex

	if pattern == "" {
		return DEFAULT_FUNC_NAME
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		panic("failed to compile regex: " + err.Error())
	}

	matches := re.FindStringSubmatch(code)
	if len(matches) >= 2 && matches[1] != "" && matches[1] != "main" {
		return strings.TrimSpace(matches[1])
	}

	return DEFAULT_FUNC_NAME
}
