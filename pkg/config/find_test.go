package config

import (
	"testing"

	"github.com/spf13/afero"
)

func TestFindConfigFile(t *testing.T) {
	t.Run("no config file", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		expepcted := ""
		got, err := FindConfigFile(fs)
		assertEqual(t, expepcted, got)
		assertIsError(t, err)
	})
	t.Run("config file in current directory", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		fs.Create("./.meteor.json")
		expected := ".meteor.json"
		got, err := FindConfigFile(fs)
		assertEqual(t, expected, got)
		assertIsNotError(t, err)
	})
}

func assertEqual(t testing.TB, expected, got string) {
	t.Helper()
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func assertIsNotError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Errorf("expected no error, but got %v", got)
	}
}

func assertIsError(t testing.TB, got error) {
	t.Helper()
	if got == nil {
		t.Errorf("expected an error, but got nil")
	}
}
