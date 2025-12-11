package config

import (
	"fmt"
	"strings"

	"github.com/suda-3156/leetcode-cli/internal/file"
)

const (
	FILEPATH_TMPL       = "./{{ .Date }}/{{ .FrontendID }}.{{ .TitleSlug }}{{ .Extension }}"
	DEFAULT_DATE_FORMAT = "2006-01-02"
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

// ResolveConfig loads the configuration from the given path and overrides
// fields with the provided parameters if they are not empty.
func ResolveConfig(configPath, langSlug, titleSlug, outPath, overwriteStr string) (*Config, error) {
	cfg, err := Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("resolve config: %w", err)
	}

	if langSlug != "" {
		cfg.LangSlug = langSlug
	}
	if titleSlug != "" {
		cfg.TitleSlug = titleSlug
	}
	if outPath != "" {
		cfg.OutPath = outPath
	}
	if overwriteStr != "" {
		switch strings.ToLower(overwriteStr) {
		case "always", "force", "true":
			cfg.Overwrite = OverwriteAlways
		case "backup":
			cfg.Overwrite = OverwriteBackup
		case "never":
			cfg.Overwrite = OverwriteNever
		default:
			cfg.Overwrite = OverwritePrompt
		}
	}

	return cfg, nil
}

func Load(configPath string) (*Config, error) {
	var data string
	var err error
	for _, path := range CONFIG_PATH_SLICE {
		if exists := file.FileExists(path); !exists {
			continue
		}
		data, err = file.Read(path)
		if err != nil {
			return nil, fmt.Errorf("config load: failed to read config file %s: %w", path, err)
		}
	}

	data = strings.TrimSpace(data)

	cfg := &Config{
		LangSlug:   "",
		TitleSlug:  "",
		OutPath:    FILEPATH_TMPL,
		DateFormat: DEFAULT_DATE_FORMAT,
		Overwrite:  OverwritePrompt,
	}

	// If no config file found, return default config
	if data == "" {
		return cfg, nil
	}

	cfg, err = parseConfig(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseConfig(data string, cfg *Config) (*Config, error) {
	cfgStr := &struct {
		LangSlug   string `yaml:"lang_slug"`
		TitleSlug  string `yaml:"title_slug"`
		OutPath    string `yaml:"out_path"`
		DateFormat string `yaml:"date_format"`
		Overwrite  string `yaml:"overwrite"`
	}{
		LangSlug:   "",
		TitleSlug:  "",
		OutPath:    FILEPATH_TMPL,
		DateFormat: DEFAULT_DATE_FORMAT,
		Overwrite:  "prompt",
	}

	if err := file.Parse(data, cfgStr); err != nil {
		return nil, fmt.Errorf("config parse: failed to parse config data: %w", err)
	}

	cfg.LangSlug = cfgStr.LangSlug
	cfg.TitleSlug = cfgStr.TitleSlug
	cfg.OutPath = cfgStr.OutPath

	switch strings.ToLower(cfgStr.Overwrite) {
	case "always", "force", "true":
		cfg.Overwrite = OverwriteAlways
	case "backup":
		cfg.Overwrite = OverwriteBackup
	case "never":
		cfg.Overwrite = OverwriteNever
	default:
		cfg.Overwrite = OverwritePrompt
	}

	return cfg, nil
}
