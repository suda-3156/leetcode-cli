package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/suda-3156/leetcode-cli/internal/generator"
)

type fileGeneratedMsg struct {
	path string
	err  error
}

func (m Model) generateFile() tea.Msg {
	outputPath := generator.GetOutputPath(
		m.presetPath,
		m.questionDetail.QuestionFrontendID,
		m.questionDetail.TitleSlug,
		m.selectedLang.LangSlug,
	)

	// Use text input value if presetPath is not set
	if m.presetPath == "" && m.textInput.Value() != "" {
		outputPath = m.textInput.Value()
	}

	err := generator.GenerateFile(outputPath, m.questionDetail, m.selectedLang)
	if err != nil {
		return fileGeneratedMsg{err: err}
	}
	return fileGeneratedMsg{path: outputPath}
}

func (m Model) updatePathInput(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, m.generateFile
		case "esc":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) viewPathInput() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Output path for \"%s\" (%s):\n\n",
		m.questionDetail.Title,
		m.selectedLang.Lang,
	))
	sb.WriteString(m.textInput.View())
	sb.WriteString("\n\nenter: confirm • esc: quit\n")
	return sb.String()
}
