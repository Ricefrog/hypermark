package utils

import (
	"fmt"
	"time"
	"strings"
)

type Bytemark struct {
	Title string
	DateTime string
	RootURL string
	Rows []string
}

func (b *Bytemark) Table() string {
	table := fmt.Sprintf(
		"| %s |\n| :-- |\n| %s |\n| %s |\n",
		b.Title,
		b.DateTime,
		b.RootURL,
	)
	for _, row := range b.Rows {
		table += fmt.Sprintf("| %s |\n", row)
	}
	table += "\n"

	return table
}

func (b *Bytemark) SetDateTime(timeUsed time.Time) {
	year, month, day := timeUsed.Date()
	hour, min, _ := timeUsed.Clock()

	dateString := fmt.Sprintf("%d/%d/%d", month, day, year)
	timeString := fmt.Sprintf("%d:%d", hour, min)
	b.DateTime = fmt.Sprintf("%s %s", dateString, timeString)
}

func (b *Bytemark) SetDateTimeNow() {
	b.SetDateTime(time.Now())
}

// Returns whether or not the title of the article contains the search
// string. Can be improved upon later -> punctuation can create annoying
// situations.
func (article Bytemark) TitleContains(keyword string) bool {
	return strings.Contains(strings.ToLower(article.Title), keyword)
}
