package frontend

import (
	"fmt"
	"log"
	"strings"
	"strconv"
	"hypermark/frontend/styles"
	"hypermark/frontend/templates"
	"hypermark/utils"
	tea "github.com/charmbracelet/bubbletea"
)

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
					"%s %s",
					styles.HRender(styles.Crimson, "Editing"),
					styles.MakeHyperpathString(state.cursorIndex),
				)

				submit := styles.CommandInfo("Submit", "enter")
				back := styles.CommandInfo("Go back", "esc")
				footer := fmt.Sprintf("%s | %s", submit, back)

				state.editHyperpath.index = state.cursorIndex
				m.initPromptAndTextInput(placeholder, prompt, footer)
				m.currentView = editHPView
			case "d":
				if len(state.hyperpaths) == 1 { break }
				m.promptMenu.prompt = fmt.Sprintf("Delete '%s'?",
					state.hyperpaths[state.cursorIndex],
				)
				m.promptMenu.options = []string{"Yes", "Cancel"}
				m.currentView = deleteHyperpathView
			case "m":
				state.moveMode = true
			case "n":
				placeholder := fmt.Sprintf(
					"hyperpath[%d]", len(state.hyperpaths),
				)
				prompt := fmt.Sprintf("%s %s",
					styles.HRender(styles.Crimson, "Creating"),
					styles.MakeHyperpathString(len(state.hyperpaths)),
				)

				submit := styles.CommandInfo("Submit", "enter")
				back := styles.CommandInfo("Go back", "esc")
				footer := fmt.Sprintf("%s | %s", submit, back)

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
		del = fmt.Sprintf(" | %s", styles.CommandInfo("Delete", "d"))
	}
	if !state.moveMode {
		move = fmt.Sprintf(" | %s", styles.CommandInfo("Move", "m"))
	} else {
		move = fmt.Sprintf(" | %s", styles.CommandInfo("Drop", "m"))
	}

	s += fmt.Sprintf("\n%s: %s%s%s\n\n",
		styles.MakeHyperpathString(state.cursorIndex),
		styles.CommandInfo("Edit", "e"),
		del,
		move,
	)

	for i, hyperpath := range state.hyperpaths {
		cursor := ""
		num := strconv.Itoa(i)
		colon := ":"
		if state.cursorIndex == i {
			cursor = templates.Cursor()
			if state.moveMode {
				hyperpath = styles.HRender(styles.OrangeRed, hyperpath)
			} else {
				hyperpath = styles.HRender(styles.ProtonPurple, hyperpath)
			}
			num = styles.HRender(styles.Crimson, num)
			colon = styles.HRender(styles.OrangeRed, colon)
		}
		s += fmt.Sprintf("%s%s%s %s\n", cursor, num, colon, hyperpath)
	}

	s += fmt.Sprintf("\n%s\n%s\n",
		styles.CommandInfo("Add new hyperpath", "n"),
		styles.CommandInfo("Go back", "esc"),
	)
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
			m.wipePromptMenu()
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
				m.loadHyperpaths() // Reload data from hyperpaths file.
			}
			m.wipePromptMenu()
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
			m.wipePromptMenu()
			m.currentView = hyperpathsView
		}
	}
	return m, nil
}

func updateDeleteHyperpath(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	stateA := &m.hyperpathsMenu
	stateB := &m.promptMenu

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl-c", "q":
			return m, tea.Quit
		case "up", "k":
			if stateB.cursorIndex > 0 {
				stateB.cursorIndex--
			}
		case "down", "j":
			if stateB.cursorIndex < len(stateB.options)-1 {
				stateB.cursorIndex++
			}
		case "enter", "esc":
			if stateB.cursorIndex == 0 {
				stateA.hyperpaths = utils.DeleteElement(
					stateA.hyperpaths,
					stateA.cursorIndex,
				)
				if stateA.cursorIndex > 0 {
					stateA.cursorIndex--
				}
				if err := utils.WriteHyperpaths(stateA.hyperpaths); err != nil {
					log.Fatal(err)
				}

				m.loadHyperpaths()
			}

			m.wipePromptMenu()
			m.currentView = hyperpathsView
		}
	}

	return m, nil
}
