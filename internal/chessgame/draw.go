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

	chessBoardImage := ebiten.NewImage(g.boardWidth(), g.boardHeight())

	chessBoardImage.Fill(color.RGBA{255, 255, 255, 255})

	if g.chessBoard.playerInCheckMate(white) || g.chessBoard.playerInCheckMate(black) {
		return
	}

	// Draw Board
	for file := range 8 {
		for rank := range 8 {

			squareX := g.squareWidth() * float64(file)
			squareY := g.squareHeight() * float64(rank)

			g.drawChessSquare(chessBoardImage, squareX, squareY, float32(g.squareWidth()), float32(g.squareHeight()), rotateChessCoord(vector2{file, rank}, rotationFactor))
			continue
		}
	}

	// Draw the clicked piece
	if g.mouseLifeCycle.selectedPiece != emptyPiece {
		mouseX, mouseY := ebiten.CursorPosition()
		g.drawChessPiece(g.selectedPiece, float64(mouseX), float64(mouseY), chessBoardImage)
	}
	promotionPopupCenter := g.graphicsState.getPromotionPopupOrigin(vector2{0, 0})
	g.drawPromotionBox(chessBoardImage, promotionPopupCenter)

	promotionPopupCenter = g.graphicsState.getPromotionPopupOrigin(vector2{0, 7})
	g.drawPromotionBox(chessBoardImage, promotionPopupCenter)

	promotionPopupCenter = g.graphicsState.getPromotionPopupOrigin(vector2{7, 7})
	g.drawPromotionBox(chessBoardImage, promotionPopupCenter)

	promotionPopupCenter = g.graphicsState.getPromotionPopupOrigin(vector2{7, 0})
	g.drawPromotionBox(chessBoardImage, promotionPopupCenter)

	if g.promotionLifeCycle.promotionInProgress {
		promotionPopupCenter := g.graphicsState.getPromotionPopupOrigin(g.promotionLifeCycle.promotionSquare)
		g.drawPromotionBox(chessBoardImage, promotionPopupCenter)
	}

	// Rotate the ChessBoard so that white is at the bottom
	// Now draw the offscreen buffer to the screen with rotation
	op := &ebiten.DrawImageOptions{}

	// Set the origin to the center of the screen
	w, h := float64(g.boardWidth()), float64(g.boardHeight())
	op.GeoM.Translate(-w/2, -h/2)

	// Rotate 180 degrees (π radians)
	op.GeoM.Rotate(math.Pi)

	// Move back to screen coordinates
	op.GeoM.Translate(w/2, h/2)

	// Draw the rotated image to the screen
	screen.DrawImage(chessBoardImage, op)
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

// Draws a white Box with width and height equal to twice the width and height of a chess square.
// The box has black borders that are rounded at the corner, and has a vertical and horizontal line running through its center in a symmetrical cross
// The box is drawn to the screen with its center at the coordinates specified by boxCenter
func (g *ChessGame) drawPromotionBox(screen *ebiten.Image, boxCenter point) {
	boxWidth := float32(g.squareWidth() * 2)
	boxHeight := float32(g.squareHeight() * 2)

	// Calculate top-left corner from center
	x := boxCenter.x - float64(boxWidth)/2
	y := boxCenter.y - float64(boxHeight)/2

	// Draw white box
	vector.DrawFilledRect(screen, float32(x), float32(y), boxWidth, boxHeight, color.RGBA{255, 255, 255, 255}, true)

	// Draw black borders with rounded corners
	borderWidth := float32(4)
	vector.StrokeRect(screen, float32(x), float32(y), boxWidth, boxHeight, borderWidth, color.RGBA{0, 0, 0, 255}, true)

	// Draw vertical line
	vector.StrokeLine(screen, float32(x+float64(boxWidth)/2), float32(y), float32(x+float64(boxWidth)/2), float32(y+float64(boxHeight)), borderWidth, color.RGBA{0, 0, 0, 255}, true)

	// Draw horizontal line
	vector.StrokeLine(screen, float32(x), float32(y+float64(boxHeight)/2), float32(x+float64(boxWidth)), float32(y+float64(boxHeight)/2), borderWidth, color.RGBA{0, 0, 0, 255}, true)
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
