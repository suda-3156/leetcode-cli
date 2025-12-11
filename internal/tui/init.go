package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/suda-3156/leetcode-cli/internal/tui/model"
)

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick,
		m.currentHandler.Enter(&m.Model),
	)
}

func Run(keyword, slug, lang, path string) (string, error) {
	input := model.Input{
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
