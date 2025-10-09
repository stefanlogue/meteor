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
			want: []huh.Option[string]{
				huh.NewOption("feat     - a new feature ", "feat"),
				huh.NewOption("fix      - a bug fix ", "fix"),
				huh.NewOption("build    - changes that affect the build system or external dependencies ", "build"),
				huh.NewOption("chore    - changes to the build process or auxiliary tools and libraries ", "chore"),
				huh.NewOption("ci       - changes to our CI configuration files and scripts ", "ci"),
				huh.NewOption("docs     - documentation only changes ", "docs"),
				huh.NewOption("perf     - a code change that improves performance ", "perf"),
				huh.NewOption("refactor - a code change that neither fixes a bug nor adds a feature ", "refactor"),
				huh.NewOption("revert   - reverts a previous commit ", "revert"),
				huh.NewOption("style    - changes that do not affect the meaning of the code ", "style"),
				huh.NewOption("test     - adding missing tests or correcting existing tests ", "test"),
			},
		},
		{
			name: "prefixes without emojis",
			prefixes: Prefixes{
				{T: "feat", D: "a new feature"},
				{T: "fix", D: "a bug fix"},
			},
			want: []huh.Option[string]{
				huh.NewOption("feat - a new feature ", "feat"),
				huh.NewOption("fix  - a bug fix ", "fix"),
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
			name:      "empty prefixes returns default with useEmojis false",
			prefixes:  Prefixes{},
			useEmojis: false,
			want: []huh.Option[string]{
				huh.NewOption("feat     - a new feature ", "feat"),
				huh.NewOption("fix      - a bug fix ", "fix"),
				huh.NewOption("build    - changes that affect the build system or external dependencies ", "build"),
				huh.NewOption("chore    - changes to the build process or auxiliary tools and libraries ", "chore"),
				huh.NewOption("ci       - changes to our CI configuration files and scripts ", "ci"),
				huh.NewOption("docs     - documentation only changes ", "docs"),
				huh.NewOption("perf     - a code change that improves performance ", "perf"),
				huh.NewOption("refactor - a code change that neither fixes a bug nor adds a feature ", "refactor"),
				huh.NewOption("revert   - reverts a previous commit ", "revert"),
				huh.NewOption("style    - changes that do not affect the meaning of the code ", "style"),
				huh.NewOption("test     - adding missing tests or correcting existing tests ", "test"),
			},
		},
		{
			name:      "empty prefixes with useEmojis true returns formatted defaults",
			prefixes:  Prefixes{},
			useEmojis: true,
			want: []huh.Option[string]{
				huh.NewOption("feat     - a new feature ", "feat"),
				huh.NewOption("fix      - a bug fix ", "fix"),
				huh.NewOption("build    - changes that affect the build system or external dependencies ", "build"),
				huh.NewOption("chore    - changes to the build process or auxiliary tools and libraries ", "chore"),
				huh.NewOption("ci       - changes to our CI configuration files and scripts ", "ci"),
				huh.NewOption("docs     - documentation only changes ", "docs"),
				huh.NewOption("perf     - a code change that improves performance ", "perf"),
				huh.NewOption("refactor - a code change that neither fixes a bug nor adds a feature ", "refactor"),
				huh.NewOption("revert   - reverts a previous commit ", "revert"),
				huh.NewOption("style    - changes that do not affect the meaning of the code ", "style"),
				huh.NewOption("test     - adding missing tests or correcting existing tests ", "test"),
			},
		},
		{
			name: "prefixes with emojis enabled",
			prefixes: Prefixes{
				{T: "feat", D: "a new feature", E: emojiPtr("✨")},
				{T: "fix", D: "a bug fix", E: emojiPtr("🐛")},
			},
			useEmojis: true,
			want: []huh.Option[string]{
				huh.NewOption("feat - a new feature ✨ ", "feat ✨"),
				huh.NewOption("fix  - a bug fix 🐛 ", "fix 🐛"),
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
				huh.NewOption("feat - a new feature ", "feat"),
				huh.NewOption("fix  - a bug fix ", "fix"),
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
				huh.NewOption("feat - a new feature ", "feat"),
				huh.NewOption("fix  - a bug fix ", "fix"),
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
				huh.NewOption("feat - a new feature ", "feat"),
				huh.NewOption("fix  - a bug fix ", "fix"),
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
				huh.NewOption("feat - a new feature ✨ ", "feat ✨"),
				huh.NewOption("fix  - a bug fix ", "fix"),
				huh.NewOption("docs - documentation 📚 ", "docs 📚"),
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

func TestGetDefaultPrefixOptions(t *testing.T) {
	expected := []huh.Option[string]{
		huh.NewOption("feat     - a new feature ", "feat"),
		huh.NewOption("fix      - a bug fix ", "fix"),
		huh.NewOption("build    - changes that affect the build system or external dependencies ", "build"),
		huh.NewOption("chore    - changes to the build process or auxiliary tools and libraries ", "chore"),
		huh.NewOption("ci       - changes to our CI configuration files and scripts ", "ci"),
		huh.NewOption("docs     - documentation only changes ", "docs"),
		huh.NewOption("perf     - a code change that improves performance ", "perf"),
		huh.NewOption("refactor - a code change that neither fixes a bug nor adds a feature ", "refactor"),
		huh.NewOption("revert   - reverts a previous commit ", "revert"),
		huh.NewOption("style    - changes that do not affect the meaning of the code ", "style"),
		huh.NewOption("test     - adding missing tests or correcting existing tests ", "test"),
	}

	got := GetDefaultPrefixOptions()
	if len(got) != len(expected) {
		t.Errorf("GetDefaultPrefixOptions() returned %d items, want %d", len(got), len(expected))
		return
	}
	for i, option := range got {
		if option.Key != expected[i].Key || option.Value != expected[i].Value {
			t.Errorf("GetDefaultPrefixOptions()[%d] = {Key: %q, Value: %q}, want {Key: %q, Value: %q}",
				i, option.Key, option.Value, expected[i].Key, expected[i].Value)
		}
	}
}

func TestGetDefaultPrefixOptionsWithEmojis(t *testing.T) {
	tests := []struct {
		name      string
		useEmojis bool
	}{
		{
			name:      "with emojis disabled returns simple format",
			useEmojis: false,
		},
		{
			name:      "with emojis enabled returns formatted with spacing",
			useEmojis: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := defaultPrefixData.OptionsWithEmojis(tt.useEmojis)
			if len(got) != 11 { // We expect 11 default prefixes
				t.Errorf("GetDefaultPrefixOptionsWithEmojis() returned %d items, want 11", len(got))
				return
			}

			// Check that all items have the expected value (type without emoji)
			expectedTypes := []string{"feat", "fix", "build", "chore", "ci", "docs", "perf", "refactor", "revert", "style", "test"}
			for i, option := range got {
				if option.Value != expectedTypes[i] {
					t.Errorf("GetDefaultPrefixOptionsWithEmojis()[%d].Value = %q, want %q", i, option.Value, expectedTypes[i])
				}

				// Check that Key contains the type and description
				if !contains(option.Key, expectedTypes[i]) {
					t.Errorf("GetDefaultPrefixOptionsWithEmojis()[%d].Key = %q should contain %q", i, option.Key, expectedTypes[i])
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
