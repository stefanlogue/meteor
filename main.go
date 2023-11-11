package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	flag "github.com/spf13/pflag"
)

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

var version = "dev"

func main() {
	flag.BoolP("version", "v", false, "show version")
	flag.Parse()
	if isFlagPassed("version") {
		fmt.Printf("meteor version %s\n", version)
		os.Exit(0)
	}

	if err := checkGitInPath(); err != nil {
		fail("Error: %s", err)
	}

	gitRoot, err := findGitDir()
	if err != nil {
		fail("Error: %s", err)
	}

	if err := os.Chdir(gitRoot); err != nil {
		fail("Could not change directory: %s", err)
	}

	prefixes, coauthors, boards, err := loadConfig()
	if err != nil {
		fail("Error: %s", err)
	}

	m := newModel(boards, prefixes, coauthors)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fail("Error: %s", err)
	}
	fmt.Println("")

	if !m.Finished() {
		fail("Aborted")
	}

	msg, body := m.CommitMessage()
	if err := commit(msg, body); err != nil {
		fail("Error with commit: %s", err)
	}
}

func fail(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
