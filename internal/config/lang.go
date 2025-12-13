package config

type LangConfig struct {
	Extension     string
	CommentPrefix string
	TemplateFile  string
}

var LanguageConfigs = map[string]LangConfig{
	"golang": {
		Extension:     ".go",
		CommentPrefix: "//",
		TemplateFile:  "golang.tmpl",
	},
	"python3": {
		Extension:     ".py",
		CommentPrefix: "#",
		TemplateFile:  "python3.tmpl",
	},
	"python": {
		Extension:     ".py",
		CommentPrefix: "#",
	},
	"javascript": {
		Extension:     ".js",
		CommentPrefix: "//",
	},
	"typescript": {
		Extension:     ".ts",
		CommentPrefix: "//",
	},
	"java": {
		Extension:     ".java",
		CommentPrefix: "//",
	},
	"cpp": {
		Extension:     ".cpp",
		CommentPrefix: "//",
	},
	"c": {
		Extension:     ".c",
		CommentPrefix: "//",
	},
	"csharp": {
		Extension:     ".cs",
		CommentPrefix: "//",
	},
	"rust": {
		Extension:     ".rs",
		CommentPrefix: "//",
	},
	"ruby": {
		Extension:     ".rb",
		CommentPrefix: "#",
	},
	"swift": {
		Extension:     ".swift",
		CommentPrefix: "//",
	},
	"kotlin": {
		Extension:     ".kt",
		CommentPrefix: "//",
	},
	"scala": {
		Extension:     ".scala",
		CommentPrefix: "//",
	},
	"php": {
		Extension:     ".php",
		CommentPrefix: "//",
	},
	"dart": {
		Extension:     ".dart",
		CommentPrefix: "//",
	},
	"sql": {
		Extension:     ".sql",
		CommentPrefix: "--",
	},
}

// Default values:
//   - Extension: .txt
//   - CommentPrefix: //
//   - Template: "default.tmpl"
func GetLangConfig(langSlug string) *LangConfig {
	config, ok := LanguageConfigs[langSlug]
	if !ok {
		config = LangConfig{
			Extension:     ".txt",
			CommentPrefix: "//",
			TemplateFile:  "default.tmpl",
		}
	}

	if config.TemplateFile == "" {
		config.TemplateFile = "default.tmpl"
	}

	return &config
}
