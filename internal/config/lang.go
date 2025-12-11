package config

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

// Default values:
//   - Extension: .txt
//   - CommentPrefix: //
func GetLangConfig(langSlug string) LangConfig {
	config, ok := LanguageConfigs[langSlug]
	if !ok {
		config = LangConfig{
			Extension:     ".txt",
			CommentPrefix: "//",
		}
	}

	return config
}
