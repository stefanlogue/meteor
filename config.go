package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
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

var defaultPrefixes = []list.Item{
	prefix{"feat", "A new feature"},
	prefix{"fix", "A bug fix"},
	prefix{"docs", "Documentation only changes"},
	prefix{"style", "Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)"},
	prefix{"refactor", "A code change that neither fixes a bug nor adds a feature"},
	prefix{"perf", "A code change that improves performance"},
	prefix{"test", "Adding missing tests or correcting existing tests"},
	prefix{"chore", "Changes to the build process or auxiliary tools and libraries such as documentation generation"},
	prefix{"revert", "Reverts a previous commit"},
	prefix{"ci", "Changes to our CI configuration files and scripts"},
}

const configFile = ".meteor.json"

func convertPrefixes(prefixes []prefix) []list.Item {
	items := []list.Item{}
	for _, prefix := range prefixes {
		items = append(items, prefix)
	}
	if len(items) == 0 {
		return defaultPrefixes
	}
	return items
}

func convertCoauthors(coauthors []coauthor) []list.Item {
	items := []list.Item{}
	for _, coauthor := range coauthors {
		items = append(items, coauthor)
	}
	if len(items) == 0 {
		return nil
	}
	items = append(items, coauthor{})
	copy(items[1:], items)
	items[0] = coauthor{"None", "no coauthors", false}
	return items
}

func convertBoards(boards []board) []list.Item {
	items := []list.Item{}
	for _, board := range boards {
		items = append(items, board)
	}
	if len(items) == 0 {
		return nil
	}
	return items
}

func loadConfigFile(path string) ([]list.Item, []list.Item, []list.Item, error) {
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

func loadConfig() ([]list.Item, []list.Item, []list.Item, error) {
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
			fmt.Println("Found config file at", filePath)
			return loadConfigFile(filePath)
		}

		targetPath += "/.."
	}
	return defaultPrefixes, nil, nil, nil
}
