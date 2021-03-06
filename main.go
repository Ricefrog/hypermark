package main

import (
	"bufio"
	"flag"
	"fmt"
	"hypermark/frontend"
	"hypermark/hackerNews"
	"hypermark/urlMode"
	"hypermark/utils"
	"log"
	"os"
)

// flags
var (
	k            string
	o            bool
	s            bool
	stdout       bool
	url          bool
	clipboardOut bool
	tui          bool
	//tuiTest      bool
	//test         bool
)

func init() {
	// k and s are mutually exclusive.
	flag.StringVar(&k, "k", "",
		"Save HN articles based on a keyword in the title.")
	flag.BoolVar(&o, "o", false,
		"Overwrite the target file instead of appending to the end.")
	flag.BoolVar(&s, "s", false, "Show all HN articles and exit.")
	flag.BoolVar(&stdout, "stdout", false, "Write output to stdout.")
	flag.BoolVar(&url, "url", false,
		"Use a URL from the system clipboard.")
	flag.BoolVar(&clipboardOut, "c", false,
		"Input will be written to the system clipboard.")
	flag.BoolVar(&tui, "tui", false,
		"Use TUI.")
	/*
	flag.BoolVar(&tuiTest, "tuiTest", false,
		"Run stubs to test TUI styling.")
	flag.BoolVar(&test, "test", false,
		"Run test stubs.")
	*/
}

func main() {
	flag.Parse()
	var outputPath *os.File
	var err error

	// outputPath is either a user-provided file or Stdout.
	// Accommodations made for system clipboard.
	outputPath, err = utils.ChooseOutputPath(
		flag.Args(), o, stdout, clipboardOut)
	if err != nil {
		if err.Error() == utils.EARLY_EXIT {
			return
		}
		log.Fatal(err)
	}
	defer outputPath.Close()

	if tui {
		frontend.SetOutputVars(
			outputPath,
			flag.Args(),
			o,
			stdout,
			clipboardOut,
		)
		frontend.ClearScreen()
		frontend.Start()
		return
	}

	/*
	if tuiTest {
		frontend.Test()
		return
	}

	if test {
		utils.TestStub()
		return
	}
	*/

	if url {
		bytemark, err := urlMode.BytemarkFromURL()
		if err != nil {
			log.Fatal(err)
		}

		output := bytemark.Table()
		writtenTo, err := utils.Write(outputPath, output, clipboardOut)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("bytemark for %s was written to %s.\n",
			bytemark.RootURL,
			writtenTo,
		)
		return
	}

	articles := hackerNews.ScrapeHN()
	if s {
		for i := 0; i < 30; i++ {
			title, sLink, cLink := hackerNews.GetHNInfo(articles[i])
			fmt.Printf("%d. %s\n%s\n%s\n\n", i+1, title, sLink, cLink)
		}
	} else if k != "" {
		fmt.Printf("Searching for articles with '%s' in the title.\n", k)

		output := ""
		articlesFound := 0
		for i := 0; i < 30; i++ {
			if articles[i].TitleContains(k) {
				output += articles[i].Table()
				articlesFound++
			}
		}
		writtenTo, err := utils.Write(outputPath, output, clipboardOut)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d articles found. Writing output to %s.\n",
			articlesFound,
			writtenTo,
		)
	} else {
		for i, article := range articles {
			fmt.Printf("%d %s\n", i+1, article.Title)
		}

		var userInput string
		fmt.Printf("\nArticles to save: (eg: 1 2 3, 1-3)\n")
		reader := bufio.NewReader(os.Stdin)
		userInput, err := reader.ReadString('\n')
		userInput = userInput[:len(userInput)-1] // remove trailing newline
		if err != nil {
			log.Fatal(err)
		}

		selections, err := utils.GetUserSelections(userInput)
		if err != nil {
			log.Fatal(err)
		}

		var output string
		for _, sel := range selections {
			output += articles[sel-1].Table()
		}

		writtenTo, err := utils.Write(outputPath, output, clipboardOut)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(
			"%d articles written to %s.\n",
			len(selections),
			writtenTo,
		)
	}
}
