package config

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/spf13/afero"
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
func FindConfigFile(fs afero.Fs, getwd func() (string, error), getHome func() (string, error)) (string, error) {
	if _, err := fs.Stat(configFile); err == nil {
		return path.Join("./", configFile), nil
	}

	homeDir, err := getHome()
	if err != nil {
		return "", fmt.Errorf("error getting home dir: %w", err)
	}

	currentDir, err := getwd()
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

		if _, err := fs.Open(filePath); err == nil {
			return filePath, nil
		}

		currentDir = filepath.Join(currentDir, "..")
	}

	xdgConfigFile := path.Join(homeDir, ".config/meteor/config.json")
	log.Debug("checking for config file", "path", xdgConfigFile)
	if _, err := fs.Stat(xdgConfigFile); err == nil {
		return xdgConfigFile, nil
	}

	return "", errors.New("no config file found")
}
