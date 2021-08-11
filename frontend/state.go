// structs, etc for the managing state of the TUI.
package frontend

import (
	"github.com/charmbracelet/bubbles/textinput"
	"hypermark/utils"
	"os"
)

type ViewType int

const (
	startView ViewType = iota
	articleView
	articlesAddedView
	bytemarksMainView
	byteManagerView
	deleteBytemarkView
	sendBytemarkView
	saveChangesView
	hyperpathsView
	editHPView
	addHPView
	createFileView
	invalidFilepathView
)

// Generic prompt and text input
type promptAndTextInput struct {
	textInput   textinput.Model
	prompt      string
	footer      string
	cursorIndex int
}

type editHyperpath struct {
	newHyperpath string
	index        int
}

type hyperpathsMenu struct {
	hyperpaths    []string
	editHyperpath editHyperpath
	moveMode      bool
	cursorIndex   int
}

type bytemarksManager struct {
	bytemarks       []utils.Bytemark
	moveMode        bool
	cursorIndex     int
	hyperpath       string
	otherHyperpaths []string
}

// Generic prompt menu
type promptMenu struct {
	prompt      string
	options     []string
	cursorIndex int
}

type articleMenu struct {
	articles    []utils.Bytemark
	selected    map[int]struct{}
	cursorIndex int
	pageIndex   int
}

type startMenu struct {
	choices     []string
	cursorIndex int
}

type outputVars struct {
	tail          []string
	overwriteFile bool
	writeToStdout bool
	clipboardOut  bool
	outputPath    *os.File
}

type model struct {
	outputVars outputVars

	currentView        ViewType // Use this to choose which view to show.
	startMenu          startMenu
	articleMenu        articleMenu
	promptMenu         promptMenu
	hyperpathsMenu     hyperpathsMenu
	bytemarksManager   bytemarksManager
	promptAndTextInput promptAndTextInput
}
