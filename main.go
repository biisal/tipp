package main

import (
	"flag"
	"fmt"

	"github.com/biisal/tipp/internal/game"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	words := flag.Int("w", 30, "Number of words to type. Max is 600")
	flag.Parse()
	if *words <= 0 || *words > 600 {
		fmt.Println("Number of words must be between 1 and 600")
		return
	}
	p := tea.NewProgram(game.InitTippModel(*words), tea.WithAltScreen())
	tea.ClearScreen()
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
