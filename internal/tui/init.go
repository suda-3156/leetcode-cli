package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/suda-3156/leetcode-cli/internal/config"
)

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.entryInitialPhase(),
	)
}

func Run(keyword, slug, lang, path string) (string, error) {
	// Create config from flags
	cfg := &config.Config{
		Language:  lang,
		TitleSlug: slug,
		OutPath:   path,
		Overwrite: config.OverwritePrompt,
	}

	m := New(keyword, cfg)
	p := tea.NewProgram(&m)

	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	fm := finalModel.(*Model)
	if fm.HasError() {
		return "", fm.GetError()
	}

	return fm.GetGeneratedPath(), nil
}
