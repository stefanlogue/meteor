package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/fatih/color"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
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

// isFlagPassed checks if a flag has been passed
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

var (
	version  = "dev"
	logLevel string
)

func init() {
	flag.BoolP("version", "v", false, "show version")
	flag.StringVarP(&logLevel, "log-level", "L", "info", "Log level (debug, info, warn, error, fatal, panic)")
	flag.Parse()
	if isFlagPassed("version") {
		fmt.Printf("meteor version %s\n", version)
		os.Exit(0)
	}

	programLevel, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal("invalid log level", "error", err)
	}

	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
	})

	logger.SetLevel(programLevel)
	log.SetDefault(logger)
}

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

	prefixes, coauthors, boards, showIntro, err := loadConfig()
	if err != nil {
		fail("Error: %s", err)
	}

	var newCommit Commit
	theme := huh.ThemeCatppuccin()
	if showIntro {
		introForm := huh.NewForm(
			huh.NewGroup(
				splashScreen(),
			),
		)
		if err := introForm.Run(); err != nil {
			fail("Error: %s", err)
		}
	}
	if len(boards) > 0 {
		boardForm := huh.NewForm(
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

	if len(newCommit.Board) > 0 && newCommit.Board != "NONE" {
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

	if len(newCommit.Board) > 0 && newCommit.Board != "NONE" {
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

	doesWantToCommit := true
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
	).WithKeyMap(&huh.KeyMap{
		Quit: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
		Text: huh.TextKeyMap{
			Next:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "next")),
			NewLine: key.NewBinding(key.WithKeys("alt+enter", "ctrl+j"), key.WithHelp("alt+enter / ctrl+j", "new line")),
			Editor:  key.NewBinding(key.WithKeys("ctrl+e"), key.WithHelp("ctrl+e", "open editor")),
			Prev:    key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
		},
		Input: huh.InputKeyMap{
			Next: key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter / tab", "next")),
		},
		Confirm: huh.ConfirmKeyMap{
			Toggle: key.NewBinding(key.WithKeys("left", "right", "h", "l"), key.WithHelp("left / right", "toggle")),
			Next:   key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter / tab", "next")),
			Prev:   key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
		},
	}).WithTheme(theme)

	err = messageForm.Run()
	if err != nil {
		fail("Error: %s", err)
	}

	if len(newCommit.Coauthors) > 0 {
		newCommit.Body = newCommit.Body + buildCoauthorString(newCommit.Coauthors)
	}

	args := flag.Args()
	rawCommitCommand, printableCommitCommand := buildCommitCommand(newCommit.Message, newCommit.Body, args)
	if doesWantToCommit {
		err := commit(rawCommitCommand)
		if err != nil {
			writeToClipboard(printableCommitCommand)
			fail(
				"\n%s\n%s\n\n%s\n\n",
				color.RedString(fmt.Sprintf("It looks like the commit failed.\nError: %s", err)),
				color.YellowString("To run it again without going through meteor's wizard, simply run the following command (I've copied it to your clipboard!):"),
				color.BlueString(printableCommitCommand),
			)
		}
	} else {
		writeToClipboard(printableCommitCommand)
		fmt.Printf(
			"\n%s\n\n%s\n%s\n\n",
			color.RedString("Commit aborted."),
			color.YellowString("I've copied the following command to your clipboard, so you can run it again later:"),
			color.BlueString(printableCommitCommand))
	}
}

// writeToClipboard writes a string to the clipboard
func writeToClipboard(s string) {
	clipboard.WriteAll(s)
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

// splashScreen returns a note with a splash screen
func splashScreen() *huh.Note {
	return huh.NewNote().
		Title("meteor").
		Description("A highly customisable command line tool\nfor writing conventional commit messages")
}

// fail prints an error message and exits with a non-zero exit code
func fail(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
