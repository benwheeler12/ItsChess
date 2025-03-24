package chessgame

import (
	"github.com/hajimehoshi/ebiten/v2"
)

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
	if g.promotionLifeCycle.promotionInProgress {
		g.handleMouseClickInPromotionPhase()
	} else {
		g.handleMouseClickInChessPhase()
	}
}

func (g *ChessGame) handleMouseRelease() {
	if g.promotionLifeCycle.promotionInProgress {
		g.handleMouseReleaseInPromotionPhase()
	} else {
		g.handleMouseReleaseInChessPhase()
	}
}

func (g *ChessGame) handleMouseClickInPromotionPhase() {

}

func (g *ChessGame) handleMouseReleaseInPromotionPhase() {

}

func (g *ChessGame) handleMouseClickInChessPhase() {
	mouseSquare, outOfBounds := g.graphicsState.getSquareOfMousePosition(ebiten.CursorPosition())

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

func (g *ChessGame) handleMouseReleaseInChessPhase() {
	// Reset mouse state upon any release
	defer func() {
		g.mouseLifeCycle.resetMouseState()
	}()

	// If the mouse was clicked on an invalid square
	if g.mouseClickedOnInvalidSquare {
		return
	}

	mouseSquare, outOfBounds := g.graphicsState.getSquareOfMousePosition(ebiten.CursorPosition())

	// If mouse if out of bounds,
	if outOfBounds {
		return
	}

	// Move is invalid
	if !contains(g.mouseLifeCycle.possibleMoveSquares, mouseSquare) {
		return
	}

	// Move the piece
	g.chessBoard.movePiece(g.mouseLifeCycle.selectedSquare, mouseSquare)
	if g.chessBoard.promotionTriggeredOnSquare(mouseSquare) {
		g.promotionLifeCycle.promotionInProgress = true
		g.promotionLifeCycle.promotionSquare = mouseSquare
	} else {
		g.whitesTurn = !g.whitesTurn
	}
}
