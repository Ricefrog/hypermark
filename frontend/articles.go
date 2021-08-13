package frontend

import (
	"fmt"
	"log"
	"strconv"
	"hypermark/frontend/styles"
	"hypermark/frontend/templates"
	"hypermark/utils"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
)

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
				state.cursorIndex = 15
			}
		case "left", "h":
			if state.pageIndex == 1 {
				state.pageIndex--
				state.cursorIndex = 0
			}
		case "enter", " ":
			if state.cursorIndex < to {
				_, ok := state.selected[state.cursorIndex]
				if ok {
					delete(state.selected, state.cursorIndex)
				} else {
					state.selected[state.cursorIndex] = struct{}{}
				}
			} else { // If "enter" is pressed on the checkout prompt
				articles := make([]utils.Bytemark, 0)
				for i, _ := range state.selected {
					// Append each HNArticle to the list.
					articles = append(articles, state.articles[i])
				}

				output := utils.BytemarksToTables(articles)
				writtenTo, err := utils.Write(
					m.outputVars.outputPath, output, m.outputVars.clipboardOut,
				)
				if err != nil {
					log.Fatal(err)
				}

				// Set up the 'articles added' prompt screen
				m.currentView = articlesAddedView
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

func updateArticlesAdded(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
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
				ClearScreen()
				m.currentView = startView
				m.Wipe()
			} else {
				return m, tea.Quit
			}
		}
	}

	return m, nil
}
