package config

import (
	"testing"
)

func TestScopes_Options(t *testing.T) {
	tests := []struct {
		name       string
		scopes     Scopes
		wantLen    int
		wantNil    bool
		wantFirst  string
		wantSecond string
	}{
		{
			name:    "empty scopes returns nil",
			scopes:  Scopes{},
			wantNil: true,
		},
		{
			name:    "nil scopes returns nil",
			scopes:  nil,
			wantNil: true,
		},
		{
			name: "single scope returns none + scope",
			scopes: Scopes{
				{Name: "api"},
			},
			wantLen:    2,
			wantFirst:  "", // "none" option has empty value
			wantSecond: "api",
		},
		{
			name: "multiple scopes returns none + all scopes",
			scopes: Scopes{
				{Name: "api"},
				{Name: "ui"},
				{Name: "db"},
			},
			wantLen:    4,
			wantFirst:  "",
			wantSecond: "api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.scopes.Options()

			if tt.wantNil {
				if got != nil {
					t.Errorf("Options() = %v, want nil", got)
				}
				return
			}

			if len(got) != tt.wantLen {
				t.Errorf("Options() returned %d items, want %d", len(got), tt.wantLen)
				return
			}

			// Check first option is "none" with empty value
			if got[0].Value != tt.wantFirst {
				t.Errorf("Options()[0].Value = %q, want %q", got[0].Value, tt.wantFirst)
			}

			if got[0].Key != "none" {
				t.Errorf("Options()[0].Key = %q, want %q", got[0].Key, "none")
			}

			// Check second option if exists
			if len(got) > 1 && got[1].Value != tt.wantSecond {
				t.Errorf("Options()[1].Value = %q, want %q", got[1].Value, tt.wantSecond)
			}
		})
	}
}

func TestScopes_Options_OrderPreserved(t *testing.T) {
	scopes := Scopes{
		{Name: "frontend"},
		{Name: "backend"},
		{Name: "database"},
	}

	got := scopes.Options()

	if len(got) != 4 {
		t.Fatalf("Options() returned %d items, want 4", len(got))
	}

	expectedOrder := []string{"", "frontend", "backend", "database"}
	for i, expected := range expectedOrder {
		if got[i].Value != expected {
			t.Errorf("Options()[%d].Value = %q, want %q", i, got[i].Value, expected)
		}
	}

	// Verify "none" is first
	if got[0].Key != "none" {
		t.Errorf("Options()[0].Key = %q, want %q", got[0].Key, "none")
	}

	// Verify other keys match their values
	for i := 1; i < len(got); i++ {
		if got[i].Key != got[i].Value {
			t.Errorf("Options()[%d].Key = %q, want %q", i, got[i].Key, got[i].Value)
		}
	}
}

func TestScopes_Strings(t *testing.T) {
	tests := []struct {
		name    string
		scopes  Scopes
		want    []string
		wantNil bool
	}{
		{
			name:    "empty scopes returns nil",
			scopes:  Scopes{},
			wantNil: true,
		},
		{
			name:    "nil scopes returns nil",
			scopes:  nil,
			wantNil: true,
		},
		{
			name: "single scope returns single string",
			scopes: Scopes{
				{Name: "api"},
			},
			want: []string{"api"},
		},
		{
			name: "multiple scopes returns all strings",
			scopes: Scopes{
				{Name: "api"},
				{Name: "ui"},
				{Name: "db"},
			},
			want: []string{"api", "ui", "db"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.scopes.Strings()

			if tt.wantNil {
				if got != nil {
					t.Errorf("Strings() = %v, want nil", got)
				}
				return
			}

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

func TestScopes_Strings_NoNoneOption(t *testing.T) {
	scopes := Scopes{
		{Name: "api"},
		{Name: "ui"},
	}

	got := scopes.Strings()

	// Verify "none" is NOT included in Strings() output
	for i, v := range got {
		if v == "" || v == "none" {
			t.Errorf("Strings()[%d] = %q, should not contain empty or 'none'", i, v)
		}
	}

	// Verify we only have the actual scope names
	want := []string{"api", "ui"}
	if len(got) != len(want) {
		t.Errorf("Strings() length = %d, want %d", len(got), len(want))
	}
}

func TestScopes_Strings_OrderPreserved(t *testing.T) {
	scopes := Scopes{
		{Name: "zulu"},
		{Name: "alpha"},
		{Name: "bravo"},
	}

	got := scopes.Strings()
	want := []string{"zulu", "alpha", "bravo"}

	if len(got) != len(want) {
		t.Fatalf("Strings() returned %d items, want %d", len(got), len(want))
	}

	for i, v := range got {
		if v != want[i] {
			t.Errorf("Strings()[%d] = %q, want %q (order should be preserved)", i, v, want[i])
		}
	}
}
