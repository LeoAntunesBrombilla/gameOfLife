package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
	"math/rand"
	"time"
)

type Game struct {
	board  *Board
	ticker *time.Ticker
}

type Cell struct {
	x, y, width, height int
	color               color.Color
	alive               bool
}

type Board struct {
	width, height int
	grid          [][]*Cell
}

func NewBoard(width, height int) *Board {
	board := &Board{
		width:  width,
		height: height,
		grid:   make([][]*Cell, height),
	}

	for y := 0; y < height; y++ {
		board.grid[y] = make([]*Cell, width)
		for x := 0; x < width; x++ {
			board.grid[y][x] = &Cell{
				x:      x * 10,
				y:      y * 10,
				width:  10,
				height: 10,
				alive:  false,
			}
		}
	}

	return board
}

func (c *Cell) Draw(target *ebiten.Image) {

	col := color.RGBA{0, 0, 0, 255}
	if c.alive {
		col = color.RGBA{255, 255, 255, 255}
	}

	img := ebiten.NewImage(c.width, c.height)
	img.Fill(col)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.x), float64(c.y))

	target.DrawImage(img, op)
}

func (b *Board) Draw(target *ebiten.Image) {
	for _, row := range b.grid {
		for _, cell := range row {
			cell.Draw(target)
		}
	}
}

func (b *Board) Step() {
	nextGrid := make([][]*Cell, b.height)
	for y := 0; y < b.height; y++ {
		nextGrid[y] = make([]*Cell, b.width)
		for x := 0; x < b.width; x++ {
			neighbors := b.counAliveNeighbors(x, y)
			cell := b.grid[y][x]
			alive := cell.alive
			if alive && (neighbors < 2 || neighbors > 3) {
				alive = false
			} else if !alive && neighbors == 3 {
				alive = true
			}

			nextGrid[y][x] = &Cell{
				x:      x * 10,
				y:      y * 10,
				width:  10,
				height: 10,
				alive:  alive,
			}
		}
	}
	b.grid = nextGrid
}

func (b *Board) counAliveNeighbors(x, y int) int {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}

			neighborX := (x + i + b.width) % b.width
			neighborY := (y + j + b.height) % b.height

			if b.grid[neighborY][neighborX].alive {
				count++
			}
		}
	}

	return count
}

func (g *Game) Update() error {
	select {
	case <-g.ticker.C:
		g.board.Step()
	default:
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.board.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {

	board := NewBoard(64, 48)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000; i++ {
		x := rand.Intn(64)
		y := rand.Intn(48)
		board.grid[y][x].alive = true
	}

	game := &Game{
		board:  board,
		ticker: time.NewTicker(100 * time.Millisecond),
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Conway Game of life")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
