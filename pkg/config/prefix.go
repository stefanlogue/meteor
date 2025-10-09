package config

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type Prefix struct {
	T string  `json:"type"`
	D string  `json:"description"`
	E *string `json:"emoji,omitempty"`
}

type Prefixes []Prefix

var defaultPrefixData = Prefixes{
	{T: "feat", D: "a new feature"},
	{T: "fix", D: "a bug fix"},
	{T: "build", D: "changes that affect the build system or external dependencies"},
	{T: "chore", D: "changes to the build process or auxiliary tools and libraries"},
	{T: "ci", D: "changes to our CI configuration files and scripts"},
	{T: "docs", D: "documentation only changes"},
	{T: "perf", D: "a code change that improves performance"},
	{T: "refactor", D: "a code change that neither fixes a bug nor adds a feature"},
	{T: "revert", D: "reverts a previous commit"},
	{T: "style", D: "changes that do not affect the meaning of the code"},
	{T: "test", D: "adding missing tests or correcting existing tests"},
}

// GetDefaultPrefixOptions returns the default prefix options without emoji formatting
// for use when no custom prefixes are defined.
func GetDefaultPrefixOptions() []huh.Option[string] {
	return defaultPrefixData.OptionsWithEmojis(false)
}

func (p *Prefixes) Options() []huh.Option[string] {
	return p.OptionsWithEmojis(false)
}

// OptionsWithEmojis formats the prefixes into huh.Option slices, optionally including emojis.
// It dynamically sizes the columns based on the content for better alignment.
func (p *Prefixes) OptionsWithEmojis(useEmojis bool) []huh.Option[string] {
	prefixes := []Prefix(*p)
	if len(prefixes) == 0 {
		return GetDefaultPrefixOptions()
	}

	// Measure with grapheme-aware width
	maxTypeWidth := 0
	maxDescWidth := 0
	for _, pr := range prefixes {
		if w := runewidth.StringWidth(pr.T); w > maxTypeWidth {
			maxTypeWidth = w
		}
		if w := runewidth.StringWidth(pr.D); w > maxDescWidth {
			maxDescWidth = w
		}
	}

	// Dynamically size the emoji column from data (then add 1 pad).
	maxEmojiWidth := 0
	if useEmojis {
		for _, pr := range prefixes {
			if pr.E != nil && *pr.E != "" {
				if w := runewidth.StringWidth(*pr.E); w > maxEmojiWidth {
					maxEmojiWidth = w
				}
			}
		}
	}

	if maxEmojiWidth == 0 {
		maxEmojiWidth = 2 // sane default
	}
	emojiColWidth := maxEmojiWidth + 1 // padding

	typeStyle := lipgloss.NewStyle().Width(maxTypeWidth).Align(lipgloss.Left)
	emojiStyle := lipgloss.NewStyle().Width(emojiColWidth).Align(lipgloss.Left)
	separatorStyle := lipgloss.NewStyle()
	descriptionStyle := lipgloss.NewStyle().PaddingRight(1).Align(lipgloss.Left)

	var items []huh.Option[string]
	for _, prefix := range prefixes {
		typeWithEmoji := prefix.T

		var desc string
		if useEmojis && prefix.E != nil && *prefix.E != "" {
			typeWithEmoji = fmt.Sprintf("%s %s", prefix.T, *prefix.E)

			typeColumn := typeStyle.Render(prefix.T)
			emojiColumn := emojiStyle.Render(*prefix.E)
			separatorColumn := separatorStyle.Render(" - ")
			descriptionColumn := descriptionStyle.Render(prefix.D)

			desc = lipgloss.JoinHorizontal(
				lipgloss.Center,
				typeColumn,
				separatorColumn,
				descriptionColumn,
				emojiColumn,
			)
		} else {
			typeColumn := typeStyle.Render(prefix.T)
			separatorColumn := separatorStyle.Render(" - ")
			descriptionColumn := descriptionStyle.Render(prefix.D)

			desc = lipgloss.JoinHorizontal(
				lipgloss.Center,
				typeColumn,
				separatorColumn,
				descriptionColumn,
			)
		}
		items = append(items, huh.NewOption(desc, typeWithEmoji))
	}
	return items
}
