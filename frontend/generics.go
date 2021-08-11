package frontend

import (
	"fmt"
	"hypermark/frontend/styles"
	"hypermark/frontend/templates"
)

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

func promptAndTextInputView(m model) string {
	state := m.promptAndTextInput

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		state.prompt,
		state.textInput.View(),
		state.footer,
	)
}
