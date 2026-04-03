package config

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
)

type CoAuthor struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Selected bool
}

type CoAuthors []CoAuthor

func (p *CoAuthors) Options() []huh.Option[string] {
	coAuthors := []CoAuthor(*p)

	if len(coAuthors) == 0 {
		return nil
	}
	items := []huh.Option[string]{}
	for _, coauthor := range coAuthors {
		desc := fmt.Sprintf("%s <%s>", coauthor.Name, coauthor.Email)
		items = append(items, huh.NewOption(desc, desc))
	}
	return items
}

// BuildCoauthorString takes a slice of selected coauthors and returns a formatted
// string which Github recognises
func BuildCoAuthorString(coauthors []string) string {
	var s strings.Builder
	s.WriteString(`


	`)

	for _, coauthor := range coauthors {
		if coauthor == "none" {
			return ""
		}
		fmt.Fprintf(&s, "\nCo-authored-by: %s", coauthor)
	}
	return s.String()
}
