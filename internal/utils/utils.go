package utils

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	InputStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#555555")). // Gray background
			Foreground(lipgloss.Color("#ffffff")). // White text
			Padding(1, 1).Margin(2, 0)
	TextViewStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff")).AlignHorizontal(lipgloss.Center)
	TypedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#0FF563"))
	MistakeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	RemainingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#444"))
)

const (
	wordsFiledir = "./words.txt"
	defaltWords  = "let the const be var"
)

func GetWordFromFile(n int) (string, error) {
	file, err := os.Open(wordsFiledir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("The words.txt file does not exits in the current directory.\nDo you want to create it? (y/n)")
			var input string
			fmt.Scanln(&input)
			if input == "y" || input == "Y" {
				file, err = os.Create(wordsFiledir)
				if err != nil {
					return "", err
				}
				file.WriteString(defaltWords)
				file.Close()
				return GetWordFromFile(n)
			} else {
				return "", err
			}
		}
		return "", err
	}
	var words []string
	var randomWords string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for w := range strings.SplitSeq(line, " ") {
			words = append(words, w)
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	wordsLen := len(words)
	if wordsLen == 0 {
		return "", fmt.Errorf("no words found in the file")
	}
	for range n {
		randomWords += words[rand.Intn(wordsLen)] + " "
	}
	return strings.TrimSpace(randomWords), nil
}

// returns viewString , totalCount , correctCount ,
func TextViewWithStats(typedText, words string) (string, int, int) {
	typedTextLen, wordsLen := len(typedText), len(words)

	if typedTextLen > wordsLen {
		typedText = typedText[:len(words)]
	}
	s, mistakCount := "", 0

	for i, w := range words {
		if i < typedTextLen {
			if byte(w) == typedText[i] {
				s += TypedStyle.Render(string(w))
			} else {
				s += MistakeStyle.Render(string(w))
				mistakCount++
			}
		} else if i == typedTextLen {
			s += RemainingStyle.Background(lipgloss.Color("#ffffff")).Render(string(w))
		} else {
			s += RemainingStyle.Render(string(w))
		}
	}
	return s, wordsLen, typedTextLen - mistakCount
}
