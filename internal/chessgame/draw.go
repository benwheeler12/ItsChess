package chessgame

import (
	_ "embed"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

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

func (g *ChessGame) Draw(screen *ebiten.Image) {

	chessBoardImage := g.chessBoardGraphic.drawChessBoard()

	op := &ebiten.DrawImageOptions{}

	_, _ = chessBoardImage, op

	screen.DrawImage(chessBoardImage, op)
}
