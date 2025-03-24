package chessgame

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// Graphics Consts
const (
	startingWindowWidth  int = 700
	startingWindowHeight int = 700

	pieceScale float64 = .9 // As a factor of Chess Square Dimensions

	// Rotation factor of the baord between matrix state and drawn state.  Rotation in degrees is rotationFactor * 90
	rotationFactor int = 2
	flipFactor     int = 1
)

// Base colors
var (
	lightSquareColor = color.RGBA{238, 238, 238, 255} // Off white
	darkSquareColor  = color.RGBA{118, 150, 86, 255}  // Green

	// Yellow tinted variants
	lightSquareHighlightColor = color.RGBA{255, 255, 200, 50}  // Light yellow
	darkSquareHighlightColor  = color.RGBA{255, 255, 200, 255} // Green with yellow tint

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

// Methods for basic dimensions as functions of window height and width
func (g *ChessGame) boardWidth() int {
	return min(g.graphicsState.windowWidth, g.graphicsState.windowHeight) + 4
}

func (g *ChessGame) boardHeight() int {
	return min(g.graphicsState.windowWidth, g.graphicsState.windowHeight) + 4
}

func (g *ChessGame) squareWidth() float64 {
	return float64(g.boardWidth() / 8)
}

func (g *ChessGame) squareHeight() float64 {
	return float64(g.boardHeight() / 8)
}

func (g *ChessGame) pieceWidth() float64 {
	return g.squareWidth() * pieceScale
}

func (g *ChessGame) pieceHeight() float64 {
	return g.squareHeight() * pieceScale
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

func tintColor(originalColor color.RGBA, tint color.RGBA, alpha float64) color.RGBA {
	// Ensure alpha is in range [0,1]
	if alpha < 0 || alpha > 1 {
		panic("alpha value not between 0 and 1")
	}

	return color.RGBA{
		uint8((float64(originalColor.R) * (1 - alpha)) + (float64(tint.R) * alpha)),
		uint8((float64(originalColor.G) * (1 - alpha)) + (float64(tint.G) * alpha)),
		uint8((float64(originalColor.B) * (1 - alpha)) + (float64(tint.B) * alpha)),
		originalColor.A,
	}

}

// square refers to the canonical chess square that this function call will draw
func (g *ChessGame) drawChessSquare(screen *ebiten.Image, x float64, y float64, width float32, height float32, square vector2) {

	// Calculate Color
	squareColor := lightSquareColor
	if g.chessBoard.isDarkSquare(square) {
		squareColor = darkSquareColor
	}

	if contains(g.mouseLifeCycle.possibleMoveSquares, square) {
		squareColor = tintColor(squareColor, color.RGBA{255, 255, 0, 1}, 0.3)
	}

	vector.DrawFilledRect(screen,
		float32(x),
		float32(y),
		width,
		height,
		squareColor,
		true)

	// drawChessPiece draws from the center of the piece.  Need to calculate the midpoint of the chess square
	if g.mouseLifeCycle.selectedSquare != square {
		squareCenterX := float64(x) + (g.squareWidth() / 2)
		squareCenterY := float64(y) + (g.squareHeight() / 2)
		g.drawChessPiece(g.chessBoard.getPiece(square), squareCenterX, squareCenterY, screen)
	}

}

func (g *ChessGame) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{255, 255, 255, 255})

	if g.chessBoard.playerInCheckMate(white) || g.chessBoard.playerInCheckMate(black) {
		return
	}

	// Draw Board
	for file := range 8 {
		for rank := range 8 {

			squareX := g.squareWidth() * float64(file)
			squareY := g.squareHeight() * float64(rank)

			g.drawChessSquare(screen, squareX, squareY, float32(g.squareWidth()), float32(g.squareHeight()), rotateCoord(vector2{file, rank}, rotationFactor))
			continue
		}
	}

	// Draw the clicked piece
	if g.mouseLifeCycle.selectedPiece != emptyPiece {
		mouseX, mouseY := ebiten.CursorPosition()
		g.drawChessPiece(g.selectedPiece, float64(mouseX), float64(mouseY), screen)
	}
}

func (g *ChessGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Window Size has changed!  Reinitialize piece images
	if outsideWidth != g.windowWidth || outsideHeight != g.windowHeight {
		// Order matter here
		g.windowWidth = outsideWidth
		g.windowHeight = outsideHeight
		g.initializePieceImages()
	}

	return outsideWidth, outsideHeight
}

func (g *ChessGame) getPieceImage(pieceBytes []byte) *ebiten.Image {
	// Open the SVG file
	/**file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()*/

	reader := bytes.NewReader(pieceBytes)

	// Decode SVG
	icon, err := oksvg.ReadIconStream(reader)
	if err != nil {
		log.Fatal(err)
	}

	// Set SVG viewport
	icon.SetTarget(0, 0, g.pieceWidth(), g.pieceHeight())

	w := g.pieceWidth()
	_ = w
	// Create an RGBA image
	img := image.NewRGBA(image.Rect(0, 0, int(math.Round(g.pieceWidth())), int(math.Round(g.pieceHeight()))))
	drawer := rasterx.NewDasher(
		int(math.Round(g.pieceWidth())),
		int(math.Round(g.pieceHeight())),
		rasterx.NewScannerGV(
			int(math.Round(g.pieceWidth())),
			int(math.Round(g.pieceHeight())),
			img,
			img.Bounds(),
		),
	)
	icon.Draw(drawer, 1.0)
	return ebiten.NewImageFromImage(img)
}

// x and y represnt the mid point of the image.
func (g *ChessGame) drawChessPiece(piece chessPiece, x float64, y float64, screen *ebiten.Image) {

	image, ok := g.pieceEbitenMap[piece]
	if !ok {
		return
	}

	x = math.Round(x - (g.pieceWidth() / 2))
	y = math.Round(y - (g.pieceHeight() / 2))

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	op.Filter = ebiten.FilterNearest

	screen.DrawImage(image, op)
}

func rotateCoord(square vector2, rotationFactor int) vector2 {
	// normalize between 0 and 3
	rotationFactor = ((rotationFactor % 4) + 4) % 4

	x, y := 7-square.x, square.y

	switch rotationFactor {
	case 0:
		return vector2{x, y}
	case 1:
		return vector2{7 - y, x}
	case 2:
		return vector2{7 - x, 7 - y}
	case 3:
		return vector2{y, 7 - x}
	}
	return vector2{x, y} // default case, shouldn't be reached
}

func (g *ChessGame) getSquareOfMousePosition(x int, y int) (vector2, bool) {
	file := int(float64(x) / g.squareWidth())
	rank := int(float64(y) / g.squareHeight())

	// Check if coordinates are within bounds
	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return nilSquare, true
	}

	// Apply rotation based on rotationFactor
	return rotateCoord(vector2{file, rank}, -rotationFactor), false
}
