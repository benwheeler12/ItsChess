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

func (g *ChessGame) initializePieceImages() {
	g.pieceEbitenMap = map[chessPiece]*ebiten.Image{
		{pawn, white}:   g.getPieceImage(whitePawnBytes),
		{pawn, black}:   g.getPieceImage(blackPawnBytes),
		{bishop, white}: g.getPieceImage(whiteBishopBytes),
		{bishop, black}: g.getPieceImage(blackBishopBytes),
		{knight, white}: g.getPieceImage(whiteKnightBytes),
		{knight, black}: g.getPieceImage(blackKnightBytes),
		{rook, white}:   g.getPieceImage(whiteRookBytes),
		{rook, black}:   g.getPieceImage(blackRookBytes),
		{queen, white}:  g.getPieceImage(whiteQueenBytes),
		{queen, black}:  g.getPieceImage(blackQueenBytes),
		{king, white}:   g.getPieceImage(whiteKingBytes),
		{king, black}:   g.getPieceImage(blackKingBytes),
	}
}
