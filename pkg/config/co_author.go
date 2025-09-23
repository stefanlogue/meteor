package config

import (
	"fmt"

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
