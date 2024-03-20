package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/charmbracelet/log"
)

const (
	configFile = ".meteor.json"
)

// FindConfigFile will find the config files based in the rules below:
// 1. If the current directory contains a .meteor.json file, it will be used.
// 2. If the current directory does not contain a .meteor.json file, the parent
// 3. IF parent doesn't contain the .meteor.json file, the search will continue until the home directory is reached.
// 4. If no .meteor.json file is found, look in ~/.config/meteor/config.json
// 5. If no .meteor.json file is found, return an error
func FindConfigFile() (string, error) {
	if _, err := os.Stat(configFile); err == nil {
		return path.Join("./", configFile), nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home dir: %w", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current dir: %w", err)
	}

	for {
		rel, _ := filepath.Rel(homeDir, currentDir)
		if rel == ".." {
			break
		}

		filePath := filepath.Join(currentDir, configFile)
		log.Debug("checking for config file", "path", filePath)

		if _, err := os.Open(filePath); err == nil {
			return filePath, nil
		}

		currentDir += "/.."
	}

	xdgConfigFile := path.Join(homeDir, ".config/meteor/config.json")
	log.Debug("checking for config file", "path", xdgConfigFile)
	if _, err := os.Stat(xdgConfigFile); err == nil {
		return xdgConfigFile, nil
	}

	return "", errors.New("no config file found")
}
