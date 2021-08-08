package frontend

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	hn "hypermark/hackerNews"
	"hypermark/utils"
	"log"
	"os"
)

var initialModel = model{
	startMenu: startMenu{
		choices: []string{
			"View hackernews articles",
			"Manage bytemarks",
			"Edit hyperpaths",
		},
	},
	articleMenu: articleMenu{
		selected: make(map[int]struct{}),
	},
}

func ClearScreen() {
	for i := 0; i < 100; i++ {
		fmt.Println()
	}
}

func SetOutputVars(
	outputPath *os.File,
	tail []string,
	overwriteFile bool,
	writeToStdout bool,
	clipboardOut bool,
) {
	state := &initialModel.outputVars

	state.outputPath = outputPath
	state.tail = tail
	state.writeToStdout = writeToStdout
	state.clipboardOut = clipboardOut
}

func (m *model) syncOutputVars() error {
	state := &m.outputVars

	outputPath, err := utils.ChooseOutputPath(
		state.tail,
		state.overwriteFile,
		state.writeToStdout,
		state.clipboardOut,
	)
	state.outputPath = outputPath
	return err
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

func (m *model) initPromptAndTextInput(
	placeholder, prompt, footer string,
) {
	ti := textinput.NewModel()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 60

	m.promptAndTextInput.textInput = ti
	m.promptAndTextInput.prompt = prompt
	m.promptAndTextInput.footer = footer
}

func (m *model) setPrompt(prompt string, options []string) {
	m.promptMenu.prompt = prompt
	m.promptMenu.options = options
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
	case articlesAddedView:
		return updateArticlesAdded(m, msg)
	case hyperpathsView:
		return updateHyperpathsMenu(m, msg)
	case editHPView:
		return updateEditHyperpath(m, msg)
	case addHPView:
		return updateEditHyperpath(m, msg)
	case createFileView:
		return updateCreateFile(m, msg)
	case invalidFilepathView:
		return updateInvalidFilepath(m, msg)
	}
	return updateStartMenu(m, msg)
}

func (m model) View() string {
	switch m.currentView {
	case startView:
		return startMenuView(m)
	case articleView:
		return articleMenuView(m)
	case articlesAddedView:
		return promptMenuView(m)
	case hyperpathsView:
		return hyperpathsMenuView(m)
	case editHPView:
		return promptAndTextInputView(m)
	case addHPView:
		return promptAndTextInputView(m)
	case createFileView:
		return promptMenuView(m)
	case invalidFilepathView:
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
