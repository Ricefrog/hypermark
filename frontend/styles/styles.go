package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ProtonPurple = "#8A2BE2"
	AquaMenthe = "#7FFFD4"
	CosmicLatte = "#FFF9E3"
	OrangeRed = "#FF4500"
	JustBlue = "#0000FF"
	Crimson = "#DC143C"
	MatureCrimson = "#9a031e"
	Burgundy = "#5f0f40"
	PalestBlue = "#C6E2FF"

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
		Foreground(lipgloss.Color("#DC143C"))

	HighlightedBlue = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#0C71E0"))

	HighlightedHotPink = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF007F"))

	HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(Crimson)).
		Background(lipgloss.Color(CosmicLatte)).
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color(ProtonPurple)).
		Padding(1)

	PromptStyle = lipgloss.NewStyle().
		Bold(true).
		Padding(2).
		Foreground(lipgloss.Color(MatureCrimson)).
		Background(lipgloss.Color(CosmicLatte)).
		BorderForeground(lipgloss.Color(JustBlue)).
		BorderStyle(lipgloss.RoundedBorder())

	CursorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(AquaMenthe))
)

func Highlighted(color string) lipgloss.Style {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(color))
	return style
}

// render string, bolded and colored
func HRender(color, message string) string {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(color))
	return style.Render(message)
}
