package phases

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/suda-3156/leetcode-cli/internal/tui/model"
)

// PhaseType represents the different phases of the TUI flow
type PhaseType int

const (
	InitialPhase PhaseType = iota
	DecideQuestionPhase
	DecideLanguagePhase
	DecidePathPhase
	OverwriteConfirmPhase
	GenerationPhase
	DonePhase
)

// PhaseHandler defines the interface for phase implementations
type PhaseHandler interface {
	// Enter is called when transitioning into this phase
	Enter(m *model.Model) tea.Cmd

	// Update handles messages for this phase
	// Returns a command and optionally a new phase to transition to
	Update(m *model.Model, msg tea.Msg) (tea.Cmd, *PhaseType)

	// View renders the current phase
	View(m *model.Model) string
}

// GetHandler returns the appropriate handler for the given phase type
func GetHandler(phaseType PhaseType) PhaseHandler {
	switch phaseType {
	case InitialPhase:
		return &InitialHandler{}
	case DecideQuestionPhase:
		return &DecideQuestionHandler{}
	case DecideLanguagePhase:
		return &DecideLanguageHandler{}
	case DecidePathPhase:
		return &DecidePathHandler{}
	case OverwriteConfirmPhase:
		return &OverwriteConfirmHandler{}
	case GenerationPhase:
		return &GenerationHandler{}
	case DonePhase:
		return &DoneHandler{}
	default:
		return &DoneHandler{}
	}
}
