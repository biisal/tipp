package main

import (
	"github.com/biisal/tipp/internal/game"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(game.InitTippModel(), tea.WithAltScreen())
	tea.ClearScreen()
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
