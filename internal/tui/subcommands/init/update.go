package init

import tea "github.com/charmbracelet/bubbletea"

// Update implements tea.Model for the confirmation prompt
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit
		case "up", "k", "left", "h":
			if m.choice > 0 {
				m.choice--
			}
		case "down", "j", "right", "l":
			if m.choice < 1 {
				m.choice++
			}
		case "y", "Y":
			m.confirmed = true
			m.quitting = true
			return m, tea.Quit
		case "n", "N":
			m.confirmed = false
			m.quitting = true
			return m, tea.Quit
		case "enter":
			m.confirmed = (m.choice == 0)
			m.quitting = true
			return m, tea.Quit
		}
	default:
		return m, nil
	}
	return m, nil
}
