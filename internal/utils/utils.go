package utils

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/biisal/tipp/words"
	"github.com/charmbracelet/lipgloss"
)

var (
	InputStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#130F1A")). // Gray background
			Foreground(lipgloss.Color("#ffffff")). // White text
			Padding(1, 1).Margin(2, 0)
	TextViewStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff")).AlignHorizontal(lipgloss.Center)
	TypedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#0FF563"))
	MistakeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	RemainingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#A8A8A8"))

	wordsFiledir = "/.config/tipp/"
)

const (
	wordsFileName = "words.txt"
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	wordsFiledir = home + wordsFiledir
}

func GetWordFromFile(n int) (string, error) {
	permission := 0644
	fullFilePath := wordsFiledir + wordsFileName
	if err := os.MkdirAll(wordsFiledir, os.ModePerm); err != nil {
		return "", err
	}
	file, err := os.OpenFile(fullFilePath, os.O_RDONLY, os.FileMode(permission))
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("The words.txt file does not exits in the current directory.\nDo you want to create it? (y/n)")
			var input string
			fmt.Scanln(&input)
			if input == "y" || input == "Y" {
				file, err = os.Create(fullFilePath)
				if err != nil {
					return "", err
				}
				file.WriteString(words.DEFAULT_WORDS)
				file.Close()
				return GetWordFromFile(n)
			} else {
				return "", err
			}
		}
	}
	defer file.Close()
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
	s, mistakeCount := "", 0

	for i, w := range words {
		if i < typedTextLen {
			if byte(w) == typedText[i] {
				s += TypedStyle.Render(string(w))
			} else {
				s += MistakeStyle.Render(string(w))
				mistakeCount++
			}
		} else if i == typedTextLen {
			s += RemainingStyle.Background(lipgloss.Color("#CB97FF")).Foreground(lipgloss.Color("#000000")).Render(string(w))
		} else {
			s += RemainingStyle.Render(string(w))
		}
	}
	return s, wordsLen, typedTextLen - mistakeCount
}
