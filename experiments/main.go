package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

func getAllFilenames(rootDir string) ([]string, error) {
	var filenames []string

	gitIgnore, err := parseGitIgnore(rootDir)
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories, we only want file paths
		if info.IsDir() {
			return nil
		}

		// Get the relative path from the root directory
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		// Check if the file should be ignored based on the .gitignore rules
		if !gitIgnore.MatchesPath(relPath) {
			filenames = append(filenames, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filenames, nil
}

func parseGitIgnore(rootDir string) (*gitignore.GitIgnore, error) {
	gitIgnorePath := filepath.Join(rootDir, ".gitignore")
	gitIgnoreData, err := os.ReadFile(gitIgnorePath)
	if err != nil {
		// Return an empty GitIgnore if .gitignore is not found or readable
		return gitignore.CompileIgnoreLines(), nil
	}
	return gitignore.CompileIgnoreLines(strings.Split(string(gitIgnoreData), "\n")...), nil
}

func main() {
	rootDir := "." // Change this to the directory you want to start from
	filenames, err := getAllFilenames(rootDir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print all filenames with their full paths
	for _, filename := range filenames {
		fmt.Println(filename)
	}
}
