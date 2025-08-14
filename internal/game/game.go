package game

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/biisal/tipp/internal/utils"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TippModel struct {
	Input      textinput.Model
	TextView   viewport.Model
	TypedText  string
	FullText   string
	Content    string
	Width      int
	Height     int
	Words      int
	ShowResult bool
	StartTime  time.Time
	EndTime    time.Time
}

func InitTippModel() *TippModel {
	input := textinput.New()
	textView := viewport.New(20, 10)
	input.Focus()
	words, err := utils.GetWordFromFile(10)
	if err != nil {
		log.Fatal("could not get words", err)
	}
	return &TippModel{
		Content:    words,
		Input:      input,
		TextView:   textView,
		ShowResult: false,
		TypedText:  "",
		FullText:   "",
		Words:      1,
		StartTime:  time.Now(),
	}
}

func (t TippModel) Init() tea.Cmd {
	return textinput.Blink
}

func (t *TippModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	t.Input, cmd = t.Input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return t, tea.Quit
		case "esc":
			if t.ShowResult {
				return t, tea.Quit
			}
			t.EndTime = time.Now()
			t.ShowResult = true
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

	if len(t.FullText) >= len(t.Content) && !t.ShowResult {
		t.EndTime = time.Now()
		t.ShowResult = true
	}
	return t, cmd
}

func (t TippModel) View() string {
	words, totalCount, correctCount := utils.TextViewWithStats(t.FullText, t.Content)
	var s string
	if !t.ShowResult {
		wordsView := utils.TextViewStyle.Width(t.Width - 5).Render(words)
		input := utils.InputStyle.
			Width(int(float32(t.Width) * 0.5)).
			Render(t.Input.View())

		topPart := lipgloss.Place(
			t.Width,
			t.Height-5,
			lipgloss.Center,
			lipgloss.Center,
			wordsView,
		)

		bottomPart := lipgloss.Place(
			t.Width,
			5,
			lipgloss.Center,
			lipgloss.Bottom,
			input,
		)
		s = lipgloss.JoinVertical(lipgloss.Top, topPart, bottomPart)
	} else {
		timeTaken := t.EndTime.Sub(t.StartTime)
		wpm := float64(len(strings.Split(t.FullText, " "))) / timeTaken.Minutes()
		accuracy := float64(correctCount) * 100 / float64(totalCount)
		columns := []table.Column{
			{Title: "Result", Width: 20},
			{Title: "", Width: 10},
		}
		rows := []table.Row{
			{"Accuracy", strconv.FormatFloat(accuracy, 'f', 2, 64) + "%"},
			{"WPM", strconv.Itoa(int(wpm))},
		}
		result := table.New(table.WithColumns(columns), table.WithRows(rows))
		result.SetHeight(3)
		resultVivew := lipgloss.NewStyle().BorderForeground(lipgloss.Color("#ffffff")).BorderStyle(lipgloss.RoundedBorder())
		s = lipgloss.Place(
			t.Width,
			t.Height,
			lipgloss.Center,
			lipgloss.Center,
			resultVivew.Render(result.View()),
		)
	}
	return s
}
