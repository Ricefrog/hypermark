package hackerNews

import (
	"fmt"
	"github.com/gocolly/colly"
	"strings"
)

const HN_URL = "https://news.ycombinator.com/"

type HNArticle struct {
	Title       string
	StoryLink   string
	CommentLink string
}

func (article HNArticle) GetInfo() []string {
	return []string{article.Title, article.StoryLink, article.CommentLink}
}

// Makes a markdown table with the article's info.
func (article HNArticle) GetTable() string {
	data := article.GetInfo()
	table := fmt.Sprintf("\n| %s |\n| :-- |\n| %s |\n| %s |\n",
		data[0], data[1], data[2])
	return table
}

// Returns whether or not the title of the article contains the search
// string. Can be improved upon later -> punctuation can create annoying
// situations.
func (article HNArticle) TitleContains(keyword string) bool {
	return strings.Contains(strings.ToLower(article.Title), keyword)
}

func ScrapeHN() []HNArticle {
	NUM_OF_ARTICLES := 30 // number of articles on the front page.
	index := 0
	cIndex := 0
	articles := make([]HNArticle, NUM_OF_ARTICLES)

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML(".athing", func(e *colly.HTMLElement) {
		articles[index].Title = e.ChildText("a.storylink")
		articles[index].StoryLink = e.ChildAttr("a.storylink", "href")
		index++
	})

	c.OnHTML(".athing + tr", func(e *colly.HTMLElement) {
		selector := "td.subtext a:nth-child(6)"
		articles[cIndex].CommentLink = e.ChildAttr(selector, "href")
		cIndex++
	})

	c.Visit(HN_URL)
	for i := 0; i < NUM_OF_ARTICLES; i++ {
		if articles[i].CommentLink != "" {
			articles[i].CommentLink = HN_URL + articles[i].CommentLink
		} else {
			articles[i].CommentLink = "No comments."
		}
	}
	return articles
}
