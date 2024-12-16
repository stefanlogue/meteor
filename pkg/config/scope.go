package config

import "github.com/charmbracelet/huh"

type Scope struct {
	Name string `json:"name"`
}

type Scopes []Scope

func (s *Scopes) Options() []huh.Option[string] {
	scopes := []Scope(*s)
	if len(scopes) == 0 {
		return nil
	}
	items := []huh.Option[string]{}
	for _, scope := range scopes {
		items = append(items, huh.NewOption(scope.Name, scope.Name))
	}
	return items
}
