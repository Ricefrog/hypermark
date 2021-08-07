package frontend

import (
	"fmt"
	"log"
	"strconv"
	"hypermark/frontend/styles"
	"hypermark/frontend/templates"
	"hypermark/utils"
	hn "hypermark/hackerNews"
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
				m.currentView = articleView
				m.initializeArticles()
			case 2:
				m.currentView = hyperpathsView
				m.loadHyperpaths()
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
			s += templates.Cursor()
		}
		s += choice+"\n"
	}
	return s
}

func updateArticleMenu(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	state := &m.articleMenu
	var from int
	var to int
	if state.pageIndex == 0 {
		from = 0
		to = 15
	} else {
		from = 15
		to = 30
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		// Change cursor-selected article
		case "up", "k":
			if state.cursorIndex > from {
				state.cursorIndex--
			}
		case "down", "j":
			if state.cursorIndex < to {
				state.cursorIndex++
			}
		// Change page
		case "right", "l":
			if state.pageIndex == 0 {
				state.pageIndex++
				state.cursorIndex = 16
			}
		case "left", "h":
			if state.pageIndex == 1 {
				state.pageIndex--
				state.cursorIndex = 0
			}
		case "enter", " ":
			if state.cursorIndex < to-1 {
				_, ok := state.selected[state.cursorIndex]
				if ok {
					delete(state.selected, state.cursorIndex)
				} else {
					state.selected[state.cursorIndex] = struct{}{}
				}
			} else { // If "enter" is pressed on the checkout prompt
				articles := make([]hn.HNArticle, 0)
				for i, _ := range state.selected {
					// Append each HNArticle to the list.
					articles = append(articles, state.articles[i])
				}

				output := utils.ArticlesToTable(articles)
				writtenTo, err := utils.Write(
					m.outputPath, output, m.clipboardOut,
				)
				if err != nil {
					log.Fatal(err)
				}

				// Set up the prompt screen
				m.currentView = promptView
				m.promptMenu.options = []string{"Continue", "Quit"}
				m.promptMenu.prompt = fmt.Sprintf(
					"%d articles written to %s.\n",
					len(articles),
					writtenTo,
				)
			}
		}
	}
	return m, nil
}

func articleMenuView(m model) string {
	state := m.articleMenu

	s := styles.HeaderStyle.Render("Top 30 on HackerNews")
	s += "\n"

	var from int
	var to int
	var page string
	if state.pageIndex == 0 {
		from = 0
		to = 15
	} else {
		from = 15
		to = 30
	}
	instr := "arrow keys/hjkl to navigate"
	page = fmt.Sprintf("\nArticles %d-%d (%s):\n", from+1, to, instr)
	s += page

	for i := from; i < to; i++ {
		article := state.articles[i]

		cursor := ""
		number := strconv.Itoa(i+1)
		if i == state.cursorIndex {
			cursor = templates.Cursor()
			number = styles.HRender(styles.JustBlue, number)
		}

		// Add highlighted style if article has been selected.
		style := lipgloss.NewStyle()
		if _, ok := state.selected[i]; ok {
			style = styles.HighlightedCrimson
		}
		title := style.Render(article.Title)
		line := fmt.Sprintf("%s%s. %s\n", cursor, number, title)

		s += line
	}

	cursor := ""
	proceed := "Proceed?"
	if state.cursorIndex == to {
		cursor = templates.Cursor()
		proceed = styles.HRender(styles.ProtonPurple, proceed)
	}

	selected := strconv.Itoa(len(state.selected))
	selected = styles.HighlightedBlue.Render(selected)

	s += fmt.Sprintf("\n%s%s articles selected. %s", cursor, selected, proceed)
	return s
}

func updatePromptMenu(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	state := &m.promptMenu

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
			if state.cursorIndex < len(state.options)-1 {
				state.cursorIndex++
			}
		case "enter", " ":
			if state.cursorIndex == 0 {
				m.currentView = startView
				m.Wipe()
			} else {
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func promptMenuView(m model) string {
	state := m.promptMenu

	var s string
	s += templates.Prompt(state.prompt)
	for i, option := range state.options {
		cursor := ""
		if state.cursorIndex == i {
			cursor = templates.Cursor()
			option = styles.HRender(styles.ProtonPurple, option)
		}
		s += fmt.Sprintf("%s%s\n", cursor, option)
	}
	return s
}

func updateHyperpathsMenu(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	state := &m.hyperpathsMenu

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
			if state.cursorIndex < len(state.hyperpaths)-1 {
				state.cursorIndex++
			}
		}
	}

	return m, nil
}

func hyperpathsMenuView(m model) string {
	state := m.hyperpathsMenu

	var s string

	var del string
	if len(state.hyperpaths) > 1{
		del = "| Delete (d)"
	}
	s += fmt.Sprintf("\nhyperpaths[%d]: Edit (e) %s\n\n",
		state.cursorIndex,
		del,
	)

	for i, hyperpath := range state.hyperpaths {
		cursor := ""
		if state.cursorIndex == i {
			cursor = templates.Cursor()
			hyperpath = styles.HRender(styles.ProtonPurple, hyperpath)
		}
		s += fmt.Sprintf("%s%d: %s\n", cursor, i, hyperpath)
	}

	s += "\nAdd new hyperpath (n)\n"
	return s
}

func updateEditHyperpath(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	state := &m.promptAndTextInput

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if state.cursorIndex > 0 {
				state.cursorIndex--
			}
		case "down", "j":
			if state.cursorIndex < len(state.hyperpaths)-1 {
				state.cursorIndex++
			}
		}
	}

	return m, nil
}
