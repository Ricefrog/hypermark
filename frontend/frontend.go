package frontend

import (
	"fmt"
	"os"
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

func (m model) View() string {
	switch m.currentView {
	case 0:
		return startMenuView(m)
	case 1:
		return articleMenuView(m)
	}
	return startMenuView(m)
}

func Start() {
	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Println("error: %v", err)
		os.Exit(1)
	}
}
