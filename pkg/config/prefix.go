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

const (
	defaultEmojiWidth = 2
	emojiPadding      = 1
)

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
	if len(*p) == 0 {
		return defaultPrefixData.OptionsWithEmojis(useEmojis)
	}

	// Single pass to measure all column widths
	measurements := measureColumnWidths(*p, useEmojis)

	// Create styles once
	styles := createColumnStyles(measurements)

	// Build options
	items := make([]huh.Option[string], 0, len(*p))
	for _, prefix := range *p {
		option := buildPrefixOption(prefix, useEmojis, styles)
		items = append(items, option)
	}
	return items
}

// columnMeasurements holds the measured widths for each column
type columnMeasurements struct {
	maxTypeWidth  int
	maxDescWidth  int
	maxEmojiWidth int
}

// columnStyles holds the lipgloss styles for each column
type columnStyles struct {
	typeStyle        lipgloss.Style
	emojiStyle       lipgloss.Style
	separatorStyle   lipgloss.Style
	descriptionStyle lipgloss.Style
}

// measureColumnWidths calculates the maximum width needed for each column in a single pass
func measureColumnWidths(prefixes []Prefix, useEmojis bool) columnMeasurements {
	var measurements columnMeasurements

	for _, prefix := range prefixes {
		if w := runewidth.StringWidth(prefix.T); w > measurements.maxTypeWidth {
			measurements.maxTypeWidth = w
		}
		if w := runewidth.StringWidth(prefix.D); w > measurements.maxDescWidth {
			measurements.maxDescWidth = w
		}

		if useEmojis && prefix.E != nil && *prefix.E != "" {
			if w := runewidth.StringWidth(*prefix.E); w > measurements.maxEmojiWidth {
				measurements.maxEmojiWidth = w
			}
		}
	}

	// Set default emoji width if none found
	if measurements.maxEmojiWidth == 0 {
		measurements.maxEmojiWidth = defaultEmojiWidth
	}

	return measurements
}

// createColumnStyles creates the lipgloss styles based on measurements
func createColumnStyles(measurements columnMeasurements) columnStyles {
	emojiColWidth := measurements.maxEmojiWidth + emojiPadding

	return columnStyles{
		typeStyle:        lipgloss.NewStyle().Width(measurements.maxTypeWidth).Align(lipgloss.Left),
		emojiStyle:       lipgloss.NewStyle().Width(emojiColWidth).Align(lipgloss.Left),
		separatorStyle:   lipgloss.NewStyle(),
		descriptionStyle: lipgloss.NewStyle().PaddingRight(1).Align(lipgloss.Left),
	}
}

// buildPrefixOption creates a single huh.Option for a prefix
func buildPrefixOption(prefix Prefix, useEmojis bool, styles columnStyles) huh.Option[string] {
	typeWithEmoji := prefix.T
	hasEmoji := useEmojis && prefix.E != nil && *prefix.E != ""

	if hasEmoji {
		typeWithEmoji = fmt.Sprintf("%s %s", prefix.T, *prefix.E)
	}

	// Build the display string
	var desc string
	if hasEmoji {
		desc = lipgloss.JoinHorizontal(
			lipgloss.Center,
			styles.typeStyle.Render(prefix.T),
			styles.separatorStyle.Render(" - "),
			styles.descriptionStyle.Render(prefix.D),
			styles.emojiStyle.Render(*prefix.E),
		)
	} else {
		desc = lipgloss.JoinHorizontal(
			lipgloss.Center,
			styles.typeStyle.Render(prefix.T),
			styles.separatorStyle.Render(" - "),
			styles.descriptionStyle.Render(prefix.D),
		)
	}

	return huh.NewOption(desc, typeWithEmoji)
}
