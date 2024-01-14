package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
)

type prefix struct {
	T string `json:"type"`
	D string `json:"description"`
}

type coauthor struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Selected bool
}

type board struct {
	Name string `json:"name"`
}

type config struct {
	Prefixes  []prefix   `json:"prefixes"`
	Coauthors []coauthor `json:"coauthors"`
	Boards    []board    `json:"boards"`
}

var defaultPrefixes = []huh.Option[string]{
	huh.NewOption("feat - a new feature", "feat"),
	huh.NewOption("fix - a bug fix", "fix"),
	huh.NewOption("docs - documentation only changes", "docs"),
	huh.NewOption("style - changes that do not affect the meaning of the code", "style"),
	huh.NewOption("refactor - a code change that neither fixes a bug nor adds a feature", "refactor"),
	huh.NewOption("perf - a code change that improves performance", "perf"),
	huh.NewOption("test - adding missing tests or correcting existing tests", "test"),
	huh.NewOption("chore - changes to the build process or auxiliary tools and libraries", "chore"),
	huh.NewOption("revert - reverts a previous commit", "revert"),
	huh.NewOption("ci - changes to our CI configuration files and scripts", "ci"),
}

const configFile = ".meteor.json"

// convertPrefixes converts the given slice of prefixes into a slice of list.Items
func convertPrefixes(prefixes []prefix) []huh.Option[string] {
	if len(prefixes) == 0 {
		return defaultPrefixes
	}
	items := []huh.Option[string]{}
	for _, prefix := range prefixes {
		desc := fmt.Sprintf("%s - %s", prefix.T, prefix.D)
		items = append(items, huh.NewOption(desc, prefix.T))
	}
	return items
}

// convertCoauthors converts the given slice of coauthors into a slice of list.Items
func convertCoauthors(coauthors []coauthor) []huh.Option[string] {
	if len(coauthors) == 0 {
		return nil
	}
	items := []huh.Option[string]{}
	for _, coauthor := range coauthors {
		desc := fmt.Sprintf("%s <%s>", coauthor.Name, coauthor.Email)
		items = append(items, huh.NewOption(desc, desc))
	}
	items = append(items, huh.Option[string]{})
	copy(items[1:], items)
	items[0] = huh.NewOption[string]("no coauthors", "none")
	return items
}

// convertBoards converts the given slice of boards into a slice of list.Items
func convertBoards(boards []board) []huh.Option[string] {
	if len(boards) == 0 {
		return nil
	}
	items := []huh.Option[string]{}
	for _, board := range boards {
		items = append(items, huh.NewOption(board.Name, board.Name))
	}
	return items
}

// loadConfigFile loads the config file from the given path, and
// converts the contents into slices of list.Items
func loadConfigFile(path string) ([]huh.Option[string], []huh.Option[string], []huh.Option[string], error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error reading config file: %w", err)
	}
	var c config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, nil, nil, fmt.Errorf("error parsing config file: %w", err)
	}
	return convertPrefixes(c.Prefixes), convertCoauthors(c.Coauthors), convertBoards(c.Boards), nil
}

// loadConfig loads the config file from the current directory or any parent
func loadConfig() ([]huh.Option[string], []huh.Option[string], []huh.Option[string], error) {
	basePath, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting home dir: %w", err)
	}
	targetPath, err := os.Getwd()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting current dir: %w", err)
	}
	for {
		rel, _ := filepath.Rel(basePath, targetPath)
		if rel == "." {
			break
		}
		filePath := filepath.Join(targetPath, configFile)
		if _, err := os.Open(filePath); err == nil {
			return loadConfigFile(filePath)
		}

		targetPath += "/.."
	}
	return defaultPrefixes, nil, nil, nil
}
