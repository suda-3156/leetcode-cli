package generator

import (
	"bytes"
	"fmt"
	"html"
	"text/template"

	"github.com/suda-3156/leetcode-cli/internal/config"
	"github.com/suda-3156/leetcode-cli/internal/parser"
)

func GenerateFileContent(date, frontendID, title, langName, langSlug, codeSnippet string) (string, error) {
	// Decode HTML entities (e.g., &gt; -> >, &lt; -> <)
	codeSnippet = html.UnescapeString(codeSnippet)

	langConfig := config.GetLangConfig(langSlug)

	p, err := parser.NewParser(langSlug, codeSnippet)
	if err != nil {
		return "", fmt.Errorf("failed to create parser: %w", err)
	}
	defer p.Close()

	typeDef, err := p.ExtractTypeDefinition()
	if err != nil {
		return "", fmt.Errorf("failed to extract type definition: %w", err)
	}

	importStmt, err := p.GenerateImportStatement()
	if err != nil {
		return "", fmt.Errorf("failed to generate import statement: %w", err)
	}

	funcName, err := p.ExtractSolutionFuncName()
	if err != nil {
		return "", fmt.Errorf("failed to extract function name: %w", err)
	}

	replaceData := &ReplaceData{
		Date:          date,
		FrontendID:    frontendID,
		Title:         title,
		LangName:      langName,
		CodeSnippet:   codeSnippet,
		FuncName:      funcName,
		CommentPrefix: langConfig.CommentPrefix,
		TypeDef:       typeDef,
		ImportStmt:    importStmt,
	}

	content, err := Replace(langConfig, replaceData)
	if err != nil {
		return "", fmt.Errorf("failed to replace template placeholders: %w", err)
	}

	return content, nil
}

type ReplaceData struct {
	Date          string
	FrontendID    string
	Title         string
	LangName      string
	CodeSnippet   string
	FuncName      string
	CommentPrefix string
	TypeDef       string
	ImportStmt    string
}

func Replace(langConfig *config.LangConfig, data *ReplaceData) (string, error) {
	tmplContent, err := GetTemplate(langConfig.TemplateFile)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("langTemplate").Option("missingkey=zero").Parse(tmplContent)
	if err != nil {
		return "", err
	}

	dataMap := map[string]string{
		"Date":            data.Date,
		"FrontendID":      data.FrontendID,
		"Title":           data.Title,
		"LangName":        data.LangName,
		"CodeSnippet":     data.CodeSnippet,
		"FuncName":        data.FuncName,
		"CommentPrefix":   data.CommentPrefix,
		"TypeDefinitions": data.TypeDef,
		"ImportStatement": data.ImportStmt,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, dataMap); err != nil {
		return "", err
	}

	return buf.String(), nil
}
