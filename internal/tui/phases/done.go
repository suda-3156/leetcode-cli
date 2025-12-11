package phases

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/suda-3156/leetcode-cli/internal/tui/model"
)

type DoneHandler struct{}

func (h *DoneHandler) Enter(m *model.Model) tea.Cmd {
	return tea.Quit
}

func (h *DoneHandler) Update(m *model.Model, msg tea.Msg) (tea.Cmd, *PhaseType) {
	return tea.Quit, nil
}

func (h *DoneHandler) View(m *model.Model) string {
	return fmt.Sprintf("✓ File generated: %s\n", m.OutPath)
}
