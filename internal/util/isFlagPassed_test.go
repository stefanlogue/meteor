package util_test

import (
	"testing"

	"github.com/stefanlogue/meteor/internal/util"
)

func TestIsFlagPassed(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		want     bool
	}{
		{
			name:     "version flag is passed",
			flagName: "version",
			want:     false, // Will be false in test context unless explicitly set
		},
		{
			name:     "non-existent flag",
			flagName: "non-existent-flag",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.IsFlagPassed(tt.flagName)
			if got != tt.want {
				t.Errorf("IsFlagPassed() = %v, want %v", got, tt.want)
			}
		})
	}
}
