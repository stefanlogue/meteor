package config

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

type Prefix struct {
	T string `json:"type"`
	D string `json:"description"`
}

type Prefixes []Prefix

var DefaultPrefixes = []huh.Option[string]{
	huh.NewOption("feat - a new feature", "feat"),
	huh.NewOption("fix - a bug fix", "fix"),
	huh.NewOption("build - changes that affect the build system or external dependencies", "build"),
	huh.NewOption("chore - changes to the build process or auxiliary tools and libraries", "chore"),
	huh.NewOption("ci - changes to our CI configuration files and scripts", "ci"),
	huh.NewOption("docs - documentation only changes", "docs"),
	huh.NewOption("perf - a code change that improves performance", "perf"),
	huh.NewOption("refactor - a code change that neither fixes a bug nor adds a feature", "refactor"),
	huh.NewOption("revert - reverts a previous commit", "revert"),
	huh.NewOption("style - changes that do not affect the meaning of the code", "style"),
	huh.NewOption("test - adding missing tests or correcting existing tests", "test"),
}

func (p *Prefixes) Options() []huh.Option[string] {
	prefixes := []Prefix(*p)

	if len(prefixes) == 0 {
		return DefaultPrefixes
	}
	var items []huh.Option[string]
	for _, prefix := range prefixes {
		desc := fmt.Sprintf("%s - %s", prefix.T, prefix.D)
		items = append(items, huh.NewOption(desc, prefix.T))
	}
	return items
}
