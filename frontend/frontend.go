package frontend

import (
	"fmt"
	"log"
	"hypermark/utils"
	tea "github.com/charmbracelet/bubbletea"
	hn "hypermark/hackerNews"
	"os"
)

var initialModel = model{
	startMenu: startMenu{
		choices: []string{
			"View hackernews articles",
			"Manage hypermarks",
			"Edit hyperpaths",
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

func (m *model) initializeArticles() {
	m.articleMenu.articles = hn.ScrapeHN()
}

func (m *model) loadHyperpaths() {
	hp, err := utils.GetAllHyperpaths()
	if err != nil {
		log.Fatal(err)
	}
	m.hyperpathsMenu.hyperpaths = hp
}

// Remove previous state
func (m *model) Wipe() {
	m.articleMenu = articleMenu{
		selected: make(map[int]struct{}),
	}
	m.promptMenu = promptMenu{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.currentView {
	case startView:
		return updateStartMenu(m, msg)
	case articleView:
		return updateArticleMenu(m, msg)
	case promptView:
		return updatePromptMenu(m, msg)
	case hyperpathsView:
		return updateHyperpathsMenu(m, msg)
	}
	return updateStartMenu(m, msg)
}

func (m model) View() string {
	switch m.currentView {
	case startView:
		return startMenuView(m)
	case articleView:
		return articleMenuView(m)
	case promptView:
		return promptMenuView(m)
	case hyperpathsView:
		return hyperpathsMenuView(m)
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
