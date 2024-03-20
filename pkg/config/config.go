package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

type Config struct {
	Prefixes  Prefixes  `json:"prefixes"`
	Coauthors CoAuthors `json:"coauthors"`
	Boards    Boards    `json:"boards"`
	ShowIntro *bool     `json:"showIntro"`
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
