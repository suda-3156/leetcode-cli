package config

import "time"

const DefaultSearchLimit = 20

type LangConfig struct {
	Extension     string
	CommentPrefix string
}

var LanguageConfigs = map[string]LangConfig{
	"golang":     {Extension: ".go", CommentPrefix: "//"},
	"python3":    {Extension: ".py", CommentPrefix: "#"},
	"python":     {Extension: ".py", CommentPrefix: "#"},
	"javascript": {Extension: ".js", CommentPrefix: "//"},
	"typescript": {Extension: ".ts", CommentPrefix: "//"},
	"java":       {Extension: ".java", CommentPrefix: "//"},
	"cpp":        {Extension: ".cpp", CommentPrefix: "//"},
	"c":          {Extension: ".c", CommentPrefix: "//"},
	"csharp":     {Extension: ".cs", CommentPrefix: "//"},
	"rust":       {Extension: ".rs", CommentPrefix: "//"},
	"ruby":       {Extension: ".rb", CommentPrefix: "#"},
	"swift":      {Extension: ".swift", CommentPrefix: "//"},
	"kotlin":     {Extension: ".kt", CommentPrefix: "//"},
	"scala":      {Extension: ".scala", CommentPrefix: "//"},
	"php":        {Extension: ".php", CommentPrefix: "//"},
	"dart":       {Extension: ".dart", CommentPrefix: "//"},
	"sql":        {Extension: ".sql", CommentPrefix: "--"},
}

func GetLangConfig(langSlug string) (LangConfig, bool) {
	config, ok := LanguageConfigs[langSlug]
	return config, ok
}

// GetDefaultOutputPath generates the default output path: ./YYYY-MM-DD/<frontendID>.<titleSlug><extension>
func GetDefaultOutputPath(frontendID, titleSlug, extension string) string {
	date := time.Now().Format("2006-01-02")
	return "./" + date + "/" + frontendID + "." + titleSlug + extension
}

func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}
