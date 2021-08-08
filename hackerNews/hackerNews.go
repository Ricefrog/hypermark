package hackerNews

import (
	"fmt"
	"github.com/gocolly/colly"
	"hypermark/utils"
)

const HN_URL = "https://news.ycombinator.com/"

func GetHNInfo(b utils.Bytemark) (title, storyLink, commentLink string) {
	title = b.Title
	storyLink = b.RootURL
	commentLink = b.Rows[0]

	return title, storyLink, commentLink
}

func ScrapeHN() []utils.Bytemark {
	NUM_OF_ARTICLES := 30 // number of articles on the front page.
	index := 0
	cIndex := 0
	articles := make([]utils.Bytemark, NUM_OF_ARTICLES)

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML(".athing", func(e *colly.HTMLElement) {
		articles[index].Title = e.ChildText("a.storylink")
		articles[index].RootURL = e.ChildAttr("a.storylink", "href")
		index++
	})

	c.OnHTML(".athing + tr", func(e *colly.HTMLElement) {
		selector := "td.subtext a:nth-child(6)"
		commentLink := e.ChildAttr(selector, "href")
		articles[cIndex].Rows = append(articles[cIndex].Rows, commentLink)
		cIndex++
	})

	c.Visit(HN_URL)
	for i := 0; i < NUM_OF_ARTICLES; i++ {
		articles[i].SetDateTimeNow()
		if len(articles[i].Rows) != 0 {
			articles[i].Rows[0] = "Comments: " + HN_URL + articles[i].Rows[0]
		} else {
			articles[i].Rows = append(articles[i].Rows, "No comments.")
		}
	}
	return articles
}
