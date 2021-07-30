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

const EARLY_EXIT = "42"

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

func doesFileExist(fileName string) bool {
	fileExists := true
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			fileExists = false
		}
	}
	return fileExists
}

func getFile(
	fileName string,
	overwriteFile bool,
) (outputPath *os.File, err error) {
	fileExists := doesFileExist(fileName)
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
	if !doesFileExist(hyperpathFilePath) {
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
	for !doesFileExist(hyperpath) {
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
