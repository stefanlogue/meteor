package config

import (
	"testing"

	"github.com/charmbracelet/huh"
)

func TestPrefixes_Options(t *testing.T) {
	tests := []struct {
		name     string
		prefixes Prefixes
		want     []huh.Option[string]
	}{
		{
			name:     "empty prefixes returns default",
			prefixes: Prefixes{},
			want:     DefaultPrefixes,
		},
		{
			name: "prefixes without emojis",
			prefixes: Prefixes{
				{T: "feat", D: "a new feature"},
				{T: "fix", D: "a bug fix"},
			},
			want: []huh.Option[string]{
				huh.NewOption("feat - a new feature", "feat"),
				huh.NewOption("fix - a bug fix", "fix"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.prefixes.Options()
			if len(got) != len(tt.want) {
				t.Errorf("Options() returned %d items, want %d", len(got), len(tt.want))
				return
			}
			for i, option := range got {
				if option.Key != tt.want[i].Key || option.Value != tt.want[i].Value {
					t.Errorf("Options()[%d] = {Key: %q, Value: %q}, want {Key: %q, Value: %q}",
						i, option.Key, option.Value, tt.want[i].Key, tt.want[i].Value)
				}
			}
		})
	}
}

func TestPrefixes_OptionsWithEmojis(t *testing.T) {
	emojiPtr := func(s string) *string { return &s }

	tests := []struct {
		name      string
		prefixes  Prefixes
		useEmojis bool
		want      []huh.Option[string]
	}{
		{
			name:      "empty prefixes returns default",
			prefixes:  Prefixes{},
			useEmojis: true,
			want:      DefaultPrefixes,
		},
		{
			name: "prefixes with emojis enabled",
			prefixes: Prefixes{
				{T: "feat", D: "a new feature", E: emojiPtr("✨")},
				{T: "fix", D: "a bug fix", E: emojiPtr("🐛")},
			},
			useEmojis: true,
			want: []huh.Option[string]{
				huh.NewOption("feat ✨ - a new feature", "feat ✨"),
				huh.NewOption("fix 🐛 - a bug fix", "fix 🐛"),
			},
		},
		{
			name: "prefixes with emojis disabled",
			prefixes: Prefixes{
				{T: "feat", D: "a new feature", E: emojiPtr("✨")},
				{T: "fix", D: "a bug fix", E: emojiPtr("🐛")},
			},
			useEmojis: false,
			want: []huh.Option[string]{
				huh.NewOption("feat - a new feature", "feat"),
				huh.NewOption("fix - a bug fix", "fix"),
			},
		},
		{
			name: "prefixes with nil emojis",
			prefixes: Prefixes{
				{T: "feat", D: "a new feature", E: nil},
				{T: "fix", D: "a bug fix", E: nil},
			},
			useEmojis: true,
			want: []huh.Option[string]{
				huh.NewOption("feat - a new feature", "feat"),
				huh.NewOption("fix - a bug fix", "fix"),
			},
		},
		{
			name: "prefixes with empty emojis",
			prefixes: Prefixes{
				{T: "feat", D: "a new feature", E: emojiPtr("")},
				{T: "fix", D: "a bug fix", E: emojiPtr("")},
			},
			useEmojis: true,
			want: []huh.Option[string]{
				huh.NewOption("feat - a new feature", "feat"),
				huh.NewOption("fix - a bug fix", "fix"),
			},
		},
		{
			name: "mixed emojis - some with, some without",
			prefixes: Prefixes{
				{T: "feat", D: "a new feature", E: emojiPtr("✨")},
				{T: "fix", D: "a bug fix", E: nil},
				{T: "docs", D: "documentation", E: emojiPtr("📚")},
			},
			useEmojis: true,
			want: []huh.Option[string]{
				huh.NewOption("feat ✨ - a new feature", "feat ✨"),
				huh.NewOption("fix - a bug fix", "fix"),
				huh.NewOption("docs 📚 - documentation", "docs 📚"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.prefixes.OptionsWithEmojis(tt.useEmojis)
			if len(got) != len(tt.want) {
				t.Errorf("OptionsWithEmojis() returned %d items, want %d", len(got), len(tt.want))
				return
			}
			for i, option := range got {
				if option.Key != tt.want[i].Key || option.Value != tt.want[i].Value {
					t.Errorf("OptionsWithEmojis()[%d] = {Key: %q, Value: %q}, want {Key: %q, Value: %q}",
						i, option.Key, option.Value, tt.want[i].Key, tt.want[i].Value)
				}
			}
		})
	}
}
