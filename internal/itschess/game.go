package itschess

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type ChessGame struct {
	chessBoard
	chessBoardGraphic
	mouseState
	mouseLifeCycle
	promotionLifeCycle
	whitesTurn bool
}

// Graphics Consts
const (
	startingWindowWidth  int = 700
	startingWindowHeight int = 700
)

func StartGame() {
	var game *ChessGame = &ChessGame{}
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
