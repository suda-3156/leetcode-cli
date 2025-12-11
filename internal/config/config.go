package config

type OverwriteOption int

const (
	OverwritePrompt = iota // 0 Ask user to overwrite, bakcup, or quit
	OverwriteAlways        // 1
	OverwriteBackup        // 2
	OverwriteNever         // 3
)

type Config struct {
	Language  string
	TitleSlug string
	OutPath   string
	Overwrite OverwriteOption
}

func Load() (*Config, error) {
	// Implementation for loading configuration
	// Resolve options, config file, and defaults here
	// For now, return a default config
	return &Config{
		Language:  "",
		TitleSlug: "",
		OutPath:   "",
		Overwrite: OverwritePrompt,
	}, nil
}
