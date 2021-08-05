package frontend

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	hn "hypermark/hackerNews"
	"os"
)

type ViewType int

const (
	startView ViewType = iota
	articleView
	promptView
)

type promptMenu struct {
	prompt      string
	options     []string
	cursorIndex int
}

type articleMenu struct {
	articles    []hn.HNArticle
	selected    map[int]struct{}
	cursorIndex int
	pageIndex   int
}

func initializeArticles(m *model) {
	m.articleMenu.articles = hn.ScrapeHN()
}

type startMenu struct {
	choices     []string
	cursorIndex int
}

type model struct {
	clipboardOut bool
	outputPath   *os.File

	currentView ViewType    // Use this to choose which view to show.
	startMenu   startMenu
	articleMenu articleMenu
	promptMenu  promptMenu
}

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
