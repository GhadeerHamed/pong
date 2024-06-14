package main

import (
	"fmt"
	"os"

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

func DrawState() {
	screen.Clear()

	Print(player1.row, player1.col, player1.width, player1.height, paddleSymbol)
	Print(player2.row, player2.col, player2.width, player2.height, paddleSymbol)

	screen.Show()
}

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
	initScreen()
	initGameState()

	DrawState()

	for {
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
			DrawState()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				screen.Fini()
				os.Exit(0)
			}
		}
	}
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
