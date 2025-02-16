package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	x float64
}

func (g *Game) Update() error {
	g.x += 2 // Move the circle dynamically
	if g.x > 320 {
		g.x = 0 // Reset when reaching the edge
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw a circle that moves horizontally
	ebitenutil.DrawCircle(screen, g.x, 100, 20, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 240
}

func main() {
	game := &Game{}
	ebiten.SetWindowSize(320, 240)
	ebiten.SetWindowTitle("Dynamic Drawing Example")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
