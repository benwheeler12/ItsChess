package itschess

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	ChessGame
}

// Graphics Consts
const (
	startingWindowWidth  int = 700
	startingWindowHeight int = 700
)

func StartGame() {
	var game *Game = &Game{}
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(startingWindowWidth, startingWindowHeight)
	ebiten.SetWindowTitle("It's Chess")
	game.chessBoard.init()
	game.chessBoardGraphic.init(math.Pi, 1, startingWindowWidth, startingWindowHeight)
	game.mouseLifeCycle.resetMouseState()
	game.promotionLifeCycle.resetPromotionLifeCycle()
	game.whitesTurn = true

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

func (g *Game) Update() error {
	g.ChessGame.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ChessGame.Draw(screen)
}
