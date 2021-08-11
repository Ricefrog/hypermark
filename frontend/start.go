package frontend

import (
	"fmt"
	"hypermark/frontend/styles"
	"hypermark/frontend/templates"
	tea "github.com/charmbracelet/bubbletea"
)

func updateStartMenu(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	state := &m.startMenu

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if state.cursorIndex > 0 {
				state.cursorIndex--
			}
		case "down", "j":
			if state.cursorIndex < len(state.choices)-1 {
				state.cursorIndex++
			}
		case "enter":
			switch state.cursorIndex {
			case 0:
				m.initializeArticles()
				m.currentView = articleView
			case 1:
				m.loadHyperpaths()
				m.currentView = bytemarksMainView
			case 2:
				m.loadHyperpaths()
				m.currentView = hyperpathsView
			}
		}
	}
	return m, nil
}

func startMenuView(m model) string {
	state := m.startMenu
	var outstr string
	if m.outputVars.clipboardOut {
		outstr = "system clipboard"
	} else {
		outstr = m.outputVars.outputPath.Name()
	}

	s := "\nhypermark\n\n"
	s += fmt.Sprintf("Writing to -> %s\n\n", outstr)
	for i, choice := range state.choices {
		if i == state.cursorIndex {
			s += templates.Cursor()
			choice = styles.HRender(styles.ProtonPurple, choice)
		}
		s += choice+"\n"
	}
	return s
}
