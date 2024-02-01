package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const maxGitRecursion = 32

// checkGitInPath checks if git is in PATH and returns an error if not
func checkGitInPath() error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git not found in PATH: %w", err)
	}
	return nil
}

// findGitDir returns the root of the git repository
func findGitDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

// getGitTicketNumber returns the most recent ticket number from the current git branch
func getGitTicketNumber(board string) string {
	cmd := exec.Command("git", "branch", "--show-current")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	match, _ := regexp.MatchString(fmt.Sprintf(`(?i)^%s-\d{4,}`, board), string(out))
	if !match {
		cmd = exec.Command("git", "log", "-1", "--grep", board, "--oneline", "--format=%s")
		out, err = cmd.Output()
		if err != nil {
			return ""
		}
		re := regexp.MustCompile(`(.*):.*`)
		ticket := re.ReplaceAllString(string(out), "$1")
		return strings.TrimSpace(ticket)
	}

	re := regexp.MustCompile(fmt.Sprintf(`(?i)((%s-)\d{4,})|(.*)`, board))
	ticket := re.ReplaceAllString(string(out), "$1")
	return strings.TrimSpace(ticket)
}

// buildCommitCommand builds the git commit command
func buildCommitCommand(msg string, body string) string {
	args := append([]string{"commit", "-m", msg}, os.Args[1:]...)
	if body != "" {
		args = append(args, "-m", body)
	}
	cmd := strings.Join(args, " ")
	return cmd
}

// commit commits the changes to git
func commit(command string) error {
	cmd := exec.Command("git", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
