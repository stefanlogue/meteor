package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	flag "github.com/spf13/pflag"
)

type Commit struct {
	Board     string
	Type      string
	Scope     string
	Message   string
	Body      string
	Coauthors []string
}

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

	var commit Commit
	mainForm := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("meteor").
				Description("A command line tool for generating conventional commit messages."),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Board").
				Description("Select the board for this commit").
				Options(boards...).
				Value(&commit.Board),
		).WithHideFunc(func() bool {
			return len(boards) < 1
		}),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Coauthors").
				Description("Select any coauthors for this commit").
				Options(coauthors...).
				Value(&commit.Coauthors),
		).WithHideFunc(func() bool {
			return len(coauthors) < 1
		}),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Type").
				Description("Select the type of change that you're committing").
				Options(prefixes...).
				Validate(func(s string) error {
					commit.Message = s + ": "
					return nil
				}).
				Value(&commit.Type),
		),
		huh.NewGroup(
			huh.NewInput().
				Value(&commit.Message).
				Title("Message").
				CharLimit(48),
			huh.NewText().
				Value(&commit.Body).
				Title("Body").
				Lines(8),
		),
	)

	err = mainForm.Run()
	if err != nil {
		fail("Error: %s", err)
	}

	if err != nil {
		fail("Error: %s", err)
	}
}

func fail(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
