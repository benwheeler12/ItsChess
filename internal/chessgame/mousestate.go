package chessgame

func (ms *mouseLifeCycle) resetMouseState() {
	ms.selectedPiece = emptyPiece
	ms.mouseClickedOnInvalidSquare = false
	ms.selectedSquare = nilSquare
	ms.possibleMoveSquares = nil
}
