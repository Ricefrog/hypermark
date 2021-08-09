package frontend

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"os"
	"hypermark/frontend/styles"
	"hypermark/frontend/templates"
	"hypermark/utils"
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
				m.initializeArticles()
				m.currentView = articleView
			case 1:
				m.loadHyperpaths()
				m.currentView = bytemarksMainView
			case 2:
				m.loadHyperpaths()
				m.currentView = hyperpathsView
			}
		}
	}
	return m, nil
}

func startMenuView(m model) string {
	state := m.startMenu
	var outstr string
	if m.outputVars.clipboardOut {
		outstr = "system clipboard"
	} else {
		outstr = m.outputVars.outputPath.Name()
	}

	s := "\nhypermark\n\n"
	s += fmt.Sprintf("Writing to -> %s\n\n", outstr)
	for i, choice := range state.choices {
		if i == state.cursorIndex {
			s += templates.Cursor()
			choice = styles.HRender(styles.ProtonPurple, choice)
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
		if !state.moveMode {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "esc":
				m.currentView = startView
				return m, nil
			case "up", "k":
				if state.cursorIndex > 0 {
					state.cursorIndex--
				}
			case "down", "j":
				if state.cursorIndex < len(state.hyperpaths)-1 {
					state.cursorIndex++
				}
			case "e":
				selectedHP := state.hyperpaths[state.cursorIndex]
				placeholder := selectedHP
				prompt := fmt.Sprintf(
					"Editing hyperpath[%d]",
					state.cursorIndex,
				)
				footer := "Submit (enter) | Go back (esc)"

				state.editHyperpath.index = state.cursorIndex
				m.initPromptAndTextInput(placeholder, prompt, footer)
				m.currentView = editHPView
			case "d":
				if len(state.hyperpaths) == 1 { break }
				state.hyperpaths = utils.DeleteElement(
					state.hyperpaths,
					state.cursorIndex,
				)
				if state.cursorIndex > 0 {
					state.cursorIndex--
				}
				err := utils.WriteHyperpaths(state.hyperpaths)
				if err != nil {
					log.Fatal(err)
				}
				m.loadHyperpaths()
			case "m":
				state.moveMode = true
			case "n":
				placeholder := fmt.Sprintf(
					"hyperpath[%d]", len(state.hyperpaths),
				)
				prompt := fmt.Sprintf("Creating %s", placeholder)
				footer := "Submit (enter) | Go back (esc)"

				state.editHyperpath.index = len(state.hyperpaths)
				m.initPromptAndTextInput(placeholder, prompt, footer)
				m.currentView = addHPView
			}
		} else {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "m", "enter", "esc":
				if err := utils.WriteHyperpaths(state.hyperpaths); err != nil {
					log.Fatal(err)
				}
				if err := m.syncOutputVars(); err != nil {
					log.Fatal(err)
				}
				state.moveMode = false
			case "up", "k":
				if state.cursorIndex > 0 {
					state.hyperpaths = utils.SwapElements(
						state.hyperpaths,
						state.cursorIndex - 1,
						state.cursorIndex,
					)
					state.cursorIndex--
				}
			case "down", "j":
				if state.cursorIndex < len(state.hyperpaths)-1 {
					state.hyperpaths = utils.SwapElements(
						state.hyperpaths,
						state.cursorIndex + 1,
						state.cursorIndex,
					)
					state.cursorIndex++
				}
			}
		}
	}

	return m, nil
}

func hyperpathsMenuView(m model) string {
	state := m.hyperpathsMenu

	var s string
	var del string
	var move string
	if len(state.hyperpaths) > 1 && !state.moveMode {
		del = " | Delete (d)"
	}
	if !state.moveMode {
		move = " | Move (m)"
	} else {
		move = " | Drop (m)"
	}

	s += fmt.Sprintf("\nhyperpaths[%d]: Edit (e)%s%s\n\n",
		state.cursorIndex,
		del,
		move,
	)

	for i, hyperpath := range state.hyperpaths {
		cursor := ""
		if state.cursorIndex == i {
			cursor = templates.Cursor()
			if state.moveMode {
				hyperpath = styles.HRender(styles.OrangeRed, hyperpath)
			} else {
				hyperpath = styles.HRender(styles.ProtonPurple, hyperpath)
			}
		}
		s += fmt.Sprintf("%s%d: %s\n", cursor, i, hyperpath)
	}

	s += "\nAdd new hyperpath (n)\nGo back (esc)\n"
	return s
}

func updateEditHyperpath(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	stateA := &m.promptAndTextInput
	stateB := &m.hyperpathsMenu.editHyperpath
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.currentView = hyperpathsView
			return m, nil
		case "enter":
			newHyperpath := stateA.textInput.Value()
			if strings.Contains(newHyperpath, "~") {
				newHyperpath = utils.ExpandTilde(newHyperpath)
			}
			stateB.newHyperpath = newHyperpath

			written, valid := utils.EditNthHyperpath(newHyperpath, stateB.index)
			if written && valid {
				if err := m.syncOutputVars(); err != nil {
					log.Fatal(err)
				}
				m.loadHyperpaths()
				m.currentView = hyperpathsView
			} else if valid {
				// Path is valid but file does not exist.
				prompt := fmt.Sprintf("%s does not exist.", newHyperpath)
				options := []string{"Create file", "Go back"}
				m.setPrompt(prompt, options)
				m.currentView = createFileView
			} else {
				// Path is completely invalid.
				prompt := fmt.Sprintf(
					"%s is not a valid filepath.", newHyperpath,
				)
				options := []string{"Go back"}
				m.setPrompt(prompt, options)
				m.currentView = invalidFilepathView
			}
		}
	}

	stateA.textInput, cmd = stateA.textInput.Update(msg)
	return m, cmd
}

func promptAndTextInputView(m model) string {
	state := m.promptAndTextInput

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		state.prompt,
		state.textInput.View(),
		state.footer,
	)
}

func updateCreateFile(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	state := &m.promptMenu
	index := m.hyperpathsMenu.editHyperpath.index
	newHyperpath := m.hyperpathsMenu.editHyperpath.newHyperpath

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			m.currentView = startView
			return m, nil
		case "up", "k":
			if state.cursorIndex > 0 {
				state.cursorIndex--
			}
		case "down", "j":
			if state.cursorIndex < len(state.options)-1 {
				state.cursorIndex++
			}
		case "enter":
			if state.cursorIndex == 0 {
				// Create the file.
				if _, err := utils.CreateFile(newHyperpath); err != nil {
					log.Fatal(err)
				}
				written, valid := utils.EditNthHyperpath(newHyperpath, index)
				if !written || !valid {
					var wrongRet string
					if !written && !valid {
						wrongRet += "written or valid"
					} else if !written {
						wrongRet += "written"
					} else {
						wrongRet += "valid"
					}
					wrongRet = fmt.Sprintf("'%s'", wrongRet)

					message := fmt.Sprintf(
						"hyperpath did not return %s.\n"+
						"newHyperpath: %s\nindex: %d",
						wrongRet,
						newHyperpath,
						index,
					)
					log.Fatal(message)
				}
				if err := m.syncOutputVars(); err != nil {
					log.Fatal(err)
				}
				m.loadHyperpaths()
			}
			m.currentView = hyperpathsView
		}
	}

	return m, nil
}

func updateInvalidFilepath(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", "esc":
			m.currentView = hyperpathsView
		}
	}
	return m, nil
}

func updateBytemarksMenu(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	state := &m.hyperpathsMenu

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl-c", "q":
			return m, tea.Quit
		case "up", "k":
			if state.cursorIndex > 0 {
				state.cursorIndex--
			}
		case "down", "j":
			if state.cursorIndex < len(state.hyperpaths)-1 {
				state.cursorIndex++
			}
		case "enter":
			// load bytemarks of the selected hyperpath into state.
			var err error
			file, err := os.OpenFile(
				state.hyperpaths[state.cursorIndex],
				os.O_RDWR,
				0666,
			)
			if err != nil {
				log.Fatal(err)
			}
			m.bytemarksManager.bytemarks, err = utils.FileToBytemarks(file)
			if err != nil {
				log.Fatal(err)
			}
			m.bytemarksManager.hyperpath = state.hyperpaths[state.cursorIndex]
			m.currentView = byteManagerView
		case "esc":
			m.currentView = startView
		}
	}
	return m, nil
}

func bytemarksMenuView(m model) string {
	state := m.hyperpathsMenu

	var s string
	s += fmt.Sprintf("\nhyperpaths[%d]: Manage bytemarks (enter)\n\n", state.cursorIndex)

	for i, hyperpath := range state.hyperpaths {
		cursor := ""
		if state.cursorIndex == i {
			cursor = templates.Cursor()
			hyperpath = styles.HRender(styles.ProtonPurple, hyperpath)
		}
		s += fmt.Sprintf("%s%d: %s\n", cursor, i, hyperpath)
	}
	return s
}

func updateBytemarksManager(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	state := &m.bytemarksManager

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl-c", "q":
			return m, tea.Quit
		case "s":
			if state.moveMode {state.moveMode = false}
			// Switch to a prompt for saving.
			m.promptMenu.prompt = fmt.Sprintf(
				"Save changes to %s?",
				state.hyperpath,
			)
			m.promptMenu.options = []string{"Save", "Cancel"}
			m.currentView = saveChangesView
		case "m":
			state.moveMode = !state.moveMode
		case "up", "k":
			if state.cursorIndex > 0 {
				if state.moveMode {
					state.bytemarks = utils.SwapBytemarks(
						state.bytemarks,
						state.cursorIndex,
						state.cursorIndex-1,
					)
				}
				state.cursorIndex--
			}
		case "down", "j":
			if state.cursorIndex < len(state.bytemarks)-1 {
				if state.moveMode {
					state.bytemarks = utils.SwapBytemarks(
						state.bytemarks,
						state.cursorIndex,
						state.cursorIndex+1,
					)
				}
				state.cursorIndex++
			}
		case "esc":
			m.currentView = bytemarksMainView
		}
	}
	return m, nil
}

func bytemarksManagerView(m model) string {
	state := m.bytemarksManager

	if len(state.bytemarks) == 0 {
		return "No bytemarks to display.\nGo back (esc)"
	}

	var s string
	var move string
	if !state.moveMode {
		move = " | Move (m)"
	} else {
		move = " | Drop (m)"
	}

	s += fmt.Sprintf("\nSave (s) | Delete (d)%s\n\n", move)

	for i, bytemark := range state.bytemarks {
		title := bytemark.Title
		cursor := ""
		if state.cursorIndex == i {
			cursor = templates.Cursor()
			if state.moveMode {
				title = styles.HRender(styles.OrangeRed, title)
			} else {
				title = styles.HRender(styles.ProtonPurple, title)
			}
		}
		s += fmt.Sprintf("%s%d: %s\n", cursor, i, title)
	}

	s += "\nCreate bytemark using system clipboard (n)\nGo back (esc)\n"
	return s
}

func updateSaveChanges(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	stateA := &m.promptMenu
	stateB := &m.bytemarksManager

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			m.currentView = byteManagerView
		case "up", "k":
			if stateA.cursorIndex > 0 {
				stateA.cursorIndex--
			}
		case "down", "j":
			if stateA.cursorIndex < len(stateA.options)-1 {
				stateA.cursorIndex++
			}
		case "enter":
			if stateA.cursorIndex == 0 {
				output := utils.BytemarksToTables(stateB.bytemarks)
				if err := os.Remove(stateB.hyperpath); err != nil {
					log.Fatal(err)
				}
				selectedFile, err := os.OpenFile(
					stateB.hyperpath,
					os.O_CREATE|os.O_RDWR,
					0666,
				)
				if err != nil {
					log.Fatal(err)
				}
				_, err = utils.Write(
					selectedFile,
					output,
					m.outputVars.clipboardOut,
				)
				if err != nil {
					log.Fatal(err)
				}
				stateB.bytemarks, err = utils.FileToBytemarks(selectedFile)
				if err != nil {
					log.Fatal(err)
				}
			}
			m.currentView = byteManagerView
		}
	}
	return m, nil
}
