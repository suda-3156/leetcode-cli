package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/suda-3156/leetcode-cli/internal/api"
	"github.com/suda-3156/leetcode-cli/internal/config"
)

type questionListMsg struct {
	questions []api.Question
	err       error
}

func (m *Model) fetchQuestionList() tea.Msg {
	resp, err := m.client.SearchQuestions(m.keyword, config.DefaultSearchLimit, 0)
	if err != nil {
		return questionListMsg{err: err}
	}
	return questionListMsg{questions: resp.Data.ProblemsetPanelQuestionList.Questions}
}

func (m *Model) updateQuestionList(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
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
				m.state = StateSearching // 詳細取得中
				return m, m.fetchQuestionDetail
			}
		case "q", "esc":
			return m, tea.Quit
		}
	default:
		return m, nil
	}
	return m, nil
}

func (m *Model) viewQuestionList() string {
	var sb strings.Builder
	sb.WriteString("Select a question:\n\n")

	for i, q := range m.questions {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}

		premium := ""
		if q.PaidOnly {
			premium = " [Premium]"
		}

		difficulty := strings.ToUpper(q.Difficulty)
		sb.WriteString(fmt.Sprintf("%s%s. %s (%s)%s\n",
			cursor,
			q.QuestionFrontendID,
			q.Title,
			difficulty,
			premium,
		))
	}

	sb.WriteString("\n↑/k: up • ↓/j: down • enter: select • q/esc: quit\n")
	return sb.String()
}
