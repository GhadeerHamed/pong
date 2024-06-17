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
const initialBalVelocityRow = 1
const initialBalVelocityCol = 2

type GameObject struct {
	row, col, width, height int
	velRow, velCol          int
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

	for !isGameOver() {
		handleUserInput(readInput(inputChan))

		UpdateState()
		DrawState()
		time.Sleep(75 * time.Millisecond)
	}

	screenWidth, screenHeight := screen.Size()
	winner := GetWinner()

	PrintStringCenter(screenHeight/2-1, screenWidth/2-1, "Game Over!")
	PrintStringCenter(screenHeight/2, screenWidth/2, winner+" wins...")

	screen.Show()
	time.Sleep(3 * time.Second)
	screen.Fini()
}

func DrawState() {
	screen.Clear()

	for _, obj := range gameObjects {
		Print(obj.row, obj.col, obj.width, obj.height, obj.symbol)
	}

	screen.Show()
}

func UpdateState() {
	for i := range gameObjects {
		gameObjects[i].row += gameObjects[i].velRow
		gameObjects[i].col += gameObjects[i].velCol
	}

	if CollidesWithWall(ball) {
		ball.velRow = -ball.velRow
	}
	if CollidesWithPlayerPaddle(ball, player1Paddle) || CollidesWithPlayerPaddle(ball, player2Paddle) {
		ball.velCol = -ball.velCol
	}
}

func CollidesWithWall(obj *GameObject) bool {
	_, screenHeight := screen.Size()
	return obj.row+obj.velRow < 0 || obj.row+obj.velRow >= screenHeight
}

func CollidesWithPlayerPaddle(obj *GameObject, p *GameObject) bool {
	hitRow := obj.row >= p.row && obj.row <= p.row+paddleHeight
	if ball.col < p.col {
		return hitRow && obj.col+obj.velCol >= p.col
	} else {
		return hitRow && obj.col+obj.velCol <= p.col
	}
}

func isGameOver() bool {
	return GetWinner() != ""
}

func GetWinner() string {
	width, _ := screen.Size()
	if ball.col < 0 {
		return "Player 2"
	} else if ball.col > width {
		return "Player 1"
	}

	return ""
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
		velRow: 0,
		velCol: 0,
	}

	player2Paddle = &GameObject{
		row:    paddleStart,
		col:    width - 1,
		width:  1,
		height: paddleHeight,
		symbol: PaddleSymbol,
		velRow: 0,
		velCol: 0,
	}

	ball = &GameObject{
		row:    height / 2,
		col:    width / 2,
		width:  1,
		height: 1,
		symbol: BallSymbol,
		velRow: initialBalVelocityRow,
		velCol: initialBalVelocityCol,
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

func PrintStringCenter(row, col int, str string) {
	col = col - len(str)/2
	PrintString(row, col, str)
}

func Print(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}
