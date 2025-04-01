package itschess

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	_ "embed"
)

type ChessGame struct {
	chessBoard
	chessBoardGraphic
	mouseState
	mouseLifeCycle
	promotionLifeCycle
	whitesTurn bool
}

func (g *ChessGame) Update() error {

	if g.chessBoard.playerInCheckMate(white) || g.chessBoard.playerInCheckMate(black) {
		return nil // Game over baby
	}

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

func (g *ChessGame) Draw(screen *ebiten.Image) {

	chessBoardImage := g.chessBoardGraphic.drawChessBoard(&g.chessBoard)

	op := &ebiten.DrawImageOptions{}

	_, _ = chessBoardImage, op

	screen.DrawImage(chessBoardImage, op)

}

func (g *ChessGame) handleMouseClick() {
	x, y := ebiten.CursorPosition()
	mousePosition := vector2{x, y}

	// Handle Out of Bounds case
	if !g.chessBoardGraphic.positionInGraphic(mousePosition) {
		return
	}

	clickedElement := g.chessBoardGraphic.getElementAtMousePosition(mousePosition)

	if clickedElement.isChessSquare {
		mouseSquare := clickedElement.square
		// Handle click on empty square
		if g.chessBoard.isEmpty(mouseSquare) {
			return
		}

		// Handle Click on Actual Piece
		clickedPiece := g.chessBoard.getPiece(mouseSquare)

		// If it is not the turn of the piece color that was clicked, do nothing
		if (g.whitesTurn && clickedPiece.color == black) ||
			(!g.whitesTurn && clickedPiece.color == white) {
			return
		}
		// Else, load the mouseLifeCycle State with info of clicked piece
		g.chessBoardGraphic.clickedSquare = mouseSquare
		g.chessBoardGraphic.possibleMoveSquares = g.chessBoard.getValidMoves(g.chessBoardGraphic.clickedSquare)
	} else if clickedElement.isPromotionSquare {
		clickedSquare := clickedElement.square
		if clickedSquare == nilSquare {
			return
		}
		g.chessBoardGraphic.clickedSquare = clickedElement.square
	}

}

func (g *ChessGame) handleMouseRelease() {
	// Reset Chessboard Graphic State
	defer func() {
		g.chessBoardGraphic.possibleMoveSquares = nil
		g.chessBoardGraphic.clickedSquare = nilSquare
	}()
	x, y := ebiten.CursorPosition()
	mousePosition := vector2{x, y}

	// Handle Out of Bounds case
	if !g.chessBoardGraphic.positionInGraphic(mousePosition) {
		return
	}

	clickedElement := g.chessBoardGraphic.getElementAtMousePosition(mousePosition)

	if clickedElement.isChessSquare {
		mouseSquare := clickedElement.square

		// Move is invalid
		if !contains(g.chessBoardGraphic.possibleMoveSquares, mouseSquare) {
			return
		}

		// Move the piece
		g.chessBoard.movePiece(g.chessBoardGraphic.clickedSquare, mouseSquare)
		if g.chessBoard.promotionTriggeredOnSquare(mouseSquare) {
			g.chessBoardGraphic.promotionSquare = mouseSquare
		} else {
			g.whitesTurn = !g.whitesTurn
		}
	} else if clickedElement.isPromotionSquare {
		clickedSquare := clickedElement.square
		// Perform Promotion
		if clickedSquare == g.chessBoardGraphic.clickedSquare {
			promotedPieceType := g.chessBoardGraphic.getPromotionPiece(clickedSquare)
			promotedPieceColor := white
			if !g.whitesTurn {
				promotedPieceColor = black
			}
			g.chessBoard.setPiece(g.chessBoardGraphic.promotionSquare, chessPiece{promotedPieceType, promotedPieceColor})
			g.whitesTurn = !g.whitesTurn
			g.chessBoardGraphic.promotionSquare = nilSquare
		}

	}

}

// Base colors
var (
	lightSquareColor = color.RGBA{238, 238, 238, 255} // Off white
	darkSquareColor  = color.RGBA{118, 150, 86, 255}  // Green

	// SVG Images
	//go:embed assets/white_pawn.svg
	whitePawnBytes []byte

	//go:embed assets/black_pawn.svg
	blackPawnBytes []byte

	//go:embed assets/white_bishop.svg
	whiteBishopBytes []byte

	//go:embed assets/black_bishop.svg
	blackBishopBytes []byte

	//go:embed assets/white_knight.svg
	whiteKnightBytes []byte

	//go:embed assets/black_knight.svg
	blackKnightBytes []byte

	//go:embed assets/white_rook.svg
	whiteRookBytes []byte

	//go:embed assets/black_rook.svg
	blackRookBytes []byte

	//go:embed assets/white_queen.svg
	whiteQueenBytes []byte

	//go:embed assets/black_queen.svg
	blackQueenBytes []byte

	//go:embed assets/white_king.svg
	whiteKingBytes []byte

	//go:embed assets/black_king.svg
	blackKingBytes []byte
)
