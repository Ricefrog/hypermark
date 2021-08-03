package frontend

import (
	"fmt"
	"strconv"
	"hypermark/frontend/styles"
	"github.com/charmbracelet/lipgloss"
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
				m.currentView = 1
				initializeArticles(&m)
			}
		}
	}
	return m, nil
}

func startMenuView(m model) string {
	state := m.startMenu
	var outstr string
	if m.clipboardOut {
		outstr = "system clipboard"
	} else {
		outstr = m.outputPath.Name()
	}

	s := "\nhypermark\n\n"
	s += fmt.Sprintf("output path: %s\n\n", outstr)
	for i, choice := range state.choices {
		if i == state.cursorIndex {
			s += "-> "
		}
		s += choice+"\n"
	}
	return s
}

func updateArticleMenu(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	state := &m.articleMenu

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
			if state.cursorIndex < len(state.articles) {
				state.cursorIndex++
			}
		case "enter":
			if state.cursorIndex < len(state.articles)-1 {
				_, ok := state.selected[state.cursorIndex]
				if ok {
					delete(state.selected, state.cursorIndex)
				} else {
					state.selected[state.cursorIndex] = struct{}{}
				}
			}
		}
	}
	return m, nil
}

func articleMenuView(m model) string {
	state := m.articleMenu

	s := styles.HeaderStyle.Render("Top 30 on HackerNews")
	s += "\n"
	for i, article := range state.articles {
		cursor := ""
		if i == state.cursorIndex {
			cursor = "-> "
		}

		// Add highlighted style if article has been selected.
		style := lipgloss.NewStyle()
		if _, ok := state.selected[i]; ok {
			style = styles.HighlightedCrimson
		}

		number := strconv.Itoa(i+1)
		title := style.Render(article.Title)
		line := fmt.Sprintf("%s%s. %s\n", cursor, number, title)

		s += line
	}

	cursor := ""
	proceed := "Proceed?"
	if state.cursorIndex == len(state.articles) {
		cursor = "-> "
		proceed = styles.HRender(styles.ProtonPurple, proceed)
	}

	selected := strconv.Itoa(len(state.selected))
	selected = styles.HighlightedBlue.Render(selected)

	s += fmt.Sprintf("\n%s%s articles selected. %s", cursor, selected, proceed)
	return s
}
