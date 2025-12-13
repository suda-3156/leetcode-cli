package config

import "time"

const (
	FILEPATH_TMPL        = "./{{ .Date }}/{{ .FrontendID }}.{{ .TitleSlug }}{{ .Extension }}"
	DEFAULT_DATE_FORMAT  = "2006-01-02"
	DEFAULT_SEARCH_LIMIT = 20
)

var (
	CONFIG_PATH_SLICE = []string{
		".leetcode-cli.yaml",
		".leetcode-cli.yml",
		"leetcode-cli.yaml",
		"leetcode-cli.yml",
	}
)

type OverwriteOption int

const (
	OverwritePrompt = iota // 0 Ask user to overwrite, bakcup, or quit
	OverwriteAlways        // 1
	OverwriteBackup        // 2
	OverwriteNever         // 3
)

type Config struct {
	LangSlug   string
	TitleSlug  string
	OutPath    string // Output path template
	DateFormat string
	Overwrite  OverwriteOption
}

func (cfg *Config) GetCurrentDate() string {
	return time.Now().Format(cfg.DateFormat)
}
