package frontend

import (
	"fmt"
	"os"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	messages []string
	currentMessage int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.currentMessage > 0 {
				m.currentMessage--
			}
		case "down", "j":
			if m.currentMessage < len(m.messages)-1 {
				m.currentMessage++
			}
		}
	}
	return m, nil
}

var headerStyle = lipgloss.NewStyle().
					Bold(true).
					Background(lipgloss.Color("201")).
					PaddingTop(2).
					PaddingLeft(4).
					Width(22)

var choiceStyle = lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("5"))
var killStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FF0000"))

func (m model) View() string {
	// Header
	s := headerStyle.Render("Welcome to the Based Store!\n\n")
	s += "\nWhat would you like to buy?\n\n"
	s += fmt.Sprintf("Your choice -> %s\n",
					  choiceStyle.Render(m.messages[m.currentMessage]))
	s += "\nPress q to join the "
	s += killStyle.Render("1%\n")
	return s
}

func Start() {
	p := tea.NewProgram(model{
		messages: []string{"grapes", "carrots", "anime pfp"},
	})
	if err := p.Start(); err != nil {
		fmt.Println("error: %v", err)
		os.Exit(1)
	}
}
