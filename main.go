package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"hypermark/hackerNews"
	"strings"
)

const EARLY_EXIT = "42"

// flags
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

func appendArticleTable(str string, article hackerNews.HNArticle) string {
	return str + article.GetTable()
}

func contains(arr []int, search int) bool {
	for _, el := range arr {
		if el == search {
			return true
		}
	}
	return false
}

func removeDuplicates(ints []int) []int {
	ret := make([]int, 0)
	for _, num := range ints {
		if !contains(ret, num) {
			ret = append(ret, num)
		}
	}
	return ret
}

func getUserSelections(userInput string) ([]int, error) {
	inputStrings := strings.Split(userInput, " ")
	selections := make([]int, 0)

	for _, str := range inputStrings {
		isRange, _ := regexp.MatchString(`\d+\-\d+`, str)
		isInt, _ := regexp.MatchString(`\d+`, str)
		if isRange {
			a := strings.Split(str, "-")
			from, _ := strconv.Atoi(a[0])
			to, _ := strconv.Atoi(a[1])

			if (to < from) || (from < 1 || to > 30) {
				errMessage := fmt.Sprintf("Invalid range: %s", str)
				return make([]int, 0), errors.New(errMessage)
			} else {
				for i := from; i <= to; i++ {
					selections = append(selections, i)
				}
			}
		} else if isInt {
			sel, _ := strconv.Atoi(str)
			if sel < 1 || sel > 30 {
				errMessage := fmt.Sprintf("Invalid selection: %s", str)
				return make([]int, 0), errors.New(errMessage)
			}
			selections = append(selections, sel)
		}
	}

	return removeDuplicates(selections), nil
}

func overwriteMode(
	fileExists bool,
	fileName string,
) (outputPath *os.File, err error) {
	if fileExists {
		var userInput string
		fmt.Printf(
			"The file '%s' will be overwritten.\n",
			fileName,
		)
		fmt.Printf("Proceed? y/n: ")
		if _, err = fmt.Scan(&userInput); err != nil {
			log.Fatal(err)
		}

		switch userInput {
		case "y":
			break
		case "n":
			fmt.Printf("Exiting program.\n")
			return outputPath, errors.New(EARLY_EXIT)
		default:
			fmt.Printf("Invalid option (%s).\n", userInput)
			return outputPath, errors.New(EARLY_EXIT)
		}
		// Wipe the data on the file.
		if err = os.Remove(fileName); err != nil {
			log.Fatal(err)
		}
	}
	// Create the file.
	outputPath, err = os.OpenFile(
		fileName,
		os.O_CREATE|os.O_WRONLY,
		0666,
	)
	return outputPath, err
}

func getFile(fileName string) (outputPath *os.File, err error) {
	fileExists := true
	if _, err = os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			fileExists = false
		}
	}

	// File will be overwritten.
	if o {
		outputPath, err = overwriteMode(fileExists, fileName)
	} else {
	// File will be appended.
		outputPath, err = os.OpenFile(
			fileName,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0666,
		)
	}
	return outputPath, err
}

func chooseOutputPath(tail []string) (outputPath *os.File, err error) {
	if len(tail) > 0 {
		outputPath, err = getFile(tail[0])
	} else {
		outputPath, err = os.Stdout, nil
	}
	return outputPath, err
}

func main() {
	flag.Parse()

	var outputPath *os.File
	var err error

	// outputPath is either a user-provided file or Stdout.
	if outputPath, err = chooseOutputPath(flag.Args()); err != nil {
		if err.Error() == EARLY_EXIT {
			return
		}
		log.Fatal(err)
	}
	defer outputPath.Close()

	articles := hackerNews.ScrapeHN()

	if s {
		for i := 0; i < 30; i++ {
			data := articles[i].GetInfo()
			fmt.Printf("%d. %s\n%s\n%s\n\n", i+1, data[0], data[1], data[2])
		}
		return
	} else if k != "" {
		fmt.Printf("Searching for articles with '%s' in the title.\n", k)

		output := ""
		articlesFound := 0
		for i := 0; i < 30; i++ {
			if articles[i].TitleContains(k) {
				output = appendArticleTable(output, articles[i])
				articlesFound++
			}
		}
		fmt.Printf("%d articles found. Writing output to %s.\n",
			articlesFound,
			outputPath.Name())
		_, err := outputPath.Write([]byte(output))
		if err != nil {
			log.Fatal(err)
		}
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

		selections, err := getUserSelections(userInput)
		if err != nil {
			log.Fatal(err)
		}

		var output string
		for _, sel := range selections {
			output = appendArticleTable(output, articles[sel-1])
		}

		_, err = outputPath.Write([]byte(output))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(
			"%d articles written to %s.\n",
			len(selections),
			outputPath.Name(),
		)
	}
}
