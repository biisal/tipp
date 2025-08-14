package game

import (
	"log"

	"github.com/biisal/tipp/internal/utils"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TippModel struct {
	Input          textinput.Model
	TextView       viewport.Model
	TypedText      string
	FullText       string
	InputCharCount int
	Content        string
	Width          int
	Height         int
	Words          int
}

func InitTippModel() TippModel {
	input := textinput.New()
	textView := viewport.New(20, 10)
	input.Focus()
	words, err := utils.GetWordFromFile(10)
	if err != nil {
		log.Fatal("could not get words", err)
	}
	return TippModel{
		Content:  words,
		Input:    input,
		TextView: textView,
		Words:    1,
	}
}

func (t TippModel) Init() tea.Cmd {
	return textinput.Blink
}

func (t TippModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Update the input first (handles backspace, ctrl+w, etc.)
	t.Input, cmd = t.Input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return t, tea.Quit
		case " ":
			if t.Input.Value() != "" {
				t.TypedText += t.Input.Value()
				t.Input.Reset()
			}
		}
	case tea.WindowSizeMsg:
		t.Width = msg.Width
		t.Height = msg.Height
	}

	t.FullText = t.TypedText + t.Input.Value()

	return t, cmd
}

func (t TippModel) View() string {
	words := utils.TextViewStyle.Width(t.Width - 5).Render(utils.TextViewWithStats(t.FullText, t.Content))

	input := utils.InputStyle.
		Width(int(float32(t.Width) * 0.5)).
		Render(t.Input.View())

	topPart := lipgloss.Place(
		t.Width,
		t.Height-5,
		lipgloss.Center,
		lipgloss.Center,
		words,
	)

	bottomPart := lipgloss.Place(
		t.Width,
		5,
		lipgloss.Center,
		lipgloss.Bottom,
		input,
	)

	return lipgloss.JoinVertical(lipgloss.Top, topPart, bottomPart)
}
