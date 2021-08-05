package templates

import (
	"hypermark/frontend/styles"
)

func Prompt(prompt string) string {
	var s string
	s += styles.PromptStyle.Render(prompt)
	s += "\n"
	return s
}

func Cursor() string {
	return styles.CursorStyle.Render("-> ")
}
