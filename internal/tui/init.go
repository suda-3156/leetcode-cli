package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.entryInitialPhase(),
	)
}

func Run(keyword, slug, lang, path string) (string, error) {
	input := Input{
		Keyword:   keyword,
		TitleSlug: slug,
		LangSlug:  lang,
		OutPath:   path,
	}

	m := New(input)
	p := tea.NewProgram(&m)

	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("run tui: %w", err)
	}

	fm := finalModel.(*Model)
	if fm.HasError() {
		return "", fm.GetError()
	}

	return fm.GetGeneratedPath(), nil
}
