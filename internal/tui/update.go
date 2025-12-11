package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	// Let the current phase handler process the message
	cmd, nextPhase := m.currentHandler.Update(&m.Model, msg)

	// If a phase transition is requested, transition and call Enter
	if nextPhase != nil {
		transitionCmd := m.transitionTo(*nextPhase)
		if cmd != nil {
			return m, tea.Batch(cmd, transitionCmd)
		}
		return m, transitionCmd
	}

	return m, cmd
}
