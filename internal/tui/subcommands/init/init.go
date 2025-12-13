package init

import tea "github.com/charmbracelet/bubbletea"

// Init implements tea.Model for the confirmation prompt
func (m *Model) Init() tea.Cmd {
	return nil
}

// Run starts the confirmation prompt and returns true if the user confirmed overwriting.
func Run(path string) (bool, error) {
	initialModel := &Model{
		path:   path,
		choice: 1, // Default to No
	}

	p := tea.NewProgram(initialModel)

	final, err := p.Run()
	if err != nil {
		return false, err
	}

	fm := final.(*Model)
	return fm.confirmed, nil
}
