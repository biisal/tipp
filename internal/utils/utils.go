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

type Word struct {
	Default  string
	FileName string
}

var WordsMap = map[string]Word{
	"eng": {words.ENGLISH_WORDS, "eng_words.txt"},
	"py":  {words.PYTHON_WORDS, "py_words.txt"},
	"go":  {words.GO_WORDS, "go_words.txt"},
	"js":  {words.JS_WORDS, "js_words.txt"},
	"c":   {words.C_WORDS, "c_words.txt"},
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	wordsFiledir = home + wordsFiledir
}

func GetWordFromFile(n int, mode string, custom string) (string, error) {
	permission := 0644
	fullFilePath := custom
	wordsFileName := Word{}
	if custom == "" {
		exits := false
		wordsFileName, exits = WordsMap[mode]
		if !exits {
			availableModes := "\n\t"
			for m := range WordsMap {
				availableModes += m + "\n\t"
			}
			return "", fmt.Errorf("invalid mode: %s,\navailable modes: %v", mode, availableModes)
		}
		fullFilePath = wordsFiledir + wordsFileName.FileName
		permission = 0666
	}
	if err := os.MkdirAll(wordsFiledir, os.ModePerm); err != nil {
		return "", err
	}
	file, err := os.OpenFile(fullFilePath, os.O_RDONLY, os.FileMode(permission))
	if err != nil {
		if os.IsNotExist(err) {
			if custom != "" {
				return "", fmt.Errorf("your custom file %s does not exits", custom)
			} else {
				fmt.Printf("The %s file does not exits in %s directory.\nDo you want to create it? (y/n)\n: ", wordsFileName.FileName, wordsFiledir)
				var input string
				fmt.Scanln(&input)
				if input == "y" || input == "Y" {
					file, err = os.Create(fullFilePath)
					if err != nil {
						return "", err
					}
					file.WriteString(wordsFileName.Default)
					file.Close()
					return GetWordFromFile(n, mode, custom)
				} else {
					return "", fmt.Errorf("Exiting program! Either create the file %s manually in %s directory or re run the command and press y", wordsFileName.FileName, wordsFiledir)
				}
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
			if w == "" {
				continue
			}
			words = append(words, strings.TrimSpace(w))
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

const (
	RED    = "\033[0;31m"
	GREEN  = "\033[0;32m"
	PURPLE = "\033[0;35m"
	RESET  = "\033[0m"
)

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
				s += fmt.Sprintf("%s%c%s", PURPLE, w, RESET)
				mistakeCount++
			}
		} else if i == typedTextLen {
			s += RemainingStyle.Background(lipgloss.Color("#CB97FF")).Foreground(lipgloss.Color("#000000")).Render(string(w))
		} else {
			s += string(w)
		}
	}
	return s, wordsLen, typedTextLen - mistakeCount
}

func GetInstructionsStyle() (string, int) {
	quteInstractioons := lipgloss.NewStyle().Foreground(lipgloss.Color("#8AA0AC")).
		Background(lipgloss.Color("#393939")).
		Margin(1, 2).
		Padding(0, 1).
		Render("press esc to show result")
	return quteInstractioons, lipgloss.Height(quteInstractioons)
}

func GetLogo() string {
	logo := lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#75FFCF")).Margin(1, 2).Padding(0, 1).Bold(true).Render(`𝗧𝗜𝗣𝗣`)
	return logo
}
