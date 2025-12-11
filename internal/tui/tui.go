package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/suda-3156/leetcode-cli/internal/config"
)

func (m *Model) Init() tea.Cmd {
	// If --slug is specified, fetch question detail directly
	if m.presetSlug != "" {
		return m.fetchQuestionDetail
	}
	return m.fetchQuestionList
}

//nolint:cyclop // The main update function.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:funlen // main update function
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Common key handling
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case questionListMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = StateError
			return m, tea.Quit
		}
		if len(msg.questions) == 0 {
			m.err = fmt.Errorf("no questions found for keyword: %s", m.keyword)
			m.state = StateError
			return m, tea.Quit
		}
		m.questions = msg.questions
		m.state = StateQuestionList
		m.cursor = 0
		return m, nil

	case questionDetailMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = StateError
			return m, tea.Quit
		}
		m.questionDetail = msg.detail
		m.languages = msg.detail.CodeSnippets
		m.cursor = 0

		// If --lang is specified, skip language selection
		if m.presetLang != "" {
			lang := findLanguageBySlug(m.languages, m.presetLang)
			if lang == nil {
				m.err = fmt.Errorf("language '%s' is not available. Available languages: %v",
					m.presetLang, getAvailableLanguages(m.languages))
				m.state = StateError
				return m, tea.Quit
			}
			m.selectedLang = lang
			m.state = StatePathInput

			// If --path is also specified, skip path input
			if m.presetPath != "" {
				return m, m.generateFile
			}

			// Set default path
			langConfig := config.GetLangConfig(m.selectedLang.LangSlug)
			defaultPath := config.GetDefaultOutputPath(
				m.questionDetail.QuestionFrontendID,
				m.questionDetail.TitleSlug,
				langConfig.Extension,
			)
			m.textInput.SetValue(defaultPath)
			m.textInput.Focus()
			return m, nil
		}

		m.state = StateLanguageSelect
		return m, nil

	case fileGeneratedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = StateError
			return m, tea.Quit
		}
		m.generatedPath = msg.path
		m.state = StateDone
		return m, tea.Quit
	}

	// State-specific updates
	switch m.state {
	case StateQuestionList:
		return m.updateQuestionList(msg)
	case StateLanguageSelect:
		return m.updateLanguageSelect(msg)
	case StatePathInput:
		return m.updatePathInput(msg)
	default:
		return m, nil
	}
}

func (m *Model) View() string {
	switch m.state {
	case StateSearching:
		if m.presetSlug != "" {
			return fmt.Sprintf("Fetching question details for '%s'...\n", m.presetSlug)
		}
		return fmt.Sprintf("Searching for '%s'...\n", m.keyword)
	case StateQuestionList:
		return m.viewQuestionList()
	case StateLanguageSelect:
		return m.viewLanguageSelect()
	case StatePathInput:
		return m.viewPathInput()
	case StateDone:
		return fmt.Sprintf("✓ File generated: %s\n", m.generatedPath)
	case StateError:
		return fmt.Sprintf("Error: %v\n", m.err)
	}
	return ""
}

func Run(keyword, slug, lang, path string) (string, error) {
	m := NewModel(keyword, slug, lang, path)
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

func FormatQuestionDisplay(frontendID, title, difficulty string, paidOnly bool) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s. %s (%s)", frontendID, title, strings.ToUpper(difficulty)))
	if paidOnly {
		sb.WriteString(" [Premium]")
	}
	return sb.String()
}
