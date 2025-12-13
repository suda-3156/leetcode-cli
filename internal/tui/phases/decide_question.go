package phases

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/suda-3156/leetcode-cli/internal/api"
	"github.com/suda-3156/leetcode-cli/internal/config"
	"github.com/suda-3156/leetcode-cli/internal/tui/model"
	"github.com/suda-3156/leetcode-cli/internal/tui/styles"
)

type DecideQuestionHandler struct{}

type fetchQuestionListMsg struct {
	questionList []api.Question
}

func (h *DecideQuestionHandler) Enter(m *model.Model) tea.Cmd {
	m.Loading = true
	m.Cursor = 0
	return h.fetchQuestionList(m)
}

func (h *DecideQuestionHandler) fetchQuestionList(m *model.Model) tea.Cmd {
	return func() tea.Msg {
		resp, err := m.Client.SearchQuestions(m.Input.Keyword, config.DEFAULT_SEARCH_LIMIT, 0)
		if err != nil {
			return errMsg{err: fmt.Errorf("search questions: %w", err)}
		}

		questions := resp.Data.ProblemsetPanelQuestionList.Questions

		if len(questions) == 0 {
			return errMsg{err: fmt.Errorf("no questions found for keyword: %s", m.Input.Keyword)}
		}

		return fetchQuestionListMsg{questionList: questions}
	}
}

//nolint:cyclop // Handlers may have complex logic
func (h *DecideQuestionHandler) Update(m *model.Model, msg tea.Msg) (tea.Cmd, *PhaseType) {
	switch msg := msg.(type) {
	case fetchQuestionListMsg:
		m.Loading = false
		m.Questions = msg.questionList
		return nil, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Questions)-1 {
				m.Cursor++
			}
		case "enter":
			if len(m.Questions) > 0 {
				if m.Questions[m.Cursor].PaidOnly {
					m.Err = fmt.Errorf("premium-only question is not supported")
					return tea.Quit, nil
				}
				m.SelectedQ = &m.Questions[m.Cursor]
				next := DecideLanguagePhase
				return nil, &next
			}
		case "q", "esc":
			return tea.Quit, nil
		}

	case errMsg:
		m.Err = msg.err
		return tea.Quit, nil
	}

	return nil, nil
}

func (h *DecideQuestionHandler) View(m *model.Model) string {
	if m.Loading {
		return fmt.Sprintf("%s Searching for '%s'...\n", m.Spinner.View(), m.Input.Keyword)
	}

	var sb strings.Builder
	sb.WriteString("Select a question:\n\n")

	for i, q := range m.Questions {
		cursor := "  "
		if m.Cursor == i {
			cursor = styles.StyleActive.Render("> ")
		}

		premium := ""
		if q.PaidOnly {
			premium = " [Premium]"
		}

		difficulty := strings.ToUpper(q.Difficulty)
		line := fmt.Sprintf("%s. %s (%s)%s",
			q.QuestionFrontendID,
			q.Title,
			difficulty,
			premium,
		)

		if m.Cursor == i {
			line = styles.StyleActive.Render(line)
		}

		sb.WriteString(fmt.Sprintf("%s%s\n", cursor, line))
	}

	sb.WriteString("\n↑/k: up • ↓/j: down • enter: select • q/esc: quit\n")
	return sb.String()
}
