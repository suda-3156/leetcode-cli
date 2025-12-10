package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/suda-3156/leetcode-cli/internal/api"
	"github.com/suda-3156/leetcode-cli/internal/config"
)

type questionDetailMsg struct {
	detail *api.QuestionDetail
	err    error
}

func (m Model) fetchQuestionDetail() tea.Msg {
	var titleSlug string
	if m.presetSlug != "" {
		titleSlug = m.presetSlug
	} else if m.selectedQ != nil {
		titleSlug = m.selectedQ.TitleSlug
	} else {
		return questionDetailMsg{err: fmt.Errorf("no question selected")}
	}

	resp, err := m.client.GetQuestionDetail(titleSlug)
	if err != nil {
		return questionDetailMsg{err: err}
	}
	return questionDetailMsg{detail: &resp.Data.Question}
}

func (m Model) updateLanguageSelect(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.languages)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.languages) > 0 {
				m.selectedLang = &m.languages[m.cursor]
				m.state = StatePathInput
				m.cursor = 0

				// Set default path
				langConfig, ok := config.GetLangConfig(m.selectedLang.LangSlug)
				if !ok {
					langConfig = config.LangConfig{Extension: ".txt"}
				}
				defaultPath := config.GetDefaultOutputPath(
					m.questionDetail.QuestionFrontendID,
					m.questionDetail.TitleSlug,
					langConfig.Extension,
				)
				m.textInput.SetValue(defaultPath)
				m.textInput.Focus()

				// Skip if --path flag is specified
				if m.presetPath != "" {
					return m, m.generateFile
				}

				return m, nil
			}
		case "q", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) viewLanguageSelect() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Select a language for \"%s\":\n\n", m.questionDetail.Title))

	for i, lang := range m.languages {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}
		sb.WriteString(fmt.Sprintf("%s%s\n", cursor, lang.Lang))
	}

	sb.WriteString("\n↑/k: up • ↓/j: down • enter: select • q/esc: quit\n")
	return sb.String()
}

func findLanguageBySlug(snippets []api.CodeSnippet, langSlug string) *api.CodeSnippet {
	for i := range snippets {
		if snippets[i].LangSlug == langSlug {
			return &snippets[i]
		}
	}
	return nil
}

func getAvailableLanguages(snippets []api.CodeSnippet) []string {
	langs := make([]string, len(snippets))
	for i, s := range snippets {
		langs[i] = s.LangSlug
	}
	return langs
}
