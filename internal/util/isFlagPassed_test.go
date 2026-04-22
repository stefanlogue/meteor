package util_test

import (
	"testing"

	flag "github.com/spf13/pflag"
	"github.com/stefanlogue/meteor/internal/util"
)

func TestIsFlagPassed(t *testing.T) {
	original := flag.CommandLine
	t.Cleanup(func() {
		flag.CommandLine = original
	})

	tests := []struct {
		name     string
		flagName string
		setFlag  bool
		want     bool
	}{
		{
			name:     "version flag is not passed",
			flagName: "version",
			want:     false,
		},
		{
			name:     "version flag is passed",
			flagName: "version",
			setFlag:  true,
			want:     true,
		},
		{
			name:     "non-existent flag",
			flagName: "non-existent-flag",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
			flag.Bool("version", false, "")
			if tt.setFlag {
				if err := flag.CommandLine.Set(tt.flagName, "true"); err != nil {
					t.Fatalf("failed setting flag %q: %v", tt.flagName, err)
				}
			}

			got := util.IsFlagPassed(tt.flagName)
			if got != tt.want {
				t.Errorf("IsFlagPassed() = %v, want %v", got, tt.want)
			}
		})
	}
}
