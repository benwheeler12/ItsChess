package itschess

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

type chessBoardGraphic struct {
	rotationTheta       float64
	reflection          int
	pieceImages         map[chessPiece]*ebiten.Image
	width               int
	height              int
	clickedSquare       vector2
	possibleMoveSquares []vector2
	promotionSquare     vector2
}

const pieceScale = .9

func (cbg *chessBoardGraphic) init(rotationTheta float64, reflection int, width int, height int) {
	cbg.rotationTheta = rotationTheta
	cbg.reflection = reflection
	cbg.width = width
	cbg.height = height
	cbg.loadPieceImages()

	cbg.clickedSquare = nilSquare
	cbg.possibleMoveSquares = nil
	cbg.promotionSquare = nilSquare
}

func (cbg *chessBoardGraphic) drawChessBoard(chessBoard *chessBoard) *ebiten.Image {
	chessBoardImage := ebiten.NewImage(cbg.boardWidth(), cbg.boardHeight())

	chessBoardImage.Fill(color.RGBA{255, 255, 255, 255})

	// TODO implement checkmate animation
	if chessBoard.playerInCheckMate(white) || chessBoard.playerInCheckMate(black) {
		return chessBoardImage
	}

	// Draw Board
	for file := range 8 {
		for rank := range 8 {

			squareX := cbg.squareWidth() * float64(file)
			squareY := cbg.squareHeight() * float64(rank)

			cbg.drawChessSquare(chessBoardImage, chessBoard, squareX, squareY, vector2{file, rank})
			continue
		}
	}

	// Draw the clicked piece
	/**if cbg.mouseLifeCycle.selectedPiece != emptyPiece {
		mouseX, mouseY := ebiten.CursorPosition()
		cbg.drawChessPiece(cbg.selectedPiece, float64(mouseX), float64(mouseY), chessBoardImage)
	}*/
	promotionPopupCenter := cbg.getPromotionPopupOrigin(vector2{0, 0})
	cbg.drawPromotionBox(chessBoardImage, promotionPopupCenter)

	promotionPopupCenter = cbg.getPromotionPopupOrigin(vector2{0, 7})
	cbg.drawPromotionBox(chessBoardImage, promotionPopupCenter)

	promotionPopupCenter = cbg.getPromotionPopupOrigin(vector2{7, 7})
	cbg.drawPromotionBox(chessBoardImage, promotionPopupCenter)

	promotionPopupCenter = cbg.getPromotionPopupOrigin(vector2{7, 0})
	cbg.drawPromotionBox(chessBoardImage, promotionPopupCenter)

	if cbg.promotionSquare != nilSquare {
		promotionPopupCenter := cbg.getPromotionPopupOrigin(cbg.promotionSquare)
		cbg.drawPromotionBox(chessBoardImage, promotionPopupCenter)
	}

	// Rotate the ChessBoard so that white is at the bottom
	// Now draw the offscreen buffer to the screen with rotation
	op := &ebiten.DrawImageOptions{}

	// Set the origin to the center of the screen
	w, h := float64(cbg.boardWidth()), float64(cbg.boardHeight())
	op.GeoM.Translate(-w/2, -h/2)
	_, _ = w, h

	// Rotate 180 degrees (π radians)
	op.GeoM.Rotate(cbg.rotationTheta)

	// Scale by -1 in the x direction to flip horizontally
	if cbg.reflection == 1 {
		op.GeoM.Scale(-1, 1)
	}

	// Move back to screen coordinates
	op.GeoM.Translate(w/2, h/2)

	rotatedImage := ebiten.NewImage(cbg.boardWidth(), cbg.boardHeight())

	rotatedImage.DrawImage(chessBoardImage, op)

	return rotatedImage
}

// Loads Piece Images based on width and height of graphic
func (cbg *chessBoardGraphic) loadPieceImages() {
	cbg.pieceImages = map[chessPiece]*ebiten.Image{
		{pawn, white}:   cbg.getPieceImage(whitePawnBytes),
		{pawn, black}:   cbg.getPieceImage(blackPawnBytes),
		{bishop, white}: cbg.getPieceImage(whiteBishopBytes),
		{bishop, black}: cbg.getPieceImage(blackBishopBytes),
		{knight, white}: cbg.getPieceImage(whiteKnightBytes),
		{knight, black}: cbg.getPieceImage(blackKnightBytes),
		{rook, white}:   cbg.getPieceImage(whiteRookBytes),
		{rook, black}:   cbg.getPieceImage(blackRookBytes),
		{queen, white}:  cbg.getPieceImage(whiteQueenBytes),
		{queen, black}:  cbg.getPieceImage(blackQueenBytes),
		{king, white}:   cbg.getPieceImage(whiteKingBytes),
		{king, black}:   cbg.getPieceImage(blackKingBytes),
	}
}

func (cbg *chessBoardGraphic) getPieceImage(pieceBytes []byte) *ebiten.Image {

	reader := bytes.NewReader(pieceBytes)

	// Decode SVG
	icon, err := oksvg.ReadIconStream(reader)
	if err != nil {
		log.Fatal(err)
	}

	// Set SVG viewport
	icon.SetTarget(0, 0, cbg.pieceWidth(), cbg.pieceHeight())

	w := cbg.pieceWidth()
	_ = w
	// Create an RGBA image
	img := image.NewRGBA(image.Rect(0, 0, int(math.Round(cbg.pieceWidth())), int(math.Round(cbg.pieceHeight()))))
	drawer := rasterx.NewDasher(
		int(math.Round(cbg.pieceWidth())),
		int(math.Round(cbg.pieceHeight())),
		rasterx.NewScannerGV(
			int(math.Round(cbg.pieceWidth())),
			int(math.Round(cbg.pieceHeight())),
			img,
			img.Bounds(),
		),
	)
	icon.Draw(drawer, 1.0)

	return ebiten.NewImageFromImage(img)
}

// Methods for basic dimensions as functions of window height and width
func (cbg *chessBoardGraphic) boardWidth() int {
	return min(cbg.width, cbg.height) + 4
}

func (cbg *chessBoardGraphic) boardHeight() int {
	return min(cbg.width, cbg.height) + 4
}

func (cbg *chessBoardGraphic) squareWidth() float64 {
	return float64(cbg.boardWidth() / 8)
}

func (cbg *chessBoardGraphic) squareHeight() float64 {
	return float64(cbg.boardHeight() / 8)
}

func (cbg *chessBoardGraphic) pieceWidth() float64 {
	return cbg.squareWidth() * pieceScale
}

func (cbg *chessBoardGraphic) pieceHeight() float64 {
	return cbg.squareHeight() * pieceScale
}
func (cbg *chessBoardGraphic) getSquareOfMousePosition(x int, y int) (vector2, bool) {
	file := int(float64(x) / cbg.squareWidth())
	rank := int(float64(y) / cbg.squareHeight())

	// Check if coordinates are within bounds
	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return nilSquare, true
	}

	// Apply rotation based on rotationFactor
	return rotateChessCoord(vector2{file, rank}, -rotationFactor), false
}

func (cbg *chessBoardGraphic) getPromotionPopupOrigin(promotionSquare vector2) point {

	// Get top corner of square
	popUpOrigin := point{float64(promotionSquare.x) * cbg.squareWidth(), float64(promotionSquare.y) * cbg.squareHeight()}

	// Translate into center of square from corner
	popUpOrigin = popUpOrigin.add(point{cbg.squareWidth() / 2, cbg.squareHeight() / 2})

	// Translate vertically away from edge of board
	if promotionSquare.y == 0 {
		// Move so popup is at edge of board
		popUpOrigin = popUpOrigin.add(point{0, cbg.squareHeight()})
	} else if promotionSquare.y == 7 { // Top of board case
		// Move so popup is at edge of board
		popUpOrigin = popUpOrigin.add(point{0, -cbg.squareHeight()})
	} else {
		panic("promotion square not on top or bottom of board")
	}

	// Translate away from sides if applicable
	if promotionSquare.x == 0 {
		popUpOrigin.x += cbg.squareWidth()
	}
	if promotionSquare.x == 7 {
		popUpOrigin.x -= cbg.squareHeight()
	}

	return popUpOrigin
}

// square refers to the canonical chess square that this function call will draw
func (cbg *chessBoardGraphic) drawChessSquare(screen *ebiten.Image, chessBoard *chessBoard, x float64, y float64, square vector2) {

	// Calculate Color
	squareColor := lightSquareColor

	if (square.x+square.y)%2 == 0 {
		squareColor = darkSquareColor
	}

	if contains(cbg.possibleMoveSquares, square) {
		squareColor = tintColor(squareColor, color.RGBA{255, 255, 0, 1}, 0.3)
	}

	vector.DrawFilledRect(screen,
		float32(x),
		float32(y),
		float32(cbg.squareWidth()),
		float32(cbg.squareHeight()),
		squareColor,
		true)

	if cbg.clickedSquare != square {
		cbg.drawChessPiece(chessBoard.getPiece(square), x, y, screen)
	}

}

// Draws a white Box with width and height equal to twice the width and height of a chess square.
// The box has black borders that are rounded at the corner, and has a vertical and horizontal line running through its center in a symmetrical cross
// The box is drawn to the screen with its center at the coordinates specified by boxCenter
func (cbg *chessBoardGraphic) drawPromotionBox(screen *ebiten.Image, boxCenter point) {
	boxWidth := float32(cbg.squareWidth() * 2)
	boxHeight := float32(cbg.squareHeight() * 2)

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

func (cbg *chessBoardGraphic) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Window Size has changed!  Reinitialize piece images
	if outsideWidth != cbg.width || outsideHeight != cbg.height {
		// Order matter here
		cbg.width = outsideWidth
		cbg.height = outsideHeight
		cbg.loadPieceImages()
	}

	return outsideWidth, outsideHeight
}

// x and y represnt the mid point of the image.
func (cbg *chessBoardGraphic) drawChessPiece(piece chessPiece, x float64, y float64, screen *ebiten.Image) {

	image, ok := cbg.pieceImages[piece]
	if !ok {
		return
	}

	op := &ebiten.DrawImageOptions{}

	// Rotate Image
	op.GeoM.Translate(-cbg.pieceWidth()/2, -cbg.pieceHeight()/2)
	op.GeoM.Rotate(cbg.rotationTheta)
	op.GeoM.Translate(cbg.pieceWidth()/2, cbg.pieceHeight()/2)

	// Translate Image
	op.GeoM.Translate(x, y)

	// Translate Additional pixels accounting for difference in square size vs piece size
	squareOffsetX := (cbg.squareWidth() - cbg.pieceWidth()) / 2
	squareOffsetY := (cbg.squareHeight() - cbg.pieceHeight()) / 2
	op.GeoM.Translate(squareOffsetX, squareOffsetY)

	op.Filter = ebiten.FilterNearest

	screen.DrawImage(image, op)
}
