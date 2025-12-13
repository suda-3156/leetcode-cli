package phases

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/suda-3156/leetcode-cli/internal/config"
	"github.com/suda-3156/leetcode-cli/internal/tui/model"
	"github.com/suda-3156/leetcode-cli/internal/tui/styles"
)

type OverwriteConfirmHandler struct{}
type autoTransitionMsg struct{} // Message to trigger automatic transition

func (h *OverwriteConfirmHandler) Enter(m *model.Model) tea.Cmd {
	m.Loading = false
	m.OverwriteChoice = 0

	// Check config for automatic decision
	switch m.Config.Overwrite {
	case config.OverwriteAlways:
		m.OverwriteChoice = 1
		return func() tea.Msg { return autoTransitionMsg{} }

	case config.OverwriteBackup:
		m.OverwriteChoice = 2
		return func() tea.Msg { return autoTransitionMsg{} }

	case config.OverwriteNever:
		m.Err = fmt.Errorf("file already exists and overwrite is disabled: %s", m.OutPath)
		return tea.Quit

	case config.OverwritePrompt:
		// Show prompt to user
		return nil
	}

	return nil
}

//nolint:cyclop // Handlers may have complex logic
func (h *OverwriteConfirmHandler) Update(m *model.Model, msg tea.Msg) (tea.Cmd, *PhaseType) {
	switch msg := msg.(type) {
	case autoTransitionMsg:
		// Automatic transition triggered by Enter()
		next := GenerationPhase
		return nil, &next

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.OverwriteChoice > 0 {
				m.OverwriteChoice--
			}
		case "down", "j":
			if m.OverwriteChoice < 3 {
				m.OverwriteChoice++
			}
		case "r":
			m.OverwriteChoice = config.OverwritePrompt // Reset choice
			m.TextInput.Focus()
			next := DecidePathPhase
			return nil, &next
		case "o":
			m.OverwriteChoice = config.OverwriteAlways
			next := GenerationPhase
			return nil, &next
		case "b":
			m.OverwriteChoice = config.OverwriteBackup
			next := GenerationPhase
			return nil, &next
		case "enter":
			switch m.OverwriteChoice {
			case config.OverwritePrompt: // Return to path input
				m.TextInput.Focus()
				next := DecidePathPhase
				return nil, &next
			case config.OverwriteAlways:
				next := GenerationPhase
				return nil, &next
			case config.OverwriteBackup:
				next := GenerationPhase
				return nil, &next
			case config.OverwriteNever:
				return tea.Quit, nil
			}
		case "q", "esc":
			return tea.Quit, nil
		}
	default:
		return nil, nil
	}

	return nil, nil
}

func (h *OverwriteConfirmHandler) View(m *model.Model) string {
	var sb strings.Builder

	sb.WriteString(styles.StyleWarningHeader.Render("File already exists"))
	sb.WriteString(fmt.Sprintf("\n\n%s\n\n", m.OutPath))
	sb.WriteString("What would you like to do?\n\n")

	choices := []string{
		"Return to path input",
		"Overwrite the existing file",
		"Create backup and overwrite",
		"Quit",
	}

	for i, choice := range choices {
		cursor := "  "
		line := choice

		if m.OverwriteChoice == config.OverwriteOption(i) {
			cursor = styles.StyleActive.Render("> ")
			line = styles.StyleActive.Render(line)
		}

		sb.WriteString(fmt.Sprintf("%s%s\n", cursor, line))
	}

	sb.WriteString("\n↑/k: up • ↓/j: down • enter: select • o: overwrite • b: backup • r: return • q/esc: quit\n")

	return sb.String()
}
