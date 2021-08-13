package frontend

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"hypermark/frontend/styles"
	"hypermark/frontend/templates"
	"hypermark/utils"
	tea "github.com/charmbracelet/bubbletea"
)

// Same size array. For convenience when displaying.
func makeOtherHyperpaths(hyperpaths []string, index int) []string {
	culled := make([]string, len(hyperpaths))
	for i, hp := range hyperpaths {
		if i != index {culled[i] = hp}
	}
	return culled
}

func firstNonEmpty(arr []string) (string, int) {
	for i, s := range arr {
		if s != "" {
			return s, i
		}
	}
	return "something went wrong\n", -1
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
				errMsg := fmt.Sprintf("Error opening hyperpaths[%d]", state.cursorIndex)
				errMsg += err.Error()
				cErr := errors.New(errMsg)
				log.Fatal(cErr)
			}
			m.bytemarksManager.bytemarks, err = utils.FileToBytemarks(file)
			if err != nil {
				log.Fatal(err)
			}
			file.Close()

			m.bytemarksManager.hyperpath = state.hyperpaths[state.cursorIndex]
			m.bytemarksManager.otherHyperpaths = makeOtherHyperpaths(
				state.hyperpaths,
				state.cursorIndex,
			)
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
	s += fmt.Sprintf("\n%s: Manage bytemarks (enter)\n\n",
		styles.MakeHyperpathString(state.cursorIndex),
	)

	for i, hyperpath := range state.hyperpaths {
		cursor := ""
		num := strconv.Itoa(i)
		colon := ":"
		if state.cursorIndex == i {
			cursor = templates.Cursor()
			hyperpath = styles.HRender(styles.ProtonPurple, hyperpath)
			num = styles.HRender(styles.Crimson, num)
			colon = styles.HRender(styles.OrangeRed, colon)
		}
		s += fmt.Sprintf("%s%s%s %s\n", cursor, num, colon, hyperpath)
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
		case "t":
			_, sendIndex := firstNonEmpty(state.otherHyperpaths)
			m.promptMenu.prompt = fmt.Sprintf("Send bytemark to hyperpath[%d]", sendIndex)
			m.promptMenu.cursorIndex = sendIndex
			m.promptMenu.options = utils.Copy(state.otherHyperpaths)
			m.currentView = sendBytemarkView
		case "p":
			state.bytemarks = utils.InsertBytemark(
				state.bytemarks,
				state.bytemarks[state.cursorIndex],
				state.cursorIndex,
			)
		case "d":
			m.promptMenu.prompt = fmt.Sprintf("Are you sure?")
			m.promptMenu.options = []string{"Yes", "Cancel"}
			m.currentView = deleteBytemarkView
		case "m":
			state.moveMode = !state.moveMode
		case "n":
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
			state.cursorIndex = 0
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
	s += fmt.Sprintf("%s: %s\n",
		styles.HRender(styles.AquaMenthe, "bytemarks"),
		styles.StylePath(state.hyperpath),
	)
	var move string
	if !state.moveMode {
		move = " | Move (m)"
	} else {
		move = " | Drop (m)"
	}

	s += fmt.Sprintf("Save changes (s) | Duplicate (p) | Send to (t) | Delete (d)%s\n\n", move)

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
		s += fmt.Sprintf("%s%s\n", cursor, title)
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
			m.wipePromptMenu()
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
				selectedFile.Close()

				file, err := os.OpenFile(
					stateB.hyperpath,
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
				file.Close()
				if err := m.syncOutputVars(); err != nil {
					log.Fatal(err)
				}
			}
			m.wipePromptMenu()
			m.currentView = byteManagerView
		}
	}
	return m, nil
}

func updateDeleteBytemark(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	stateA := &m.promptMenu
	stateB := &m.bytemarksManager

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl-c", "q":
			return m, tea.Quit
		case "esc":
			m.wipePromptMenu()
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
				stateB.bytemarks = utils.DeleteBytemark(
					stateB.bytemarks,
					stateB.cursorIndex,
				)
				stateB.cursorIndex--
			}
			m.wipePromptMenu()
			m.currentView = byteManagerView
		}
	}
	return m, nil
}

func updateSendBytemark(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	stateA := &m.promptMenu
	stateB := &m.bytemarksManager

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl-c", "q":
			return m, tea.Quit
		case "esc":
			m.wipePromptMenu()
			m.currentView = byteManagerView
		case "up", "k":
			if stateA.cursorIndex > 0 {
				stateA.cursorIndex--
				if stateA.options[stateA.cursorIndex] == "" {
					if stateA.cursorIndex - 1 >= 0 {
						stateA.cursorIndex--
					} else {
						stateA.cursorIndex++
					}
				}
			}
		case "down", "j":
			if stateA.cursorIndex < len(stateA.options)-1 {
				stateA.cursorIndex++
				if stateA.options[stateA.cursorIndex] == "" {
					if stateA.cursorIndex + 1 <= len(stateA.options) {
						stateA.cursorIndex++
					} else {
						stateA.cursorIndex--
					}
				}
			}
		case "enter":
			writeTo, err := utils.GetFile(stateA.options[stateA.cursorIndex], false)
			if err != nil {
				log.Fatal(err)
			}
			bytemark := stateB.bytemarks[stateB.cursorIndex]
			_, err = utils.Write(writeTo, bytemark.Table(), false)
			if err != nil {
				log.Fatal(err)
			}

			m.promptMenu.prompt = fmt.Sprintf(
				"Sent bytemark to %s",
				stateA.options[stateA.cursorIndex],
			)
			m.promptMenu.options = []string{"Okay"}
			m.promptMenu.cursorIndex = 0
			m.currentView = sentConfirmationView
		}
	}
	return m, nil
}

func sendBytemarkPromptView(m model) string {
	stateA := m.promptMenu
	//stateB := m.bytemarksManager

	var s string
	s += fmt.Sprintf("Send bytemark to hyperpath[%d] (enter)\n\n", stateA.cursorIndex)
	for i, hyperpath := range stateA.options {
		if hyperpath != "" {
			cursor := ""
			if stateA.cursorIndex == i {
				cursor = templates.Cursor()
				hyperpath = styles.HRender(styles.ProtonPurple, hyperpath)
			}
			s += fmt.Sprintf("%s%s\n", cursor, hyperpath)
		}
	}
	s += "\nGo back (esc)\n"
	return s
}

func updateSentConfirmation(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", "esc":
			m.wipePromptMenu()
			m.currentView = byteManagerView
		}
	}
	return m, nil
}
