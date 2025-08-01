package config

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

const (
	homeDir = "/home/user"
)

func TestFindConfigFile(t *testing.T) {
	t.Run("no config file", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		currentDir := "/home/user/project"
		expected := ""
		got, err := FindConfigFile(fs,
			func() (string, error) { return currentDir, nil },
			func() (string, error) { return homeDir, nil },
		)
		assertEqual(t, expected, got)
		assertIsError(t, err)
	})
	t.Run("config file in current directory", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		currentDir := ""
		writeErr := afero.WriteFile(fs, ".meteor.json", []byte("{}"), 0644)
		assertIsNotError(t, writeErr)
		expected := ".meteor.json"
		got, err := FindConfigFile(fs,
			func() (string, error) { return currentDir, nil },
			func() (string, error) { return homeDir, nil },
		)
		assertEqual(t, expected, got)
		assertIsNotError(t, err)
	})
	t.Run("config file in nested directory", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		currentDir := "/home/user/project"
		nestedDir := filepath.Join(currentDir, "nested")
		configPath := filepath.Join(nestedDir, ".meteor.json")
		fs.MkdirAll(nestedDir, 0755)
		writeErr := afero.WriteFile(fs, configPath, []byte("{}"), 0644)
		assertIsNotError(t, writeErr)
		expected := ""
		got, err := FindConfigFile(fs,
			func() (string, error) { return currentDir, nil },
			func() (string, error) { return homeDir, nil },
		)
		assertEqual(t, expected, got)
		assertIsError(t, err)
	})
	t.Run("config file in parent directory", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		currentDir := "/home/user/project"
		configPath := filepath.Join(homeDir, ".meteor.json")
		fs.MkdirAll(currentDir, 0755)
		content := "{}"
		writeErr := afero.WriteFile(fs, configPath, []byte(content), 0644)
		assertIsNotError(t, writeErr)
		expected := configPath
		got, err := FindConfigFile(fs,
			func() (string, error) { return currentDir, nil },
			func() (string, error) { return homeDir, nil },
		)
		assertEqual(t, expected, got)
		assertIsNotError(t, err)
	})
	t.Run("config file in xdg config dir", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		currentDir := "/home/user/project"
		configPath := filepath.Join(homeDir, ".config/meteor/config.json")
		fs.MkdirAll(filepath.Dir(configPath), 0755)
		content := "{}"
		writeErr := afero.WriteFile(fs, configPath, []byte(content), 0644)
		assertIsNotError(t, writeErr)
		expected := configPath
		got, err := FindConfigFile(fs,
			func() (string, error) { return currentDir, nil },
			func() (string, error) { return homeDir, nil },
		)
		assertEqual(t, expected, got)
		assertIsNotError(t, err)
	})
	t.Run("WD not in home directory", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		currentDir := "/var/www/hosts/project"
		configPath := filepath.Join(homeDir, ".config/meteor/config.json")
		fs.MkdirAll(filepath.Dir(configPath), 0755)
		content := "{}"
		writeErr := afero.WriteFile(fs, configPath, []byte(content), 0644)
		assertIsNotError(t, writeErr)
		expected := configPath
		got, err := FindConfigFile(fs,
			func() (string, error) { return currentDir, nil },
			func() (string, error) { return homeDir, nil },
		)
		assertEqual(t, expected, got)
		assertIsNotError(t, err)
	})
	t.Run("error in current directory function", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		expected := ""
		got, err := FindConfigFile(fs,
			func() (string, error) { return "", fmt.Errorf("error getting cwd") },
			func() (string, error) { return homeDir, nil },
		)
		assertEqual(t, expected, got)
		assertIsError(t, err)
	})
	t.Run("error in home directory function", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		currentDir := "/home/user/project"
		expected := ""
		got, err := FindConfigFile(fs,
			func() (string, error) { return currentDir, nil },
			func() (string, error) { return "", fmt.Errorf("error getting home dir") },
		)
		assertEqual(t, expected, got)
		assertIsError(t, err)
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
