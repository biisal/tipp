package main

import (
	"flag"
	"fmt"

	"github.com/biisal/tipp/internal/game"
	"github.com/biisal/tipp/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	words := flag.Int("w", 30, "Number of words to type. Max is 300")
	availableModes := "\n\t"
	for m := range utils.WordsMap {
		availableModes += m + "\n\t"
	}
	mode := flag.String("m", "eng", "Mode to use. Options are"+availableModes)
	custom := flag.String("c", "", "Custom words file path")
	flag.Parse()
	if *words <= 0 || *words > 300 {
		fmt.Println("Number of words must be between 1 and 300")
		return
	}
	p := tea.NewProgram(game.InitTippModel(*words, *mode, *custom), tea.WithAltScreen())
	tea.ClearScreen()
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
