package generator

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/suda-3156/leetcode-cli/internal/config"
)

func GenerateFileContent(date, frontendID, title, langName, langSlug, codeSnippet string) string {
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
	LangConfig config.LangConfig
	Template   string
}

func NewTemplateGen(langSlug string) *TemplateGen {
	langConfig := config.GetLangConfig(langSlug)

	var tmpl string
	switch langSlug {
	case "golang":
		tmpl = goTemplate
	case "python", "python3":
		tmpl = pythonTemplate
	default:
		tmpl = defaultTemplate
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
		Header      string
		CodeSnippet string
	}{
		Header:      header,
		CodeSnippet: codeSnippet,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

const goTemplate = `
{{ .Header }}

package main

{{ .CodeSnippet }}

func main() {
	testCases := []struct {
		// Define your test case structure here
		input struct {}
		want  string	
	}{
		// Add your test cases here
	}

	for _, tc := range testCases {
		result := FUNC_NAME(tc.input)
		if result != tc.want {
			fmt.Printf("Test failed for input %v: got %v, want %v\n", tc.input, result, tc.want)
		} else {
			fmt.Printf("Test passed for input %v\n", tc.input)
		}
	}
}
`

const pythonTemplate = `
{{ .Header }}

from typing import List, Any, TypedDict

{{ .CodeSnippet }}

TestCase = TypedDict("TestCase", {"input": Any, "want": Any})


def main():
    test_cases: List[TestCase] = [
		# Add your test cases here
        {"input": {}, "want": ""}
    ]

    s = Solution()

    for tc in test_cases:
        result = s.FuncName(tc["input"])
		if result != tc["want"]:
			print(f"Test failed for input {tc["input"]}: got {result}, want {tc["want"]}")
		else:
			print(f"Test passed for input {tc["input"]}")


if __name__ == "__main__":
	main()
`

const defaultTemplate = `
{{ .Header }}

{{ .CodeSnippet }}
`
