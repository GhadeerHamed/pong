package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

const PaddleSymbol = 0x2588
const BallSymbol = 0x25CF
const paddleHeight = 4

type GameObject struct {
	row, col, width, height int
	symbol                  rune
}

var screen tcell.Screen
var player1Paddle *GameObject
var player2Paddle *GameObject
var ball *GameObject

var gameObjects []*GameObject

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

func DrawState() {
	screen.Clear()

	for _, obj := range gameObjects {
		Print(obj.row, obj.col, obj.width, obj.height, obj.symbol)
	}

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

	player1Paddle = &GameObject{
		row:    paddleStart,
		col:    0,
		width:  1,
		height: paddleHeight,
		symbol: PaddleSymbol,
	}

	player2Paddle = &GameObject{
		row:    paddleStart,
		col:    width - 1,
		width:  1,
		height: paddleHeight,
		symbol: PaddleSymbol,
	}

	ball = &GameObject{
		row:    height / 2,
		col:    width / 2,
		width:  1,
		height: 1,
		symbol: BallSymbol,
	}

	gameObjects = []*GameObject{
		player1Paddle, player2Paddle, ball,
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

func handleUserInput(key string) {
	_, height := screen.Size()

	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[w]" && player1Paddle.row > 0 {
		player1Paddle.row--
	} else if key == "Rune[s]" && player1Paddle.row+player1Paddle.height < height {
		player1Paddle.row++
	} else if key == "Up" && player2Paddle.row > 0 {
		player2Paddle.row--
	} else if key == "Down" && player2Paddle.row+player2Paddle.height < height {
		player2Paddle.row++
	}
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

