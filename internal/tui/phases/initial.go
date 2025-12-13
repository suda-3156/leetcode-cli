package phases

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/suda-3156/leetcode-cli/internal/config"
	"github.com/suda-3156/leetcode-cli/internal/tui/model"
)

type InitialHandler struct{}

type loadConfigMsg struct {
	config *config.Config
}

type errMsg struct {
	err error
}

func (h *InitialHandler) Enter(m *model.Model) tea.Cmd {
	m.Loading = true
	return h.loadConfig(m)
}

func (h *InitialHandler) loadConfig(m *model.Model) tea.Cmd {
	return func() tea.Msg {
		cfg, err := config.ResolveConfig(
			m.Input.ConfigPath,
			m.Input.LangSlug,
			m.Input.TitleSlug,
			m.Input.OutPath,
			m.Input.OverwriteStr,
		)
		if err != nil {
			return errMsg{err: fmt.Errorf("load config: %w", err)}
		}
		return loadConfigMsg{config: cfg}
	}
}

func (h *InitialHandler) Update(m *model.Model, msg tea.Msg) (tea.Cmd, *PhaseType) {
	switch msg := msg.(type) {
	case loadConfigMsg:
		m.Loading = false
		m.Config = msg.config

		// Determine next phase based on config
		if m.Config.TitleSlug != "" {
			// Skip question selection if slug is provided
			next := DecideLanguagePhase
			return nil, &next
		}

		// Need to search for questions
		next := DecideQuestionPhase
		return nil, &next

	case errMsg:
		m.Err = msg.err
		return tea.Quit, nil
	}

	return nil, nil
}

func (h *InitialHandler) View(m *model.Model) string {
	return fmt.Sprintf("%s Loading configuration...\n", m.Spinner.View())
}
