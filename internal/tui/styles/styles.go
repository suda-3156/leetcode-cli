package styles

import "github.com/charmbracelet/lipgloss"

var (
	// colors
	ColorMain = lipgloss.Color("#2bff00ff")
	ColorErr  = lipgloss.Color("#ff0000")
	ColorWarn = lipgloss.Color("#ffa200")

	// styles
	StyleSpinner       = lipgloss.NewStyle().Foreground(ColorMain)
	StyleActive        = lipgloss.NewStyle().Bold(true)
	StyleActionHeader  = lipgloss.NewStyle().Bold(true).Padding(0, 1).Foreground(ColorMain)
	StyleErrorHeader   = lipgloss.NewStyle().Bold(true).Padding(0, 1).Foreground(ColorErr)
	StyleWarningHeader = lipgloss.NewStyle().Bold(true).Padding(0, 1).Foreground(ColorWarn)
)
