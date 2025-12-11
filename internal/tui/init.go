package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.entryInitialPhase(),
	)
}

// func Run(keyword, slug, lang, path string) (string, error) {
// 	m := NewModel(keyword, slug, lang, path)
// 	p := tea.NewProgram(&m)

// 	finalModel, err := p.Run()
// 	if err != nil {
// 		return "", err
// 	}

// 	fm := finalModel.(*Model)
// 	if fm.HasError() {
// 		return "", fm.GetError()
// 	}

// 	return fm.GetGeneratedPath(), nil
// }

// func FormatQuestionDisplay(frontendID, title, difficulty string, paidOnly bool) string {
// 	var sb strings.Builder
// 	sb.WriteString(fmt.Sprintf("%s. %s (%s)", frontendID, title, strings.ToUpper(difficulty)))
// 	if paidOnly {
// 		sb.WriteString(" [Premium]")
// 	}
// 	return sb.String()
// }
