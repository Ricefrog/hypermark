package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color("201")).
		PaddingTop(2).
		PaddingLeft(4).
		Width(22)

	choiceStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("5"))

	killStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF0000"))

	HighlightedCrimson = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#dc143c"))
)

