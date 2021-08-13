package urlMode

import (
	//"fmt"
	"github.com/atotto/clipboard"
	"github.com/gocolly/colly"
	"hypermark/utils"
)

func createBytemark(url string) (utils.Bytemark, error) {
	bytemark := utils.Bytemark{RootURL: url}
	var retErr error
	c := colly.NewCollector()

	/*
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	*/

	c.OnHTML("title", func(e *colly.HTMLElement) {
		bytemark.Title = e.Text
	})

	c.OnError(func(r *colly.Response, err error) {
		retErr = err
	})

	c.Visit(url)

	bytemark.SetDateTimeNow()
	return bytemark, retErr
}

func BytemarkFromURL() (utils.Bytemark, error) {
	url, err := clipboard.ReadAll()
	if err != nil {
		return utils.Bytemark{}, err
	}

	return createBytemark(url)
}
