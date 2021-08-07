package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"strconv"
	"regexp"
	"hypermark/hackerNews"
	"github.com/atotto/clipboard"
)

const (
	EARLY_EXIT = "42"
	HP_FILEPATH = "./hyperpaths"
)

func ArticlesToTable(articles []hackerNews.HNArticle) string {
	var output string
	for _, article := range articles {
		output = AppendArticleTable(output, article)
	}
	return output
}

func Write(
	outputPath *os.File,
	output string,
	clipboardOut bool,
) (string, error) {
	if !clipboardOut {
		_, err := outputPath.Write([]byte(output))
		return outputPath.Name(), err
	} else {
		err := clipboard.WriteAll(output)
		return "system clipboard", err
	}
}

func expandTilde(path string) string {
	usr, _ := user.Current()
	return strings.ReplaceAll(path, "~", usr.HomeDir)
}

func AppendArticleTable(str string, article hackerNews.HNArticle) string {
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

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	return fileInfo.IsDir()
}

func PathExists(path string) bool {
	pathExists := true
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			pathExists = false
		}
	}
	return pathExists
}

func removeBasename(path string) (retPath string) {
	parts := strings.Split(path, "/")
	for i := 0; i < len(parts)-1; i++ {
		retPath += parts[i]+"/"
	}
	return retPath
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

func GetUserSelections(userInput string) ([]int, error) {
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
		fmt.Printf("Proceed? Y/n: ")
		if _, err = fmt.Scan(&userInput); err != nil {
			log.Fatal(err)
		}

		switch strings.ToLower(userInput) {
		case "n":
			fmt.Printf("Exiting program.\n")
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

func getFile(
	fileName string,
	overwriteFile bool,
) (outputPath *os.File, err error) {
	fileExists := PathExists(fileName)
	// File will be overwritten.
	if overwriteFile {
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

func getHyperpathFromUser() (string, bool, error) {
	var hyperpath string
	var err error
	fmt.Printf("Enter hyperpath: ")
	if _, err = fmt.Scan(&hyperpath); err != nil {
		return "", false, err
	} else if strings.Contains(hyperpath, "~") {
		hyperpath = expandTilde(hyperpath)
	}
	fmt.Println()
	return hyperpath, true, nil
}

func getHyperpath() (string, error) {
	var hyperpathFile *os.File
	var err error
	var hyperpathChanged bool
	hyperpathFilePath := "env/hyperpath"

	// Create the hyperpathFile if it doesn't exist.
	if !PathExists(hyperpathFilePath) {
		hyperpathFile, err = os.OpenFile(
			hyperpathFilePath,
			os.O_CREATE,
			0666,
		)
	} else {
		hyperpathFile, err = os.OpenFile(
			hyperpathFilePath,
			os.O_RDONLY,
			0666,
		)
	}
	if err != nil {
		return "", err
	}

	data := make([]byte, 1024)
	var bytesRead int
	bytesRead, err = hyperpathFile.Read(data)
	if bytesRead != 0 && err != nil {
		return "", err
	}

	hyperpath := string(data[:bytesRead])
	if bytesRead == 0 {
		// Make this a separate function.
		var userInput string
		fmt.Printf("No hyperpath specified.\n")
		fmt.Printf("Would you like to set it now? Y/n: ")
		if _, err = fmt.Scan(&userInput); err != nil {
			log.Fatal(err)
		}

		switch strings.ToLower(userInput) {
		case "n":
			fmt.Printf("Exiting program.\n")
			return "", errors.New(EARLY_EXIT)
		}

		hyperpath, hyperpathChanged, err  = getHyperpathFromUser()
		if err != nil {
			return "", err
		}
	}
	for !PathExists(hyperpath) {
		fmt.Printf("\nInvalid file path: %s\n", hyperpath)
		hyperpath, hyperpathChanged, err  = getHyperpathFromUser()
		if err != nil {
			return "", err
		}
	}
	if hyperpathChanged {
		// The hyperpath on file was either wrong or non-existent.
		// Make a new hyperpath file and write a valid hyperpath to it.
		if err = os.Remove(hyperpathFilePath); err != nil {
			return "", err
		}
		hyperpathFile, err = os.OpenFile(
			hyperpathFilePath,
			os.O_CREATE|os.O_WRONLY,
			0666,
		)
		hyperpathFile.WriteString(hyperpath)
	}
	return hyperpath, nil
}

func ChooseOutputPath(
	tail []string,
	overwriteFile bool,
	writeToStdout bool,
	clipboardOut bool,
) (outputPath *os.File, err error) {
	if writeToStdout || clipboardOut {
		outputPath, err = os.Stdout, nil
	} else if len(tail) > 0 {
		// A specific file was specified.
		outputPath, err = getFile(tail[0], overwriteFile)
	} else {
		// Use hyperpath.
		var hyperpath string
		if hyperpath, err = getHyperpath(); err != nil {
			return outputPath, err
		}
		outputPath, err = getFile(hyperpath, overwriteFile)
	}
	return outputPath, err
}

func writeHyperpaths(hyperpaths []string) error {
	var err error
	// remove and create hyperpaths file
	if err = os.Remove(HP_FILEPATH); err != nil {
		return err
	}
	hyperpathFile, err := os.OpenFile(
		HP_FILEPATH,
		os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		return err
	}

	for i, hyperpath := range hyperpaths {
		line := fmt.Sprintf("%d: %s\n", i, hyperpath)
		hyperpathFile.WriteString(line)
	}
	return nil
}

func changeNthHyperpath(path string, n int) error {
	hyperpaths, err := GetAllHyperpaths()
	if err != nil {
		return err
	}
	if n > len(hyperpaths) {
		errorString := fmt.Sprintf("Cannot edit hyperpath[%d]: Out of range", n)
		return errors.New(errorString)
	}

	// Append new path or replace nth path.
	if n == len(hyperpaths) {
		hyperpaths = append(hyperpaths, path)
	} else {
		hyperpaths[n] = path
	}

	return writeHyperpaths(hyperpaths)
}

func pruneForHyperpaths(rawString string) []string {
	hyperpaths := make([]string, 0)
	hyperpathRegex := regexp.MustCompile(`\d+: .+`)

	lines := strings.Split(rawString, "\n")
	for _, line := range lines {
		match := hyperpathRegex.FindString(line)
		if match != "" {
			split := strings.Split(match, " ")
			hyperpaths = append(hyperpaths, split[1])
		}
	}

	return hyperpaths
}

func GetAllHyperpaths() ([]string, error) {
	var hyperpathsFile *os.File
	var err error

	if !PathExists(HP_FILEPATH) {
		return []string{}, errors.New("hyperpaths file does not exist.")
	} else {
		hyperpathsFile, err = os.OpenFile(
			HP_FILEPATH,
			os.O_RDONLY,
			0666,
		)
		if err != nil {
			return []string{}, err
		}
	}

	data := make([]byte, 1024)
	var bytesRead int
	bytesRead, err = hyperpathsFile.Read(data)
	if bytesRead == 0 {
		return []string{}, errors.New("hyperpaths file is empty.")
	}
	if err != nil {
		return []string{}, err
	}

	hyperpaths := pruneForHyperpaths(string(data))
	if len(hyperpaths) == 0 {
		return hyperpaths, errors.New("No hyperpaths found.")
	}
	return hyperpaths, nil
}

func EditNthHyperpath(path string, n int) (written, valid bool) {
	if strings.Contains(path, "~") {
		path = expandTilde(path)
	}
	// Check if the file exists. If so edit the hyperpath.
	if PathExists(path) && !isDirectory(path) {
		err := changeNthHyperpath(path, n)
		if err != nil {
			log.Fatal(err)
		}
		written, valid = true, true
		return
	}
	// If the file does not exist, check if the path to its
	// directory is valid. 
	pathToDir := removeBasename(path)
	if PathExists(pathToDir) && isDirectory(pathToDir) {
		valid = true
		return
	}
	return
}

func TestStub() {
	path := "/home/severian/terminus_est/tester.md"
	after := removeBasename(path)
	fmt.Printf("path: %s\nafter: %s\n", path, after)
}
