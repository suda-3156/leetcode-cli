package generator

import (
	"embed"
)

//go:embed template/*.tmpl
var templateFS embed.FS

// GetTemplate returns the template content for the given filename.
// If the file is not found, it returns the default template.
func GetTemplate(filename string) (string, error) {
	if filename == "" {
		filename = "default.tmpl"
	}

	content, err := templateFS.ReadFile("template/" + filename)
	if err != nil {
		// Fallback to default template
		content, err = templateFS.ReadFile("template/default.tmpl")
		if err != nil {
			return "", err
		}
	}

	return string(content), nil
}
