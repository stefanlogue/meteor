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

type LoadConfigReturn struct {
	MessageTemplate           string
	MessageWithTicketTemplate string
	Prefixes                  []huh.Option[string]
	Coauthors                 []huh.Option[string]
	Boards                    []huh.Option[string]
	CommitTitleCharLimit      int
	ShowIntro                 bool
}

// loadConfig loads the config file from the current directory or any parent
func loadConfig(fs afero.Fs) (LoadConfigReturn, error) {
	filePath, err := config.FindConfigFile(fs)
	if err != nil {
		log.Debug("Error finding config file", "error", err)
		return LoadConfigReturn{
			MessageTemplate:           defaultMessageTemplate,
			MessageWithTicketTemplate: defaultMessageWithTicketTemplate,
			Prefixes:                  config.DefaultPrefixes,
			CommitTitleCharLimit:      defaultCommitTitleCharLimit,
			ShowIntro:                 true,
		}, nil
	}

	log.Debug("found config file", "path", filePath)

	c := config.New()

	err = c.LoadFile(filePath)
	if err != nil {
		return LoadConfigReturn{
			MessageTemplate:           defaultMessageTemplate,
			MessageWithTicketTemplate: defaultMessageWithTicketTemplate,
			CommitTitleCharLimit:      defaultCommitTitleCharLimit,
			ShowIntro:                 true,
		}, fmt.Errorf("error parsing config file: %w", err)
	}

	if c.ShowIntro == nil {
		showIntro := true
		c.ShowIntro = &showIntro
	}

	if c.CommitTitleCharLimit == nil || *c.CommitTitleCharLimit < defaultCommitTitleCharLimit {
		commitTitleCharLimit := defaultCommitTitleCharLimit
		c.CommitTitleCharLimit = &commitTitleCharLimit
	}

	messageTemplate := defaultMessageTemplate
	if c.MessageTemplate != nil {
		messageTemplate, err = config.ConvertTemplate(*c.MessageTemplate)
		if err != nil {
			log.Error("Error converting message template", "error", err)
			messageTemplate = defaultMessageTemplate
		}
	}
	c.MessageTemplate = &messageTemplate

	messageWithTicketTemplate := defaultMessageWithTicketTemplate
	if c.MessageWithTicketTemplate != nil {
		messageWithTicketTemplate, err = config.ConvertTemplate(*c.MessageWithTicketTemplate)
		if err != nil {
			log.Error("Error converting message with ticket template", "error", err)
			messageWithTicketTemplate = defaultMessageWithTicketTemplate
		}
	}
	c.MessageWithTicketTemplate = &messageWithTicketTemplate

	return LoadConfigReturn{
		MessageTemplate:           messageTemplate,
		MessageWithTicketTemplate: messageWithTicketTemplate,
		Prefixes:                  c.Prefixes.Options(),
		Coauthors:                 c.Coauthors.Options(),
		Boards:                    c.Boards.Options(),
		CommitTitleCharLimit:      *c.CommitTitleCharLimit,
		ShowIntro:                 *c.ShowIntro,
	}, nil
}
