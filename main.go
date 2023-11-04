package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
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
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
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
