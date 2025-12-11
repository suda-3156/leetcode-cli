package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/suda-3156/leetcode-cli/internal/api"
	"github.com/suda-3156/leetcode-cli/internal/config"
	"github.com/suda-3156/leetcode-cli/internal/generator"
)

type errMsg struct{ err error }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg.err
		return m, nil
	}

	switch m.phase {
	case InitialPhase:
		return m.updateInitialPhase(msg)
	case DecideQuestionPhase:
		return m.updateDecideQuestionPhase(msg)
	case DecideLanguagePhase:
		return m.updateDecideLanguagePhase(msg)
	case DecidePathPhase:
		return m.updateDecidePathPhase(msg)
	case GenerationPhase:
		return m.updateGenerationPhase(msg)
	case DonePhase:
		return m, nil
	}
	return m, nil
}

/**
 * Initial Phase
 *  - Load config
 */
func (m *Model) entryInitialPhase() tea.Cmd {
	m.phase = InitialPhase
	return m.loadConfig
}

type loadConfigMsg struct{ config *config.Config }

func (m *Model) loadConfig() tea.Msg {
	cfg, err := config.Load()
	if err != nil {
		return errMsg{err}
	}

	return loadConfigMsg{config: cfg}
}

func (m *Model) updateInitialPhase(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case loadConfigMsg:
		m.config = msg.config
		return m, m.entryDecideQuestionPhase()
	default:
		return m, nil
	}
}

/**
 * Decide Question Phase
 *  - Search questions or use slug from config
 */
func (m *Model) entryDecideQuestionPhase() tea.Cmd {
	m.phase = DecideQuestionPhase

	if m.config.TitleSlug != "" {
		return m.entryDecideLanguagePhase()
	}

	return m.fetchQuestionList
}

type fetchQuestionListMsg struct {
	questionList []api.Question
}

func (m *Model) fetchQuestionList() tea.Msg {
	resp, err := m.client.SearchQuestions(m.keyword, config.DefaultSearchLimit, 0)
	if err != nil {
		return errMsg{err}
	}

	questions := resp.Data.ProblemsetPanelQuestionList.Questions

	if len(questions) == 0 {
		return errMsg{fmt.Errorf("no questions found for keyword: %s", m.keyword)}
	}

	return fetchQuestionListMsg{
		questionList: questions,
	}
}

func (m *Model) updateDecideQuestionPhase(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case fetchQuestionListMsg:
		m.questions = msg.questionList
		m.cursor = 0
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.questions)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.questions) > 0 {
				m.selectedQ = &m.questions[m.cursor]
				return m, m.entryDecideLanguagePhase()
			}
		case "q", "esc":
			return m, tea.Quit
		}
	default:
		return m, nil
	}
	return m, nil
}

/**
 * Decide Language Phase
 *  - Fetch question detail by slug
 *  - Show language list and handle selection or use lang from config
 */
func (m *Model) entryDecideLanguagePhase() tea.Cmd {
	m.phase = DecideLanguagePhase
	return m.fetchQuestionDetail
}

type questionDetailMsg struct {
	detail *api.QuestionDetail
}

func (m *Model) fetchQuestionDetail() tea.Msg {
	titleSlug := m.selectedQ.TitleSlug

	resp, err := m.client.GetQuestionDetail(titleSlug)
	if err != nil {
		return errMsg{err}
	}
	return questionDetailMsg{detail: &resp.Data.Question}
}

func (m *Model) updateDecideLanguagePhase(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case questionDetailMsg:
		m.questionDetail = msg.detail

		if m.config.Language != "" {
			langSnippet := findCodeSnippetByLang(m.questionDetail.CodeSnippets, m.config.Language)
			if langSnippet != nil {
				m.selectedLang = langSnippet
				return m, m.entryDecidePathPhase()
			}
		}

		m.languages = m.questionDetail.CodeSnippets
		m.cursor = 0
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.questions)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.questions) > 0 {
				m.selectedLang = &m.languages[m.cursor]
				m.cursor = 0
				return m, m.entryDecidePathPhase()
			}
		case "q", "esc":
			return m, tea.Quit
		}
	default:
		return m, nil
	}
	return m, nil
}

func findCodeSnippetByLang(snippets []api.CodeSnippet, lang string) *api.CodeSnippet {
	for _, snippet := range snippets {
		if snippet.LangSlug == lang {
			return &snippet
		}
	}
	return nil
}

/**
 * Decide Path Phase
 *  - Show path input and handle input
 */
func (m *Model) entryDecidePathPhase() tea.Cmd {
	m.phase = DecidePathPhase
	if m.config.OutPath != "" {
		m.textInput.SetValue(m.config.OutPath)
		return m.entryGenerationPhase()
	}
	return nil
}

func (m *Model) updateDecidePathPhase(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, m.entryGenerationPhase()
		case "esc":
			return m, tea.Quit
		}
	default:
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

/**
 * Generation Phase
 *  - Generate file and show result
 */
func (m *Model) entryGenerationPhase() tea.Cmd {
	m.phase = GenerationPhase
	return m.generateFile
}

type fileGeneratedMsg struct {
	path string
}

func (m *Model) generateFile() tea.Msg {
	outputPath := generator.GetOutputPath(
		m.config.OutPath,
		m.questionDetail.QuestionFrontendID,
		m.questionDetail.TitleSlug,
		m.selectedLang.LangSlug,
	)

	// Use text input value if presetPath is not set
	if m.config.OutPath == "" && m.textInput.Value() != "" {
		outputPath = m.textInput.Value()
	}

	err := generator.GenerateFile(outputPath, m.questionDetail, m.selectedLang)
	if err != nil {
		return errMsg{err}
	}
	return fileGeneratedMsg{
		path: outputPath,
	}
}

func (m *Model) updateGenerationPhase(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case fileGeneratedMsg:
		m.phase = DonePhase
		m.generatedPath = msg.path
		return m, nil
	default:
		return m, nil
	}
}
