package config

type LangConfig struct {
	Extension     string
	CommentPrefix string
	FuncDefRegex  string
	TemplateFile  string
}

var LanguageConfigs = map[string]LangConfig{
	"golang": {
		Extension:     ".go",
		CommentPrefix: "//",
		FuncDefRegex:  `func (\w+)\s*\(`,
		TemplateFile:  "golang.tmpl",
	},
	"python3": {
		Extension:     ".py",
		CommentPrefix: "#",
		// ```py
		// # Definition for a binary tree node.
		// # class TreeNode:
		// #     def __init__(self, val=0, left=None, right=None):
		// #         self.val = val
		// #         self.left = left
		// #         self.right = right
		//
		// class Solution:
		//     def preorderTraversal(self, root: Optional[TreeNode]) -> List[int]:
		// ```
		// We want to capture "preorderTraversal" here.
		FuncDefRegex: `(?s)class\s+Solution:.*?def (\w+)\s*\(`,
		TemplateFile: "python.tmpl",
	},
	"python": {
		Extension:     ".py",
		CommentPrefix: "#",
		FuncDefRegex:  `(?s)class\s+Solution:.*?def (\w+)\s*\(`,
		TemplateFile:  "python.tmpl",
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
//   - FuncDefRegex: ""
//   - Template: "default.tmpl"
func GetLangConfig(langSlug string) *LangConfig {
	config, ok := LanguageConfigs[langSlug]
	if !ok {
		config = LangConfig{
			Extension:     ".txt",
			CommentPrefix: "//",
			FuncDefRegex:  "",
			TemplateFile:  "default.tmpl",
		}
	}

	if config.TemplateFile == "" {
		config.TemplateFile = "default.tmpl"
	}

	return &config
}
