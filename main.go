package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"strings"
	"flag"
	"os"
	"log"
)

const HN_URL = "https://news.ycombinator.com/"

var k string
var o bool
var s bool
func init() {
	// k and s are mutually exclusive.
	flag.StringVar(&k, "k", "",
				   "Save articles based on a keyword in the title.")
	flag.BoolVar(&o, "o", false,
			     "Overwrite the target file instead of appending to the end.")
	flag.BoolVar(&s, "s", false, "Show all articles and exit.")
}

func scrapeHN() ([]string, []string, []string) {
	index := 0
	cIndex := 0
	titles := make([]string, 30)
	storylinks := make([]string, 30)
	commentlinks := make([]string, 30)

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML(".athing", func(e *colly.HTMLElement) {
		titles[index] = e.ChildText("a.storylink")
		storylinks[index] = e.ChildAttr("a.storylink", "href")
		index++
	})

	c.OnHTML(".athing + tr", func(e *colly.HTMLElement) {
		commentlinks[cIndex] = e.ChildAttr("td.subtext a:nth-child(6)", "href")
		cIndex++
	})

	c.Visit(HN_URL)
	for i, _ := range titles {
		if commentlinks[i] != "" {
			commentlinks[i] = HN_URL + commentlinks[i]
		} else {
			commentlinks[i] = "No comments."
		}
	}
	return titles, storylinks, commentlinks
}

func appendArticle(str, title, storylink, commentlink string) string {
	article := fmt.Sprintf("\n| %s |\n| :-- |\n| %s |\n| %s |\n",
							title, storylink, commentlink)
	return str + article
}

func main() {
	flag.Parse()

	var outputPath *os.File
	tail := flag.Args()
	if len(tail) > 0 {
		var err error
		fileName := tail[0]
		fileExists := true
		if _, err = os.Stat(fileName); err != nil {
			if os.IsNotExist(err) {
				fileExists = false
			}
		}

		if o {
			outputPath, err = os.OpenFile(
									fileName,
									os.O_WRONLY|os.O_CREATE,
									0666,
								)
			if fileExists {
				var choice string
				fmt.Printf(
					"The file '%s' will be overwritten.\n",
					fileName,
				)
				fmt.Printf("Proceed? y/n: ")
				if _, err = fmt.Scan(&choice); err != nil {
					log.Fatal(err)
				}

				switch choice {
				case "y":
					break
				case "n":
					fmt.Printf("Exiting program.\n")
					return
				default:
					fmt.Printf("Invalid option (%s).\n", choice)
					return
				}
			}
		} else {
			fmt.Println("Appending to file.")
			outputPath, err = os.OpenFile(
				fileName,
				os.O_APPEND|os.O_CREATE|os.O_WRONLY,
				0666,
			)
		}
		if err != nil {
			log.Fatal(err)
		}

		defer outputPath.Close()
	} else {
		outputPath = os.Stdout
	}

	titles, storylinks, commentlinks := scrapeHN()

	if s {
		for i := 0; i < 30; i++ {
			fmt.Printf("%d. %s\n%s\n%s\n\n",
						i+1,
						titles[i],
						storylinks[i],
						commentlinks[i])
		}
	} else if k != "" {
		fmt.Printf("Searching for articles with '%s' in the title.\n", k)
		output := ""
		for i := 0; i < 30; i++ {
			if strings.Contains(strings.ToLower(titles[i]), k) {
				output = appendArticle(output,
									   titles[i],
									   storylinks[i],
									   commentlinks[i])
			}
		}
		fmt.Printf("Writing output to %s.\n", outputPath.Name())
		_, err := outputPath.Write([]byte(output))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("default mode")
	}
}
