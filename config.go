package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/afero"
	"github.com/stefanlogue/meteor/pkg/config"
)

const (
	defaultCommitTitleCharLimit      = 48
	defaultMessageTemplate           = "{{.Type}}{{if .Scope}}({{.Scope}}){{end}}{{if .IsBreakingChange}}!{{end}}: {{.Message}}"
	defaultMessageWithTicketTemplate = "{{.TicketNumber}}{{if .Scope}}({{.Scope}}){{end}}{{if .IsBreakingChange}}!{{end}}: <{{.Type}}> {{.Message}}"
)

// loadConfig loads the config file from the current directory or any parent
func loadConfig(fs afero.Fs) ([]huh.Option[string], []huh.Option[string], []huh.Option[string], bool, int, string, string, error) {
	filePath, err := config.FindConfigFile(fs)
	if err != nil {
		log.Debug("Error finding config file", "error", err)
		return config.DefaultPrefixes, nil, nil, true, defaultCommitTitleCharLimit, defaultMessageTemplate, defaultMessageWithTicketTemplate, nil
	}

	log.Debug("found config file", "path", filePath)

	c := config.New()

	err = c.LoadFile(filePath)
	if err != nil {
		return nil, nil, nil, true, defaultCommitTitleCharLimit, defaultMessageTemplate, defaultMessageWithTicketTemplate, fmt.Errorf("error parsing config file: %w", err)
	}

	if c.ShowIntro == nil {
		showIntro := true
		c.ShowIntro = &showIntro
	}

	if c.CommitTitleCharLimit == nil || *c.CommitTitleCharLimit < defaultCommitTitleCharLimit {
		commitTitleCharLimit := defaultCommitTitleCharLimit
		c.CommitTitleCharLimit = &commitTitleCharLimit
	}

	var messageTemplate, messageWithTicketTemplate string
	if c.MessageTemplate == nil {
		messageTemplate = defaultMessageTemplate
	} else {
		messageTemplate, err = config.ConvertTemplate(*c.MessageTemplate)
		if err != nil {
			log.Error("Error converting message template", "error", err)
			messageTemplate = defaultMessageTemplate
		}
	}
	c.MessageTemplate = &messageTemplate

	if c.MessageWithTicketTemplate == nil {
		messageWithTicketTemplate = defaultMessageWithTicketTemplate
	} else {
		messageWithTicketTemplate, err = config.ConvertTemplate(*c.MessageWithTicketTemplate)
		if err != nil {
			log.Error("Error converting message with ticket template", "error", err)
			messageWithTicketTemplate = defaultMessageWithTicketTemplate
		}
	}
	c.MessageWithTicketTemplate = &messageWithTicketTemplate

	return c.Prefixes.Options(), c.Coauthors.Options(), c.Boards.Options(), *c.ShowIntro, *c.CommitTitleCharLimit, messageTemplate, messageWithTicketTemplate, nil
}
