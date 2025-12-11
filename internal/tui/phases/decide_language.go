package phases

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/suda-3156/leetcode-cli/internal/api"
	"github.com/suda-3156/leetcode-cli/internal/tui/model"
	"github.com/suda-3156/leetcode-cli/internal/tui/styles"
)

type DecideLanguageHandler struct{}

type questionDetailMsg struct {
	detail *api.QuestionDetail
}

func (h *DecideLanguageHandler) Enter(m *model.Model) tea.Cmd {
	m.Loading = true
	m.Cursor = 0
	return h.fetchQuestionDetail(m)
}

func (h *DecideLanguageHandler) fetchQuestionDetail(m *model.Model) tea.Cmd {
	return func() tea.Msg {
		// Get titleSlug from selected question or config
		titleSlug := m.Config.TitleSlug
		if titleSlug == "" && m.SelectedQ != nil {
			titleSlug = m.SelectedQ.TitleSlug
		}

		if titleSlug == "" {
			return errMsg{err: fmt.Errorf("no title slug available")}
		}

		resp, err := m.Client.GetQuestionDetail(titleSlug)
		if err != nil {
			return errMsg{err: fmt.Errorf("get question detail: %w", err)}
		}
		return questionDetailMsg{detail: &resp.Data.Question}
	}
}

//nolint:cyclop // Handlers may have complex logic
func (h *DecideLanguageHandler) Update(m *model.Model, msg tea.Msg) (tea.Cmd, *PhaseType) {
	switch msg := msg.(type) {
	case questionDetailMsg:
		m.Loading = false
		m.QuestionDetail = msg.detail

		// Check if language is pre-specified in config
		if m.Config.LangSlug != "" {
			langSnippet := findCodeSnippetByLang(m.QuestionDetail.CodeSnippets, m.Config.LangSlug)
			if langSnippet != nil {
				m.SelectedLang = langSnippet
				next := DecidePathPhase
				return nil, &next
			}
			// If specified language not found, fall through to show list with error message
		}

		m.Languages = m.QuestionDetail.CodeSnippets
		return nil, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Languages)-1 {
				m.Cursor++
			}
		case "enter":
			if len(m.Languages) > 0 {
				m.SelectedLang = &m.Languages[m.Cursor]
				m.Cursor = 0
				next := DecidePathPhase
				return nil, &next
			}
		// case "h", "left":
		// 	// Go back to question selection if we came from there
		// 	if m.SelectedQ != nil {
		// 		m.Cursor = 0
		// 		next := DecideQuestionPhase
		// 		return nil, &next
		// 	}
		case "q", "esc":
			return tea.Quit, nil
		}

	case errMsg:
		m.Err = msg.err
		return tea.Quit, nil
	}

	return nil, nil
}

//nolint:cyclop // Handlers may have complex logic
func (h *DecideLanguageHandler) View(m *model.Model) string {
	if m.Loading {
		slug := m.Config.TitleSlug
		if slug == "" && m.SelectedQ != nil {
			slug = m.SelectedQ.TitleSlug
		}
		return fmt.Sprintf("%s Fetching question details for '%s'...\n", m.Spinner.View(), slug)
	}

	var sb strings.Builder

	// Show error if config language was invalid
	if m.Config.LangSlug != "" && m.QuestionDetail != nil {
		if findCodeSnippetByLang(m.QuestionDetail.CodeSnippets, m.Config.LangSlug) == nil {
			sb.WriteString(styles.StyleErrorHeader.Render(
				fmt.Sprintf("Warning: Language '%s' not available for this question\n", m.Config.LangSlug),
			))
			sb.WriteString("\n")
		}
	}

	if m.QuestionDetail != nil {
		sb.WriteString(fmt.Sprintf("Select a language for %q:\n\n", m.QuestionDetail.Title))
	} else {
		sb.WriteString("Select a language:\n\n")
	}

	for i, lang := range m.Languages {
		cursor := "  "
		line := lang.Lang

		if m.Cursor == i {
			cursor = styles.StyleActive.Render("> ")
			line = styles.StyleActive.Render(line)
		}

		sb.WriteString(fmt.Sprintf("%s%s\n", cursor, line))
	}

	helpText := "\n↑/k: up • ↓/j: down • enter: select"
	if m.SelectedQ != nil {
		helpText += " • h/left: back"
	}
	helpText += " • q/esc: quit\n"
	sb.WriteString(helpText)

	return sb.String()
}

func findCodeSnippetByLang(snippets []api.CodeSnippet, lang string) *api.CodeSnippet {
	for _, snippet := range snippets {
		if snippet.LangSlug == lang {
			return &snippet
		}
	}
	return nil
}
