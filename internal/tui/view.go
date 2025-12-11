package tui

import (
	"fmt"

	"github.com/suda-3156/leetcode-cli/internal/tui/styles"
)

func (m *Model) View() string {
	// Show error if present
	if m.Err != nil {
		return m.viewError()
	}

	// Delegate to current phase handler
	return m.currentHandler.View(&m.Model)
}

func (m *Model) viewError() string {
	header := styles.StyleErrorHeader.Render("Error")
	return fmt.Sprintf("%s\n\n%v\n", header, m.Err)
}
