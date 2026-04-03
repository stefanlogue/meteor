package config_test

import (
	"testing"

	"github.com/stefanlogue/meteor/pkg/config"
)

func TestBuildCoAuthorString(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		coauthors []string
		want      string
	}{
		{
			name:      "empty slice returns empty string",
			coauthors: []string{},
			want:      "\n\n\n\t",
		},
		{
			name:      "none option returns empty string",
			coauthors: []string{"none"},
			want:      "",
		},
		{
			name:      "single coauthor",
			coauthors: []string{"John Doe <john@example.com>"},
			want:      "\n\n\n\t\nCo-authored-by: John Doe <john@example.com>",
		},
		{
			name: "multiple coauthors",
			coauthors: []string{
				"John Doe <john@example.com>",
				"Jane Smith <jane@example.com>",
			},
			want: "\n\n\n\t\nCo-authored-by: John Doe <john@example.com>\nCo-authored-by: Jane Smith <jane@example.com>",
		},
		{
			name: "none in middle of list",
			coauthors: []string{
				"John Doe <john@example.com>",
				"none",
				"Jane Smith <jane@example.com>",
			},
			want: "",
		},
		{
			name: "multiple coauthors with different formats",
			coauthors: []string{
				"John Doe <john@example.com>",
				"Jane Smith <jane@example.com>",
				"Bob Johnson <bob@example.com>",
			},
			want: "\n\n\n\t\nCo-authored-by: John Doe <john@example.com>\nCo-authored-by: Jane Smith <jane@example.com>\nCo-authored-by: Bob Johnson <bob@example.com>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := config.BuildCoAuthorString(tt.coauthors)
			if got != tt.want {
				t.Errorf("BuildCoAuthorString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildCoauthorString_NoneShortCircuits(t *testing.T) {
	// Verify that "none" immediately returns empty string
	// even if there are other coauthors after it
	coauthors := []string{
		"John Doe <john@example.com>",
		"none",
		"This should not appear <test@example.com>",
	}

	got := config.BuildCoAuthorString(coauthors)
	if got != "" {
		t.Errorf("buildCoauthorString() with 'none' = %q, want empty string", got)
	}
}
