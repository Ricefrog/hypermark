package frontend

import (
	"fmt"
	"os"
	"strconv"
	"hypermark/frontend/styles"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	hn "hypermark/hackerNews"
)

type articleMenu struct {
	articles []hn.HNArticle
	selected map[int]struct{}

	cursorIndex int
}

func initializeArticles(m *model) {
	m.articleMenu.articles = hn.ScrapeHN()
}

type startMenu struct {
	choices []string
	cursorIndex int
}

type model struct {
	clipboardOut bool
	outputPath *os.File

	currentView int
	startMenu startMenu
	articleMenu articleMenu
}

var initialModel = model{
	startMenu: startMenu{
		choices: []string{
			"View hackernews articles",
			"Manage hypermarks",
			"Change hyperpaths",
		},
	},
	articleMenu: articleMenu{
		selected: make(map[int]struct{}),
	},
}

func SetOutputPath(outputPath *os.File, clipboardOut bool) {
	if clipboardOut {
		initialModel.clipboardOut = true
	} else {
		initialModel.outputPath = outputPath
	}
}

func ClearScreen() {
	for i := 0; i < 100; i++ {
		fmt.Println()
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.currentView {
	case 0:
		return updateStartMenu(m, msg)
	case 1:
		return updateArticleMenu(m, msg)
	}
	return updateStartMenu(m, msg)
}

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
			if state.cursorIndex < len(state.articles)-1 {
				state.cursorIndex++
			}
		case "enter":
			_, ok := state.selected[state.cursorIndex]
			if ok {
				delete(state.selected, state.cursorIndex)
			} else {
				state.selected[state.cursorIndex] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.currentView {
	case 0:
		return startMenuView(m)
	case 1:
		return articleMenuView(m)
	}
	return startMenuView(m)
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

func articleMenuView(m model) string {
	state := m.articleMenu

	s := "\nTop 30 on HackerNews\n\n"
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
	return s
}

func Start() {
	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Println("error: %v", err)
		os.Exit(1)
	}
}
