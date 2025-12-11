package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/suda-3156/leetcode-cli/internal/api"
	"github.com/suda-3156/leetcode-cli/internal/tui/model"
	"github.com/suda-3156/leetcode-cli/internal/tui/phases"
	"github.com/suda-3156/leetcode-cli/internal/tui/styles"
)

type Model struct {
	model.Model
	currentPhase   phases.PhaseType
	currentHandler phases.PhaseHandler
}

func New(input model.Input) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter output path..."
	ti.CharLimit = 256
	ti.Width = 60

	sp := spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(styles.StyleSpinner),
	)

	m := Model{
		Model: model.Model{
			Spinner:   sp,
			Loading:   true,
			Input:     input,
			TextInput: ti,
			Client:    api.NewClient(),
		},
		currentPhase:   phases.InitialPhase,
		currentHandler: phases.GetHandler(phases.InitialPhase),
	}

	return m
}

func (m *Model) transitionTo(phaseType phases.PhaseType) tea.Cmd {
	m.currentPhase = phaseType
	m.currentHandler = phases.GetHandler(phaseType)
	return m.currentHandler.Enter(&m.Model)
}
