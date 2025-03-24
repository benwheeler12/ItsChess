package chessgame

import "github.com/hajimehoshi/ebiten/v2"

func StartGame() {
	var game *ChessGame = &ChessGame{}
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(startingWindowWidth, startingWindowHeight)
	ebiten.SetWindowTitle("Ben's Cutie Pie Chess Game")

	game.windowWidth = startingWindowWidth
	game.windowHeight = startingWindowHeight

	game.initializeBoard()
	game.initializePieceImages()
	game.mouseLifeCycle.resetMouseState()
	game.promotionLifeCycle.resetPromotionLifeCycle()
	game.whitesTurn = true

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
