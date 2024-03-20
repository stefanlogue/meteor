package config

import "github.com/charmbracelet/huh"

type Board struct {
	Name string `json:"name"`
}

type Boards []Board

func (p *Boards) Options() []huh.Option[string] {
	boards := []Board(*p)

	if len(boards) == 0 {
		return nil
	}
	items := []huh.Option[string]{}
	for _, board := range boards {
		items = append(items, huh.NewOption(board.Name, board.Name))
	}
	return items
}
