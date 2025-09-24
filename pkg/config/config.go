// Package config handles loading and parsing the configuration file for the application.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

type Config struct {
	PushAfterCommit           *bool     `json:"pushAfterCommit"`
	ShowIntro                 *bool     `json:"showIntro"`
	CommitTitleCharLimit      *int      `json:"commitTitleCharLimit"`
	CommitBodyCharLimit       *int      `json:"commitBodyCharLimit"`
	CommitBodyLineLength      *int      `json:"commitBodyLineLength"`
	MessageTemplate           *string   `json:"messageTemplate"`
	MessageWithTicketTemplate *string   `json:"messageWithTicketTemplate"`
	Prefixes                  Prefixes  `json:"prefixes"`
	Coauthors                 CoAuthors `json:"coauthors"`
	Boards                    Boards    `json:"boards"`
	Scopes                    Scopes    `json:"scopes"`
	ReadContributorsFromGit   *bool     `json:"readContributorsFromGit"`
}

// New returns a new Config
func New() *Config {
	return &Config{}
}

func (c *Config) LoadFile(filePath string) error {
	log.Debug("loading config file", "path", filePath)

	if filePath == "" {
		return errors.New("no path provided")
	}

	f, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	if err := json.Unmarshal(f, &c); err != nil {
		return fmt.Errorf("error parsing the json file: %w", err)
	}

	return nil
}
