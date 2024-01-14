package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	flag "github.com/spf13/pflag"
)

type Commit struct {
	Board            string
	TicketNumber     string
	Type             string
	Scope            string
	Message          string
	Body             string
	Coauthors        []string
	IsBreakingChange bool
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

	var newCommit Commit
	theme := huh.ThemeCatppuccin()
	if len(boards) > 0 {
		boardForm := huh.NewForm(
			huh.NewGroup(
				splashScreen(),
			),
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Board").
					Description("Select the board for this commit").
					Options(boards...).
					Value(&newCommit.Board),
			).WithHideFunc(func() bool {
				return len(boards) < 1
			}),
		).WithTheme(theme)

		err = boardForm.Run()
		if err != nil {
			fail("Error: %s", err)
		}
	}

	if len(newCommit.Board) > 0 {
		ticketNumber := getGitTicketNumber(newCommit.Board)

		if ticketNumber == "" {
			newCommit.TicketNumber = fmt.Sprintf("%s-", newCommit.Board)
		} else {
			newCommit.TicketNumber = ticketNumber
		}

		ticketNumberForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Ticket number").
					Description("The ticket number associated with this commit").
					CharLimit(10).
					Value(&newCommit.TicketNumber),
			).WithHideFunc(func() bool {
				return len(boards) < 1
			}),
		).WithTheme(theme)

		err = ticketNumberForm.Run()
		if err != nil {
			fail("Error: %s", err)
		}
	}

	mainForm := huh.NewForm(
		huh.NewGroup(
			splashScreen(),
		).WithHideFunc(func() bool {
			return len(boards) > 0
		}),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Type").
				Description("Select the type of change that you're committing").
				Options(prefixes...).
				Value(&newCommit.Type),
			huh.NewConfirm().
				Title("Breaking Change").
				Description("Is this a breaking change?").
				Affirmative("Yes!").
				Negative("Nope.").
				Value(&newCommit.IsBreakingChange),
			huh.NewInput().
				Title("Scope").
				Description("Specify a scope of the change").
				CharLimit(16).
				Value(&newCommit.Scope),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Coauthors").
				Description("Select any coauthors for this commit").
				Options(coauthors...).
				Value(&newCommit.Coauthors),
		).WithHideFunc(func() bool {
			return len(coauthors) < 1
		}),
	).WithTheme(theme)

	err = mainForm.Run()
	if err != nil {
		fail("Error: %s", err)
	}

	if len(newCommit.Board) > 0 {
		if newCommit.IsBreakingChange {
			if len(newCommit.Scope) > 0 {
				newCommit.Message = fmt.Sprintf("%s(%s)!: <%s> ", newCommit.TicketNumber, newCommit.Scope, newCommit.Type)
			} else {
				newCommit.Message = fmt.Sprintf("%s!: <%s> ", newCommit.TicketNumber, newCommit.Type)
			}
		} else {
			if len(newCommit.Scope) > 0 {
				newCommit.Message = fmt.Sprintf("%s(%s): <%s> ", newCommit.TicketNumber, newCommit.Scope, newCommit.Type)
			} else {
				newCommit.Message = fmt.Sprintf("%s: <%s> ", newCommit.TicketNumber, newCommit.Type)
			}
		}
	} else {
		if newCommit.IsBreakingChange {
			if len(newCommit.Scope) > 0 {
				newCommit.Message = fmt.Sprintf("%s(%s)!: ", newCommit.Type, newCommit.Scope)
			} else {
				newCommit.Message = fmt.Sprintf("%s!: ", newCommit.Type)
			}
		} else {
			if len(newCommit.Scope) > 0 {
				newCommit.Message = fmt.Sprintf("%s(%s): ", newCommit.Type, newCommit.Scope)
			} else {
				newCommit.Message = fmt.Sprintf("%s: ", newCommit.Type)
			}
		}
	}

	var doesWantToCommit bool
	messageForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(&newCommit.Message).
				Title("Message").
				CharLimit(48),
			huh.NewText().
				Value(&newCommit.Body).
				Title("Body").
				Lines(8),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Ready to commit?").
				Affirmative("Yes!").
				Negative("No.").
				Value(&doesWantToCommit),
		),
	).WithTheme(theme)

	err = messageForm.Run()
	if err != nil {
		fail("Error: %s", err)
	}

	if len(newCommit.Coauthors) > 0 {
		newCommit.Body = newCommit.Body + buildCoauthorString(newCommit.Coauthors)
	}

	if doesWantToCommit {
		err := commit(newCommit.Message, newCommit.Body)
		if err != nil {
			fail("Error: %s", err)
		}
	}
}

// buildCoauthorString takes a slice of selected coauthors and returns a formatted
// string which Github recognises
func buildCoauthorString(coauthors []string) string {
	s := `


	`

	for _, coauthor := range coauthors {
		if coauthor == "none" {
			return ""
		}
		s += fmt.Sprintf("\nCo-authored-by: %s", coauthor)
	}
	return s
}

func splashScreen() *huh.Note {
	return huh.NewNote().
		Title("meteor").
		Description("A highly customisable command line tool\nfor writing conventional commit messages")
}

func fail(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
