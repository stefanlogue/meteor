package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/charmbracelet/log"
	"github.com/fatih/color"
	"github.com/stefanlogue/meteor/internal/util"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"github.com/spf13/afero"
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
	version            = "dev"
	debugMode          bool
	skipIntro          bool
	skipBreakingChange bool
	addAll             bool
	FS                 afero.Fs     = afero.NewOsFs()
	AFS                *afero.Afero = &afero.Afero{Fs: FS}
)

const (
	AsGitEditor = "as-git-editor"
	ErrorString = "Error: %s"
	ShiftTab    = "shift+tab"
)

func init() {
	flag.BoolP("version", "v", false, "show version")
	flag.BoolP(AsGitEditor, "e", false, "used as GIT_EDITOR")
	flag.BoolVarP(&skipIntro, "skip-intro", "s", false, "skip intro splash")
	flag.BoolVarP(&debugMode, "debug", "D", false, "enable debug mode")
	flag.BoolVarP(&skipBreakingChange, "skip-breaking-change", "b", false, "skip breaking change prompt")
	flag.BoolVarP(&addAll, "add-all", "a", false, "add all files to commit")
	flag.Parse()
	if isFlagPassed("version") {
		fmt.Printf("meteor version %s\n", version)
		os.Exit(0)
	}

	programLevel := log.InfoLevel
	if debugMode {
		programLevel = log.DebugLevel
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
	gitPath, err := getGitPath()
	if err != nil {
		fail(ErrorString, err)
	}

	gitRoot, err := findGitDir(gitPath)
	if err != nil {
		fail(ErrorString, err)
	}

	if err := os.Chdir(gitRoot); err != nil {
		fail("Could not change directory: %s", err)
	}

	config, err := loadConfig(AFS)
	if err != nil {
		fail(ErrorString, err)
	}

	if !config.AddAll && addAll {
		config.AddAll = true
	}

	var newCommit Commit
	theme := huh.ThemeCatppuccin()
	if config.ShowIntro && (isFlagPassed("skip-intro") && !skipIntro) {
		introForm := huh.NewForm(
			huh.NewGroup(
				splashScreen(),
			),
		)
		if err := introForm.Run(); err != nil {
			fail(ErrorString, err)
		}
	}
	if len(config.Boards) > 0 {
		boardForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Board").
					Description("Select the board for this commit").
					Options(config.Boards...).
					Value(&newCommit.Board),
			).WithHideFunc(func() bool {
				return len(config.Boards) < 1
			}),
		).WithTheme(theme)

		err = boardForm.Run()
		if err != nil {
			fail(ErrorString, err)
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
					CharLimit(24).
					Value(&newCommit.TicketNumber),
			).WithHideFunc(func() bool {
				return len(config.Boards) < 1
			}),
		).WithTheme(theme)

		err = ticketNumberForm.Run()
		if err != nil {
			fail(ErrorString, err)
		}
	}

	typeInput := huh.NewSelect[string]().
		Title("Type").
		Description("Select the type of change that you're committing").
		Options(config.Prefixes...).
		Value(&newCommit.Type)

	// if the user has specified scopes in their config, use a select input, otherwise use a text input
	var scopeInput huh.Field
	if len(config.Scopes) > 0 {
		scopeInput = huh.NewSelect[string]().
			Title("Scope").
			Description("Choose a scope for the changes").
			Options(config.Scopes...).
			Value(&newCommit.Scope)
	} else {
		scopeInput = huh.NewInput().
			Title("Scope").
			Description("Specify a scope of the changes").
			CharLimit(16).
			Value(&newCommit.Scope)
	}

	// if the user has specified for asking breaking change, add a confirm input to the main group
	var mainGroup *huh.Group
	if skipBreakingChange {
		mainGroup = huh.NewGroup(typeInput, scopeInput)
	} else {
		mainGroup = huh.NewGroup(
			typeInput,
			huh.NewConfirm().
				Title("Breaking Change").
				Description("Is this a breaking change?").
				Affirmative("Yes!").
				Negative("Nope.").
				Value(&newCommit.IsBreakingChange),
			scopeInput,
		)
	}
	coAuthors := config.Coauthors
	if config.ReadContributorsFromGit {
		additional, err := getComitters([]string{})
		if err != nil {
			fail(ErrorString, err)
		} else {
			for _, s := range additional {
				coAuthors = append(coAuthors, huh.NewOption(s, s))
			}
		}
	}
	if len(coAuthors) > 0 {
		coAuthors = util.PrependItem(coAuthors, huh.NewOption("no coauthors", "none"))
	}
	mainForm := huh.NewForm(
		mainGroup,
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Coauthors").
				Description("Select any coauthors for this commit").
				Options(coAuthors...).
				Value(&newCommit.Coauthors),
		).WithHideFunc(func() bool {
			return len(coAuthors) < 1
		}),
	).WithTheme(theme)

	err = mainForm.Run()
	if err != nil {
		fail(ErrorString, err)
	}

	var tmpl *template.Template
	if len(newCommit.Board) > 0 && newCommit.Board != "NONE" {
		tmpl = template.Must(template.New("message").Parse(config.MessageWithTicketTemplate))
	} else {
		tmpl = template.Must(template.New("message").Parse(config.MessageTemplate))
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, newCommit)
	if err != nil {
		fail(ErrorString, err)
	}
	newCommit.Message = buf.String()

	doesWantToCommit := true
	messageForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(&newCommit.Message).
				Title("Message").
				CharLimit(config.CommitTitleCharLimit),
			huh.NewText().
				Value(&newCommit.Body).
				Title("Body").
				CharLimit(config.CommitBodyCharLimit).
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
			Prev:    key.NewBinding(key.WithKeys(ShiftTab), key.WithHelp(ShiftTab, "back")),
		},
		Input: huh.InputKeyMap{
			Next: key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter / tab", "next")),
		},
		Confirm: huh.ConfirmKeyMap{
			Toggle: key.NewBinding(key.WithKeys("left", "right", "h", "l"), key.WithHelp("left / right", "toggle")),
			Prev:   key.NewBinding(key.WithKeys(ShiftTab), key.WithHelp(ShiftTab, "back")),
			Submit: key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter / tab", "submit")),
		},
	}).WithTheme(theme)

	err = messageForm.Run()
	if err != nil {
		fail(ErrorString, err)
	}

	if config.CommitBodyLineLength >= minimumCommitBodyLineLength {
		newCommit.Body = wordWrap(newCommit.Body, config.CommitBodyLineLength)
	}

	if len(newCommit.Coauthors) > 0 {
		newCommit.Body = newCommit.Body + buildCoauthorString(newCommit.Coauthors)
	}

	args := flag.Args()

	var commitFile string

	// If we're operating in GIT_EDITOR="meteor --as-git-editor" mode, the first argument is the path (.git/COMMIT_EDITMSG)
	// where we should write the git commit message, so we shift that from args before constructing the end-user command line
	if isFlagPassed(AsGitEditor) {
		commitFile = args[0]
		args = args[1:]
	}

	if config.AddAll {
		rawCommitCommand, printableCommitCommand := buildCommitCommand(newCommit.Message, newCommit.Body, args)
		err = commit(rawCommitCommand)
		if err != nil {
			writeToClipboard(printableCommitCommand)
			fail(
				"\n%s\n%s\n\n%s\n\n",
				color.RedString(fmt.Sprintf("It looks like the add all failed.\nError: %s", err)),
				color.YellowString("To run it again without going through meteor's wizard, simply run the following command (I've copied it to your clipboard!):"),
				color.BlueString(printableCommitCommand),
			)
			return
		}
	}
	rawCommitCommand, printableCommitCommand := buildCommitCommand(newCommit.Message, newCommit.Body, args)

	if isFlagPassed(AsGitEditor) {
		// We intent to do the commit
		if doesWantToCommit {
			// Write the commit message file (.git/COMMIT_EDITMSG) in same format as git would have,
			// the message, a blank line, and a body - if body is empty, trailing newlines will be removed

			if err := os.WriteFile(commitFile, bytes.TrimRight([]byte(newCommit.Message+"\n\n"+newCommit.Body), "/n"), os.FileMode(os.O_WRONLY)); err != nil {
				// In case of failure, give the regular error-ish output to the end-user so no inputs are lost
				writeToClipboard(printableCommitCommand)

				fail(
					"\n%s\n%s\n\n%s\n\n",
					color.RedString(fmt.Sprintf("It looks like the commit failed.\nError: %s", err)),
					color.YellowString("To run it again without going through meteor's wizard, simply run the following command (I've copied it to your clipboard!):"),
					color.BlueString(printableCommitCommand),
				)

				return
			}

			// we wrote the commit message file, nothing left for us to do, success!

			return
		}

		// end-user decided to abort the commit, which mean we don't write the git commit message file (.git/COMMIT_EDITMSG)
		// which will make git abort the operation

		writeToClipboard(printableCommitCommand)
		fmt.Printf(
			"\n%s\n\n%s\n%s\n\n",
			color.RedString("Commit aborted."),
			color.YellowString("I've copied the following command to your clipboard, so you can run it again later:"),
			color.BlueString(printableCommitCommand))

		return
	}

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
	if err := clipboard.WriteAll(s); err != nil {
		fail("Failed to copy to clipboard: %s", err)
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
