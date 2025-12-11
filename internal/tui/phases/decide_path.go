package phases

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/suda-3156/leetcode-cli/internal/file"
	"github.com/suda-3156/leetcode-cli/internal/generator"
	"github.com/suda-3156/leetcode-cli/internal/tui/model"
)

type DecidePathHandler struct{}

func (h *DecidePathHandler) Enter(m *model.Model) tea.Cmd {
	m.Loading = false

	// Generate default path
	outPath := generator.GetOutputPath(
		m.Config,
		m.QuestionDetail.TitleSlug,
		m.QuestionDetail.QuestionFrontendID,
		m.SelectedLang.LangSlug,
	)

	m.OutPath = outPath

	// Show input for user to enter/modify path
	m.TextInput.SetValue(outPath)
	m.TextInput.Focus()

	return nil
}

func (h *DecidePathHandler) Update(m *model.Model, msg tea.Msg) (tea.Cmd, *PhaseType) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.OutPath = m.TextInput.Value()

			// Check if file exists
			if file.FileExists(m.OutPath) {
				next := OverwriteConfirmPhase
				return nil, &next
			}

			// File doesn't exist, go directly to generation
			next := GenerationPhase
			return nil, &next

		case "h", "left":
			// Go back to language selection
			m.Cursor = 0
			next := DecideLanguagePhase
			return nil, &next

		case "esc":
			return tea.Quit, nil
		}
	}

	var cmd tea.Cmd
	m.TextInput, cmd = m.TextInput.Update(msg)
	return cmd, nil
}

func (h *DecidePathHandler) View(m *model.Model) string {
	title := "unknown question"
	lang := "unknown language"

	if m.QuestionDetail != nil {
		title = m.QuestionDetail.Title
	}
	if m.SelectedLang != nil {
		lang = m.SelectedLang.Lang
	}

	return fmt.Sprintf(
		"Output path for %q (%s):\n\n%s\n\nenter: confirm • h/left: back • esc: quit\n",
		title,
		lang,
		m.TextInput.View(),
	)
}
