package chessgame

// Methods for basic dimensions as functions of window height and width
func (gr *graphicsState) boardWidth() int {
	return min(gr.windowWidth, gr.windowHeight) + 4
}

func (gr *graphicsState) boardHeight() int {
	return min(gr.windowWidth, gr.windowHeight) + 4
}

func (gr *graphicsState) squareWidth() float64 {
	return float64(gr.boardWidth() / 8)
}

func (gr *graphicsState) squareHeight() float64 {
	return float64(gr.boardHeight() / 8)
}

func (gr *graphicsState) pieceWidth() float64 {
	return gr.squareWidth() * pieceScale
}

func (gr *graphicsState) pieceHeight() float64 {
	return gr.squareHeight() * pieceScale
}
func (gr *graphicsState) getSquareOfMousePosition(x int, y int) (vector2, bool) {
	file := int(float64(x) / gr.squareWidth())
	rank := int(float64(y) / gr.squareHeight())

	// Check if coordinates are within bounds
	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return nilSquare, true
	}

	// Apply rotation based on rotationFactor
	return rotateChessCoord(vector2{file, rank}, -rotationFactor), false
}

func (gr *graphicsState) getPromotionPopupOrigin(promotionSquare vector2) point {
	//minEdgeOfBoardDistance := gr.squareWidth() / 4

	//promotionSquare = rotateChessCoord(promotionSquare, -rotationFactor)

	// Get top corner of square
	popUpOrigin := point{float64(promotionSquare.x) * gr.squareWidth(), float64(promotionSquare.y) * gr.squareHeight()}

	// Translate into center of square from corner
	popUpOrigin = popUpOrigin.add(point{gr.squareWidth() / 2, gr.squareHeight() / 2})

	// Translate vertically away from edge of board
	if promotionSquare.y == 0 {
		// Move so popup is at edge of board
		popUpOrigin = popUpOrigin.add(point{0, gr.squareHeight()})
	} else if promotionSquare.y == 7 { // Top of board case
		// Move so popup is at edge of board
		popUpOrigin = popUpOrigin.add(point{0, -gr.squareHeight()})
	} else {
		panic("promotion square not on top or bottom of board")
	}

	// Translate away from sides if applicable
	if promotionSquare.x == 0 {
		popUpOrigin.x += gr.squareWidth()
	}
	if promotionSquare.x == 7 {
		popUpOrigin.x -= gr.squareHeight()
	}

	return popUpOrigin
}
