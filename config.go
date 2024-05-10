package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/afero"
	"github.com/stefanlogue/meteor/pkg/config"
)

const defaultCommitTitleCharLimit = 48

// loadConfig loads the config file from the current directory or any parent
func loadConfig(fs afero.Fs) ([]huh.Option[string], []huh.Option[string], []huh.Option[string], bool, int, error) {
	filePath, err := config.FindConfigFile(fs)
	if err != nil {
		log.Debug("Error finding config file", "error", err)
		return config.DefaultPrefixes, nil, nil, true, defaultCommitTitleCharLimit, nil
	}

	log.Debug("found config file", "path", filePath)

	c := config.New()

	err = c.LoadFile(filePath)
	if err != nil {
		return nil, nil, nil, true, defaultCommitTitleCharLimit, fmt.Errorf("error parsing config file: %w", err)
	}

	if c.ShowIntro == nil {
		showIntro := true
		c.ShowIntro = &showIntro
	}

	if c.CommitTitleCharLimit == nil || *c.CommitTitleCharLimit < defaultCommitTitleCharLimit {
		commitTitleCharLimit := defaultCommitTitleCharLimit
		c.CommitTitleCharLimit = &commitTitleCharLimit
	}

	return c.Prefixes.Options(), c.Coauthors.Options(), c.Boards.Options(), *c.ShowIntro, *c.CommitTitleCharLimit, nil
}
