package config

import (
	"testing"
)

func TestPrefixes_Options(t *testing.T) {
	tests := []struct {
		name     string
		prefixes Prefixes
		want     int
		wantType string
	}{
		{
			name:     "empty prefixes returns defaults",
			prefixes: Prefixes{},
			want:     len(DefaultSelectablePrefixes),
			wantType: "feat",
		},
		{
			name:     "nil prefixes returns defaults",
			prefixes: nil,
			want:     len(DefaultSelectablePrefixes),
			wantType: "feat",
		},
		{
			name: "custom prefixes returns custom options",
			prefixes: Prefixes{
				{T: "custom", D: "a custom prefix"},
			},
			want:     1,
			wantType: "custom",
		},
		{
			name: "multiple custom prefixes",
			prefixes: Prefixes{
				{T: "feature", D: "a new feature"},
				{T: "bugfix", D: "a bug fix"},
				{T: "hotfix", D: "a critical fix"},
			},
			want:     3,
			wantType: "feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.prefixes.Options()

			if len(got) != tt.want {
				t.Errorf("Options() returned %d items, want %d", len(got), tt.want)
			}

			if len(got) > 0 {
				// Check that the first option has the expected value
				firstValue := got[0].Value
				if firstValue != tt.wantType {
					t.Errorf("Options() first value = %v, want %v", firstValue, tt.wantType)
				}
			}
		})
	}
}

func TestPrefixes_Options_Format(t *testing.T) {
	prefixes := Prefixes{
		{T: "feat", D: "a new feature"},
	}

	got := prefixes.Options()

	if len(got) != 1 {
		t.Fatalf("Options() returned %d items, want 1", len(got))
	}

	// Verify the option format matches "type - description"
	expectedKey := "feat - a new feature"
	expectedValue := "feat"

	if got[0].Key != expectedKey {
		t.Errorf("Options() key = %q, want %q", got[0].Key, expectedKey)
	}

	if got[0].Value != expectedValue {
		t.Errorf("Options() value = %q, want %q", got[0].Value, expectedValue)
	}
}

func TestPrefixes_Strings(t *testing.T) {
	tests := []struct {
		name     string
		prefixes Prefixes
		want     []string
	}{
		{
			name:     "empty prefixes returns defaults",
			prefixes: Prefixes{},
			want:     DefaultPrefixes,
		},
		{
			name:     "nil prefixes returns defaults",
			prefixes: nil,
			want:     DefaultPrefixes,
		},
		{
			name: "custom prefixes returns custom strings",
			prefixes: Prefixes{
				{T: "custom", D: "a custom prefix"},
			},
			want: []string{"custom"},
		},
		{
			name: "multiple custom prefixes",
			prefixes: Prefixes{
				{T: "feature", D: "a new feature"},
				{T: "bugfix", D: "a bug fix"},
				{T: "hotfix", D: "a critical fix"},
			},
			want: []string{"feature", "bugfix", "hotfix"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.prefixes.Strings()

			if len(got) != len(tt.want) {
				t.Errorf("Strings() returned %d items, want %d", len(got), len(tt.want))
				return
			}

			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("Strings()[%d] = %v, want %v", i, v, tt.want[i])
				}
			}
		})
	}
}

func TestPrefixes_Strings_IgnoresDescriptions(t *testing.T) {
	prefixes := Prefixes{
		{T: "feat", D: "description 1"},
		{T: "fix", D: "description 2"},
	}

	got := prefixes.Strings()
	want := []string{"feat", "fix"}

	if len(got) != len(want) {
		t.Fatalf("Strings() returned %d items, want %d", len(got), len(want))
	}

	for i, v := range got {
		if v != want[i] {
			t.Errorf("Strings()[%d] = %v, want %v", i, v, want[i])
		}
	}
}

func TestDefaultPrefixes_Length(t *testing.T) {
	if len(DefaultPrefixes) != 11 {
		t.Errorf("DefaultPrefixes length = %d, want 11", len(DefaultPrefixes))
	}
}

func TestDefaultSelectablePrefixes_Length(t *testing.T) {
	if len(DefaultSelectablePrefixes) != 11 {
		t.Errorf("DefaultSelectablePrefixes length = %d, want 11", len(DefaultSelectablePrefixes))
	}
}

func TestDefaultPrefixes_MatchesSelectablePrefixes(t *testing.T) {
	if len(DefaultPrefixes) != len(DefaultSelectablePrefixes) {
		t.Errorf("DefaultPrefixes and DefaultSelectablePrefixes have different lengths: %d vs %d",
			len(DefaultPrefixes), len(DefaultSelectablePrefixes))
	}

	// Verify each default prefix has a corresponding selectable option
	for i, prefix := range DefaultPrefixes {
		if DefaultSelectablePrefixes[i].Value != prefix {
			t.Errorf("DefaultSelectablePrefixes[%d].Value = %q, want %q",
				i, DefaultSelectablePrefixes[i].Value, prefix)
		}
	}
}
