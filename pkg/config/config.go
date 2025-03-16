package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

type Config struct {
	ShowIntro                 *bool     `json:"showIntro"`
	CommitTitleCharLimit      *int      `json:"commitTitleCharLimit"`
	MessageTemplate           *string   `json:"messageTemplate"`
	MessageWithTicketTemplate *string   `json:"messageWithTicketTemplate"`
	Prefixes                  Prefixes  `json:"prefixes"`
	Coauthors                 CoAuthors `json:"coauthors"`
	Boards                    Boards    `json:"boards"`
	Scopes                    Scopes    `json:"scopes"`
	AskBreakingChange         *bool     `json:"askBreakingChange"`
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
