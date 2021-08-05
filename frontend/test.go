package frontend

import (
	"fmt"
	"os"
	"hypermark/frontend/templates"
	tea "github.com/charmbracelet/bubbletea"
)

type testModel struct {
	output string
}

func (m testModel) Init() tea.Cmd {
	return nil
}

func (m testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m testModel) View() string {
	var s string
	prompt := "5 articles written to /home/severian/terminus_est/tester.md"
	s += templates.Prompt(prompt)
	s += "Continue\n"
	s += "Quit\n"
	return s
}

func Test() {
	p := tea.NewProgram(testModel{})
	if err := p.Start(); err != nil {
		fmt.Println("error: %v", err)
		os.Exit(1)
	}
}
