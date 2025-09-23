package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/alessio/shellescape"
)

// getGitPath returns the path to the git executable if git is in PATH and returns an error if not
func getGitPath() (string, error) {
	path, err := exec.LookPath("git")
	if err != nil {
		return "", fmt.Errorf("git not found in PATH: %w", err)
	}
	return path, nil
}

// findGitDir returns the root of the git repository
func findGitDir(gitPath string) (string, error) {
	cmd := exec.Command(gitPath, "rev-parse", "--show-toplevel")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

func checkBoardMatchesBranch(board string, msg string) bool {
	match, _ := regexp.MatchString(fmt.Sprintf(`(?i)%s-\d{1,}`, board), msg)
	return match
}

// getGitTicketNumber returns the most recent ticket number from the current git branch
func getGitTicketNumber(board string) string {
	cmd := exec.Command("git", "branch", "--show-current")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	match := checkBoardMatchesBranch(board, string(out))
	if !match {
		cmd = exec.Command("git", "log", "-1", "--grep", board, "--oneline", "--format=%s")
		out, err = cmd.Output()
		if err != nil {
			return ""
		}
	}

	ticket := getTicketNumberFromString(string(out), board)
	return strings.TrimSpace(ticket)
}

func getTicketNumberFromString(str string, sub string) string {
	expr := fmt.Sprintf(`(?i).*((%s-)\d{1,})|(.*)`, sub)

	re := regexp.MustCompile(expr)
	return re.ReplaceAllString(str, "$1")
}

// buildCommitCommand builds the git commit command
func buildCommitCommand(msg string, body string, osArgs []string) ([]string, string) {
	args := append([]string{"commit", " -m", msg}, osArgs...)
	if body != "" {
		args = append(args, "-m", body)
	}
	return args, fmt.Sprintf("git %v", shellescape.QuoteCommand(args))
}

// getComitters returns a list of comitters from the git log
func getComitters(osArgs []string) ([]string, error) {
	args := append([]string{"log", "--pretty=format:\"%an <%ae>\""}, osArgs...)
	getComitters, err := getOutput(args)
	if err != nil {
		return nil, err
	}
	result := []string{}
	seen := map[string]bool{}
	for _, comitter := range getComitters {
		if _, ok := seen[comitter]; !ok {
			seen[comitter] = true
			result = append(result, comitter)
		}
	}
	return result, nil
}

// commit commits the changes to git
func commit(command []string) error {
	cmd := exec.Command("git", command...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// getOutput returns the output of a git command
func getOutput(command []string) ([]string, error) {
	cmd := exec.Command("git", command...)

	var b []byte
	buf := bytes.NewBuffer(b)
	cmd.Stdin = os.Stdin
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return strings.Split(string(buf.Bytes()), "\n"), nil
}
