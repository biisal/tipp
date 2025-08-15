package game

import (
	"fmt"
	"log"
	"strconv"
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

func InitTippModel(wordsLen int) *TippModel {
	input := textinput.New()
	textView := viewport.New(20, 10)
	input.Focus()
	words, err := utils.GetWordFromFile(wordsLen)
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
	if !t.ShowResult {
		t.Input, cmd = t.Input.Update(msg)
	}
	prevText := t.FullText
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
			lastIdx := len(t.FullText) - 1
			if t.Input.Value() != "" && lastIdx > 0 && len(t.Content) > lastIdx+1 && string(t.Content[lastIdx+1]) == " " {
				t.TypedText += t.Input.Value()
				t.Input.Reset()
			}
		case "q":
			if t.ShowResult {
				return t, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		t.Width = msg.Width
		t.Height = msg.Height
	}
	t.FullText = t.TypedText + t.Input.Value()

	if prevText == "" && !t.ShowResult && t.FullText != "" {
		t.StartTime = time.Now()
	}
	if len(t.FullText) >= len(t.Content) && !t.ShowResult {
		t.EndTime = time.Now()
		t.ShowResult = true
	}

	return t, cmd
}

func (t TippModel) View() string {
	words, _, correctCount := utils.TextViewWithStats(t.FullText, t.Content)
	var s string
	logo := lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#75FFCF")).Margin(1, 2).Padding(0, 1).Bold(true).Render("TIPP")
	topPart := lipgloss.Place(
		t.Width,
		t.Height*20/100,
		lipgloss.Left,
		lipgloss.Top,
		logo,
	)
	if !t.ShowResult {
		wordsView := utils.TextViewStyle.Width(t.Width - 5).Render(words)
		input := utils.InputStyle.
			Width(int(float32(t.Width) * 0.5)).
			Render(t.Input.View())

		middlePart := lipgloss.Place(
			t.Width,
			t.Height*50/100,
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
		// bottomPart = lipgloss.JoinVertical(lipgloss.Center, middlePart, bottomPart)
		s += lipgloss.JoinVertical(lipgloss.Top, topPart, middlePart, bottomPart)
		// s += logo + "\n" + middlePart + "\n" + bottomPart
	} else {
		timeTaken := t.EndTime.Sub(t.StartTime)
		fullTextLen := len(t.FullText)
		wpm := 0
		if t.FullText != "" {
			minutes := timeTaken.Minutes()
			if minutes > 0 {
				wpm = int(float64(len(t.FullText)) / 5.0 / minutes)
			}
		}

		accuracy := 0.0
		if fullTextLen > 0 {
			accuracy = float64(correctCount) / float64(fullTextLen) * 100
		}
		columns := []table.Column{
			{Title: "Result", Width: 20},
			{Title: "", Width: 10},
		}
		rows := []table.Row{
			{"Accuracy", strconv.FormatFloat(accuracy, 'f', 2, 64) + "%"},
			{"WPM", strconv.Itoa(wpm)},
			{"Correct KeyPresses", fmt.Sprintf("%d/%d", correctCount, fullTextLen)},
		}
		result := table.New(table.WithColumns(columns), table.WithRows(rows))
		result.SetHeight(1 + len(rows))
		resultVivew := lipgloss.NewStyle().BorderForeground(lipgloss.Color("#ffffff")).BorderStyle(lipgloss.RoundedBorder())
		bottomPart := lipgloss.Place(
			t.Width,
			t.Height*70/100,
			lipgloss.Center,
			lipgloss.Center,
			resultVivew.Render(result.View()),
		)
		s += lipgloss.JoinVertical(lipgloss.Top, topPart, bottomPart)
	}
	return s
}
