package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

const paddleSymbol = 0x2588
const paddleHeight = 4

type Paddle struct {
	row, col, width, height int
}

var screen tcell.Screen
var player1 *Paddle
var player2 *Paddle

func main() {
	initScreen()
	initGameState()

	inputChan := initUserInput()

	for {
		DrawState()
		time.Sleep(50 * time.Millisecond)

		key := readInput(inputChan)

		handleUserInput(key)

	}
}
func handleUserInput(key string) {
	_, height := screen.Size()

	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[w]" && player1.row > 0 {
		player1.row--
	} else if key == "Rune[s]" && player1.row+player1.height < height {
		player1.row++
	} else if key == "Up" && player2.row > 0 {
		player2.row--
	} else if key == "Down" && player2.row+player2.height < height {
		player2.row++
	}
}

func DrawState() {
	screen.Clear()

	Print(player1.row, player1.col, player1.width, player1.height, paddleSymbol)
	Print(player2.row, player2.col, player2.width, player2.height, paddleSymbol)

	screen.Show()
}

func initUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputChan <- ev.Name()
			}
		}
	}()

	return inputChan
}

func initScreen() {
	var err error

	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err = screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
}

func initGameState() {
	width, height := screen.Size()

	paddleStart := height/2 - paddleHeight/2

	player1 = &Paddle{
		row:    paddleStart,
		col:    0,
		width:  1,
		height: paddleHeight,
	}

	player2 = &Paddle{
		row:    paddleStart,
		col:    width - 1,
		width:  1,
		height: paddleHeight,
	}
}

func readInput(inputChan chan string) string {
	var key string
	select {
	case key = <-inputChan:
	default:
		key = ""
	}

	return key
}

func PrintString(row, col int, str string) {
	for _, c := range str {
		screen.SetContent(col, row, c, nil, tcell.StyleDefault)
		col += 1
	}
}

func Print(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}
