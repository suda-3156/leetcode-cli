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

	tmpl := NewTemplateGen(langSlug)

	comment := tmpl.LangConfig.CommentPrefix

	header := fmt.Sprintf("%s %s\n%s %s. %s\n%s %s\n\n",
		comment, date,
		comment, frontendID, title,
		comment, langName,
	)

	content, err := tmpl.Replace(header, codeSnippet)
	if err != nil {
		panic("failed to generate file content: " + err.Error())
	}

	return content
}

type TemplateGen struct {
	LangConfig *config.LangConfig
	Template   string
}

func NewTemplateGen(langSlug string) *TemplateGen {
	langConfig := config.GetLangConfig(langSlug)

	tmpl, ok := templateMap[langSlug]
	if !ok {
		tmpl = templateMap["default"]
	}

	return &TemplateGen{
		LangConfig: langConfig,
		Template:   tmpl,
	}
}

func (t *TemplateGen) Replace(header, codeSnippet string) (string, error) {
	tmpl, err := template.New("langTemplate").Parse(t.Template)
	if err != nil {
		return "", err
	}

	data := struct {
		Header       string
		CodeSnippet  string
		FunctionName string
	}{
		Header:       header,
		CodeSnippet:  codeSnippet,
		FunctionName: t.GetFunctionName(codeSnippet),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

const DEFAULT_FUNC_NAME = "FUNC_NAME"

func (t *TemplateGen) GetFunctionName(code string) string {
	pattern := t.LangConfig.FuncDefRegex

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
