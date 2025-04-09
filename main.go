package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type application struct {
	rootDir    string
	searchStr  string
	replaceStr string
	dryRun     bool
}

func getSerchStr(dir string) (string, error) {
	path := dir + "/go.mod"
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "module") {
			return strings.Split(scanner.Text(), " ")[1], nil
		} else {
			return "", errors.New("module not found")
		}
	}

	return "", nil
}

func getCurrentDir() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return currentDir, nil
}

func processFile(path, searchStr, replaceStr string, dryRun bool) (bool, int, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, 0, err
	}

	file, err := os.Open(path)
	if err != nil {
		return false, 0, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return false, 0, err
	}

	// Convert to string for easier handling
	text := string(content)

	// Check if the search string exists in the file
	if !strings.Contains(text, searchStr) {
		return false, 0, nil // No matches found
	}

	// Count the number of replacements
	count := strings.Count(text, searchStr)

	// Replace all occurrences
	newText := strings.ReplaceAll(text, searchStr, replaceStr)

	// If this is a dry run, don't actually modify the file
	if dryRun {
		return true, count, nil
	}

	// Write the modified content back to the file
	err = os.WriteFile(path, []byte(newText), info.Mode())
	if err != nil {
		return false, 0, err
	}

	return true, count, nil
}

func (app *application) search() error {
	fmt.Println("Searching for: ", app.searchStr)
	fmt.Println("in the directory: ", app.rootDir)
	fmt.Println("and replacing with: ", app.replaceStr)
	excludeDirs := []string{".git"}

	err := filepath.Walk(app.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// skipping excluded directories
		if info.IsDir() {
			for _, dir := range excludeDirs {
				if strings.Contains(path, dir) {
					return filepath.SkipDir
				}
			}
		}

		// only process files with the .go extension
		if filepath.Ext(path) == ".go" {
			changed, count, err := processFile(path, app.searchStr, app.replaceStr, app.dryRun)
			if err != nil {
				return err
			}

			if changed {
				fmt.Printf("Modified %s (%d replacements)\n", path, count)
			}
		}

		return nil
	})

	return err
}

func main() {
	// Define flags
	var dryRun bool
	flag.BoolVar(&dryRun, "dry-run", false, "Perform a dry run without modifying files")
	flag.BoolVar(&dryRun, "dryrun", false, "Perform a dry run without modifying files (alias)")

	replaceStr := flag.String("rs", "", "The term to replace")
	help := flag.Bool("help", false, "Show help")

	flag.Parse()

	if *help {
		fmt.Println("Be sure to read the README.md file for more information")
		fmt.Println("Usage: find-replace [options]")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *replaceStr == "" {
		log.Fatalln("replace string[-rs] is required, but not provided")
	}
	// rootDir represents where the application is being run from
	rootDir, err := getCurrentDir()
	if err != nil {
		log.Fatalln(err)
	}

	// next we need to open up the go.mod file because that is the term that we are searching for
	searchStr, err := getSerchStr(rootDir)
	if err != nil {
		log.Fatalln(err)
	}

	app := application{
		rootDir:    rootDir,
		searchStr:  searchStr,
		replaceStr: *replaceStr,
		dryRun:     dryRun,
	}

	err = app.search()
	if err != nil {
		log.Fatalln(err)
	}
}
