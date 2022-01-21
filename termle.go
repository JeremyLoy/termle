package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

const totalWordles = 2315

var (
	//go:embed guesses.txt
	//go:embed answers.txt
	files embed.FS

	dayFlag    = flag.Int("day", currentDay(), "select a specific wordle by day")
	randomFlag = flag.Bool("random", false, "pick a random wordle")

	firstDay = time.Date(2021, time.June, 19, 0, 0, 0, 0, time.UTC)
	valid    = regexp.MustCompile(`^[A-Za-z]{5}$`)
)

func main() {
	flag.Parse()

	var day int
	if *randomFlag {
		day = randomDay()
	} else {
		day = *dayFlag
	}

	answer := answerForDay(day)
	validGuesses := guessesSet()
	b := newBoard()
	currentTurn := 0
	turnsRemaining := 6

	s := bufio.NewScanner(os.Stdin)

	printTurn(b)
	for s.Scan() {
		guess := strings.ToUpper(strings.TrimSpace(s.Text()))
		if !valid.MatchString(guess) {
			printTurnWithError(b, "Please enter a 5 letter word")
			continue
		}
		if _, ok := validGuesses[guess]; !ok {
			printTurnWithError(b, "Not in word list")
			continue
		}
		b.addGuess(currentTurn, guess, answer)
		printTurn(b)
		if guess == answer {
			fmt.Print("you won!")
			return
		}
		turnsRemaining--
		currentTurn++
		if turnsRemaining == 0 {
			fmt.Println("you lose!")
			fmt.Println("Answer was ", answer)
			return
		}
	}
	if s.Err() != nil {
		panic(s.Err())
	}
}

func green(l string) string {
	return "\033[37;102m" + l + "\033[0m"
}
func yellow(l string) string {
	return "\033[37;103m" + l + "\033[0m"
}
func white(l string) string {
	return "\033[0;107m" + l + "\033[0m"
}
func black(l string) string {
	return "\033[37;100m" + l + "\033[0m"
}
func clearBoard() {
	fmt.Print("\033c")
}
func prompt() {
	fmt.Print(">")
}

type board [][]string

func newBoard() board {
	b := make(board, 6)
	for i := range b {
		b[i] = make([]string, 5)
		for j := range b[i] {
			b[i][j] = black("_")
		}
	}
	return b
}

func (b board) addGuess(turn int, guess, answer string) {
	for i, c := range guess {
		if c == rune(answer[i]) {
			b[turn][i] = green(string(c))
		} else if strings.ContainsRune(answer, c) {
			b[turn][i] = yellow(string(c))
		} else {
			b[turn][i] = black(string(c))
		}
	}
}

func (b board) print() {
	for _, y := range b {
		for _, x := range y {
			fmt.Print(" " + x)
		}
		fmt.Println()
	}
}

func printTurn(b board) {
	clearBoard()
	b.print()
	prompt()
}

func printTurnWithError(b board, err string) {
	clearBoard()
	b.print()
	fmt.Println(err)
	prompt()
}

func guessesSet() map[string]struct{} {
	guessesFile, err := files.Open("guesses.txt")
	if err != nil {
		panic(err)
	}
	validGuesses := make(map[string]struct{})
	guessesReader := bufio.NewScanner(guessesFile)

	for guessesReader.Scan() {
		validGuesses[strings.ToUpper(guessesReader.Text())] = struct{}{}
	}
	if guessesReader.Err() != nil {
		panic(err)
	}
	return validGuesses
}

func answerForDay(day int) string {
	answersFile, err := files.Open("answers.txt")
	if err != nil {
		panic(err)
	}
	seeker := answersFile.(io.ReadSeeker)
	_, err = seeker.Seek(int64(day*7), io.SeekStart)
	if err != nil {
		panic(err)
	}
	answer := make([]byte, 5)
	_, err = seeker.Read(answer)
	if err != nil {
		panic(err)
	}
	return strings.ToUpper(string(answer))
}

func currentDay() int {
	return int(time.Now().UTC().Sub(firstDay).Hours()/24) - 1
}

func randomDay() int {
	return rand.Intn(totalWordles)
}
