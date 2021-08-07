// structs, etc for the managing state of the TUI.
package frontend

import (
	"github.com/charmbracelet/bubbles/textinput"
	hn "hypermark/hackerNews"
	"os"
)

type ViewType int

const (
	startView ViewType = iota
	articleView
	promptView
	hyperpathsView
	editHyperpathView
)

// prompt and text input
type promptAndTextInput struct {
	textInput   textinput.Model
	prompt      string
	options     []string
	cursorIndex int
}

/*
	For CLI mode, hyperpaths[0] is what is always used/edited.
	This view should let the user:
	- Edit and add hyperpaths, with checking for whether or not the file exists
	  and allows for the user to create that file if they wish.
	- Swap the indices of hyperpaths
	- Return to the start screen
*/
type hyperpathsMenu struct {
	hyperpaths  []string
	cursorIndex int
}

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

type startMenu struct {
	choices     []string
	cursorIndex int
}

type model struct {
	clipboardOut bool
	outputPath   *os.File

	currentView        ViewType // Use this to choose which view to show.
	startMenu          startMenu
	articleMenu        articleMenu
	promptMenu         promptMenu
	hyperpathsMenu     hyperpathsMenu
	promptAndTextInput promptAndTextInput
}
