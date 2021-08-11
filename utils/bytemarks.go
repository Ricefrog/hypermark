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
	b.Title = strings.ReplaceAll(b.Title, "|", "\\|") // escape | for markdown tables
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

func DeleteBytemark(original []Bytemark, index int) []Bytemark {
	deleted := make([]Bytemark, 0)
	for i, element := range original {
		if i != index {
			deleted = append(deleted, element)
		}
	}
	return deleted
}

func SwapBytemarks(original []Bytemark, indexA, indexB int) []Bytemark {
	swapped := make([]Bytemark, len(original))
	for i := 0; i < len(original); i++ {
		if i == indexA {
			swapped[i] = original[indexB]
		} else if i == indexB {
			swapped[i] = original[indexA]
		} else {
			swapped[i] = original[i]
		}
	}
	return swapped
}

func InsertBytemark(original []Bytemark, bm Bytemark, index int) []Bytemark {
	inserted := make([]Bytemark, len(original)+1)
	insertedIndex := 0
	for i := 0; i < len(original); i++ {
		if insertedIndex == index {
			inserted[insertedIndex] = bm
			insertedIndex++
		}
		inserted[insertedIndex] = original[i]
		insertedIndex++
	}
	return inserted
}
