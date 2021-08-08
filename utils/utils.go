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
	"github.com/atotto/clipboard"
)

const (
	EARLY_EXIT = "42"
	HP_FILEPATH = "./hyperpaths"
)

func DeleteElement(original []string, index int) []string {
	deleted := make([]string, 0)
	for i, element := range original {
		if i != index {
			deleted = append(deleted, element)
		}
	}
	return deleted
}

func SwapElements(original []string, indexA, indexB int) []string {
	swapped := make([]string, len(original))
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

func AppendArticleTable(str string, article Bytemark) string {
	return str + article.Table()
}

func ArticlesToTable(articles []Bytemark) string {
	var output string
	for _, article := range articles {
		output = AppendArticleTable(output, article)
	}
	return output
}

func ExpandTilde(path string) string {
	usr, _ := user.Current()
	return strings.ReplaceAll(path, "~", usr.HomeDir)
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

func asciiCodes(str string) []int {
	ret := make([]int, 0)
	for i := 0; i < len(str); i++ {
		ret = append(ret, int(str[i]))
	}
	return ret
}

func RemoveEmptyStrings(arr []string) (retArr []string) {
	for _, str := range arr {
		// All whitespace or some sort of null character.
		if len(strings.TrimSpace(str)) != 0 && int(str[0]) != 0 {
			retArr = append(retArr, str)
		}
	}
	return
}

func CreateFile(path string) (*os.File, error) {
	outputPath, err := os.OpenFile(
		path,
		os.O_CREATE,
		0666,
	)
	return outputPath, err
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

func trimRows(raw []string) (trimmed []string) {
	for _, field := range raw {
		trimmed = append(trimmed, field[2:len(field)-2])
	}
	return
}

func tableToBytemark(table string) Bytemark {
	fields := trimRows(DeleteElement(strings.Split(table, "\n"), 1))
	bytemark := Bytemark{
		Title: fields[0],
		DateTime: fields[1],
		RootURL: fields[2],
	}
	if len(fields) > 3 {
		for i := 3; i < len(fields); i++ {
			bytemark.Rows = append(bytemark.Rows, fields[i])
		}
	}
	return bytemark
}

func FileToBytemarks(file *os.File) ([]Bytemark, error) {
	bytemarks := make([]Bytemark, 0)

	data := make([]byte, 1024)
	bytesRead, err := file.Read(data)
	if bytesRead == 0 {
		return bytemarks, nil
	}
	if err != nil {
		return bytemarks, err
	}

	tables := RemoveEmptyStrings(strings.Split(string(data), "\n\n"))
	for _, table := range tables {
		bytemarks = append(bytemarks, tableToBytemark(table))
	}
	return bytemarks, nil
}

func TestStub() {
	file, err := os.OpenFile(
		"/home/severian/terminus_est/tester.md",
		os.O_RDONLY,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}

	bytemarks, _ := FileToBytemarks(file)
	for _, b := range bytemarks {
		fmt.Printf(
			"Title: %s\nDateTime: %s\nRootURL: %s\nRows: %v\n",
			b.Title,
			b.DateTime,
			b.RootURL,
			b.Rows,
		)
	}
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
	fmt.Printf("Enter hyperpath[0]: ")
	if _, err = fmt.Scan(&hyperpath); err != nil {
		return "", false, err
	} else if strings.Contains(hyperpath, "~") {
		hyperpath = ExpandTilde(hyperpath)
	}
	fmt.Println()
	return hyperpath, true, nil
}

func getMainHyperpath() (string, error) {
	var err error
	var hyperpathChanged bool

	// Create the hyperpathFile if it doesn't exist.
	if !PathExists(HP_FILEPATH) {
		_, err = os.OpenFile(
			HP_FILEPATH,
			os.O_CREATE,
			0666,
		)
	}
	if err != nil {
		return "", err
	}

	allHyperpaths, err := GetAllHyperpaths()
	if err != nil {
		return "", err
	}

	mainHyperpath := allHyperpaths[0]
	if mainHyperpath == "" {
		// Make this a separate function.
		var userInput string
		fmt.Printf("No hyperpath[0] specified.\n")
		fmt.Printf("Would you like to set it now? Y/n: ")
		if _, err = fmt.Scan(&userInput); err != nil {
			log.Fatal(err)
		}

		switch strings.ToLower(userInput) {
		case "n":
			fmt.Printf("Exiting program.\n")
			return "", errors.New(EARLY_EXIT)
		}

		mainHyperpath, hyperpathChanged, err  = getHyperpathFromUser()
		if err != nil {
			return "", err
		}
	}
	for !PathExists(mainHyperpath) {
		fmt.Printf("\nInvalid file path: %s\n", mainHyperpath)
		mainHyperpath, hyperpathChanged, err  = getHyperpathFromUser()
		if err != nil {
			return "", err
		}
	}
	if hyperpathChanged {
		// hyperpath[0] on file was either wrong or non-existent.
		// Insert the new hyperpath[0] into the hyperpath file.
		err := changeNthHyperpath(mainHyperpath, 0)
		if err != nil {
			return "", err
		}
	}
	return mainHyperpath, nil
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
		if hyperpath, err = getMainHyperpath(); err != nil {
			return outputPath, err
		}
		outputPath, err = getFile(hyperpath, overwriteFile)
	}
	return outputPath, err
}

func WriteHyperpaths(hyperpaths []string) error {
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

	return WriteHyperpaths(hyperpaths)
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
		return []string{""}, nil
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
