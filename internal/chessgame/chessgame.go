package chessgame

import (
	"github.com/hajimehoshi/ebiten/v2"
)

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
	game.whitesTurn = true

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

func (g *ChessGame) Update() error {

	g.lastMouseState = g.mousePressed

	g.mousePressed = ebiten.IsMouseButtonPressed((ebiten.MouseButtonLeft))

	if g.mousePressed != g.lastMouseState {
		if g.mousePressed {
			g.handleMouseClick()
		} else {
			g.handleMouseRelease()
		}
	}

	return nil
}

func (g *ChessGame) handleMouseClick() {
	mouseSquare, outOfBounds := g.getSquareOfMousePosition(ebiten.CursorPosition())

	// Handle Out of Bounds case
	if outOfBounds {
		g.mouseClickedOnInvalidSquare = true
		return
	} else {
		g.mouseClickedOnInvalidSquare = false
	}

	// Handle click on empty square
	if g.chessBoard.isEmpty(mouseSquare) {
		g.mouseClickedOnInvalidSquare = true
		return
	} else {
		g.mouseClickedOnInvalidSquare = false
	}

	// Handle Click on Actual Piece
	clickedPiece := g.chessBoard.getPiece(mouseSquare)

	// If it is not the turn of the piece color that was clicked, do nothing
	if (g.whitesTurn && clickedPiece.color == black) ||
		(!g.whitesTurn && clickedPiece.color == white) {
		g.mouseClickedOnInvalidSquare = true
		return
	}
	// Else, load the mouseLifeCycle State with info of clicked piece
	g.mouseLifeCycle.selectedPiece = g.chessBoard.getPiece(mouseSquare)
	g.mouseLifeCycle.selectedSquare = mouseSquare
	g.mouseLifeCycle.possibleMoveSquares = g.chessBoard.getValidMoves(
		g.mouseLifeCycle.selectedSquare,
		g.mouseLifeCycle.selectedPiece,
	)

}

func (g *ChessGame) handleMouseRelease() {

	// Reset mouse state upon any release
	defer func() {
		g.mouseLifeCycle.resetMouseState()
	}()

	// If the mouse was clicked on an invalid square, just return
	if g.mouseClickedOnInvalidSquare {
		return
	}

	mouseSquare, outOfBounds := g.getSquareOfMousePosition(ebiten.CursorPosition())

	// If mouse if out of bounds, return piece to its original square
	if outOfBounds {
		g.board[g.mouseLifeCycle.selectedSquare.x][g.mouseLifeCycle.selectedSquare.y] = g.mouseLifeCycle.selectedPiece
		g.mouseLifeCycle.selectedPiece = emptyPiece
		return
	}

	// If move is valid, move the piece to the new square
	if contains(g.mouseLifeCycle.possibleMoveSquares, mouseSquare) {
		g.chessBoard.movePiece(g.mouseLifeCycle.selectedSquare, mouseSquare)
		g.whitesTurn = !g.whitesTurn
	}
}
