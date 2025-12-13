package init

import (
	"fmt"
	"strings"

	"github.com/suda-3156/leetcode-cli/internal/tui/styles"
)

// View implements tea.Model for the confirmation prompt
func (m *Model) View() string {
	if m.quitting {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(styles.StyleWarningHeader.Render(
		fmt.Sprintf("\nConfig file already exists: %s\n\n", m.path),
	))
	sb.WriteString("\n")

	sb.WriteString("Do you want to overwrite it?\n\n")

	// Render choices
	choices := []string{"Yes", "No"}
	for i, choice := range choices {
		cursor := "  "
		if m.choice == i {
			cursor = styles.StyleActive.Render("> ")
			choice = styles.StyleActive.Render(choice)
		}
		sb.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}

	sb.WriteString("\n↑/k: up • ↓/j: down • enter: select • q/esc: quit\n")

	return sb.String()
}
