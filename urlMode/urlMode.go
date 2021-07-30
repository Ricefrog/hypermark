package urlMode

import (
	"fmt"
	"log"
	"github.com/gocolly/colly"
	"github.com/atotto/clipboard"
)

type GenericDatamark struct {
	URL string
	Title string
	Table string
}

func (b *GenericDatamark) createTable() {
	b.Table = fmt.Sprintf("\n| %s |\n| :-- |\n| %s |\n", b.Title, b.URL)
}

func TestingStub() {
	// test if reading and writing from the clipboard works.
	fromClipboard, err := clipboard.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("From clipboard: %s\n", fromClipboard)

	toClipboard := "crepuscular"
	fmt.Printf("Writing %s to clipboard.\n", toClipboard)
	err = clipboard.WriteAll(toClipboard)
	if err != nil {
		log.Fatal(err)
	}
}

func createDatamark(url string) GenericDatamark {
	datamark := GenericDatamark{URL: url}
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		datamark.Title = e.Text
	})

	c.Visit(url)
	datamark.createTable()
	return datamark
}

func ClipToClip() {
	url, err := clipboard.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	datamark := createDatamark(url)
	err = clipboard.WriteAll(datamark.Table)
	if err != nil {
		log.Fatal(err)
	}
}

func DatamarkFromURL() (GenericDatamark, error) {
	url, err := clipboard.ReadAll()
	return createDatamark(url), err
}
