package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/stefanlogue/meteor/pkg/config"
)

const configFile = ".meteor.json"

// loadConfigFile loads the config file from the given path, and
// converts the config file into a slice of huh.Option[string]
func loadConfigFile(path string) ([]huh.Option[string], []huh.Option[string], []huh.Option[string], bool, error) {
	c := config.New()

	err := c.LoadFile(path)
	if err != nil {
		return nil, nil, nil, true, fmt.Errorf("error parsing config file: %w", err)
	}

	if c.ShowIntro == nil {
		showIntro := true
		c.ShowIntro = &showIntro
	}

	return c.Prefixes.Option(), c.Coauthors.Options(), c.Boards.Options(), *c.ShowIntro, nil
}

// loadConfig loads the config file from the current directory or any parent
func loadConfig() ([]huh.Option[string], []huh.Option[string], []huh.Option[string], bool, error) {
	basePath, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, nil, true, fmt.Errorf("error getting home dir: %w", err)
	}
	targetPath, err := os.Getwd()
	if err != nil {
		return nil, nil, nil, true, fmt.Errorf("error getting current dir: %w", err)
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
	return config.DefaultPrefixes, nil, nil, true, nil
}
