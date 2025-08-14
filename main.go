package main

import (
	"flag"
	"github.com/biisal/tipp/internal/game"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	words := flag.Int("words", 30, "Number of words to type")
	flag.Parse()
	p := tea.NewProgram(game.InitTippModel(*words), tea.WithAltScreen())
	tea.ClearScreen()
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
